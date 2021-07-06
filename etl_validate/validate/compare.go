package main

import (
	"bufio"
	"container/list"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)
var wgbig sync.WaitGroup //定义一个同步等待的组
var wgitem sync.WaitGroup //定义一个同步等待的组
var fieldSeperator = ","
func compare()(){

	var dir string
	flag.StringVar(&dir,"d","目录名","目录名")
	flag.Parse()
	//对比的表名
	var tables []string = []string{"tb_cen_account_o_storeinven","tb_cen_storenotavailableqty","tb_gos_stock_stockpreemption"}
	//对比的表列名
	var fields [][]string = [][]string{[]string{"pk", "fk", "lineid", "lastmodifytime", "createtime", "branchid", "prodid", "invbalqty", "invbalamt", "storeid", "deleteflag", "note", "version"},
		[]string{"pk", "createtime", "lastmodifytime", "version", "branchid", "storeid", "prodid", "notavailableqty", "preassignedqty", "runningno", "note"},
		[]string{"pk", "fk", "createtime", "lastmodifytime","version", "lineid", "branchid", "deleteflag", "note", "preemptionpreemption", "prodid", "lotno", "quantity", "rowguid", "billid", "whseid", "storeid", "billguid", "opid", "custid", "custno", "custname"}}
	//需要对比的表列名 对比的列名才会打印出来
	var compareFields [][]string = [][]string{[]string{"branchid","prodid","invbalqty", "invbalamt", "storeid", "deleteflag",  "version","lastmodifytime"},
		[]string{"version", "branchid", "storeid", "prodid", "notavailableqty", "preassignedqty","lastmodifytime","deleteflag" },
		[]string{"version", "branchid", "storeid", "prodid", "quantity","lastmodifytime","deleteflag"}}
	for idx,t := range tables{
		wgbig.Add(1)
		go compareInner(dir,t,fields[idx],compareFields[idx])
	}
	wgbig.Wait()
	println("finished")

}

// 比较  erp_tb_cen_account_o_storeinven
func compareInner(dir string,tablename string,fields []string,compareFields []string)(){
	defer wgbig.Done()
	wgitem.Add(2)
	//var filename string = "tb_cen_account_o_storeinven"
	var midData map[string]string = make(map[string]string)
	var erpData map[string] string = make(map[string]string)
	var compareIndex []int = getCompareIndex(&fields,&compareFields)
	//文件内容加载到map
	go load2Map(&midData,dir+"/"+tablename)
	go load2Map(&erpData,dir+"/erp_"+tablename)
	wgitem.Wait()
	println("数据准备完毕，开始对比。。")
	//本地化校验库多余 和本地化校验库丢失数据列表
	var midMissingDatas = list.New()
	var midMoreDatas = list.New()
	var diffDatas = list.New()
	var fieldtitle string = "pk\t"
	for _,v := range compareIndex{
		fieldtitle += fields[v]+"\t"+fields[v]+"\t"
	}
	diffDatas.PushBack(fieldtitle)

	var fieldLen = len(fields)
	for pk,vMid := range midData{
		var vErp = erpData[pk]
		if vErp == ""{
			//本地化校验库pk不在erp库中 计入本地化校验库多余数据
			midMoreDatas.PushBack(pk)
		}else{
			vMidSplit := strings.Split(vMid, fieldSeperator)
			vErpsplit := strings.Split(vErp, fieldSeperator)
			if len(vMidSplit) != fieldLen || len(vErpsplit) != fieldLen{
				fmt.Printf("数据格式错误,table:%v,pk:%v",tablename,pk)
				continue
			}

			//计算数据签名（需要比较的字段拼接)
			signMid,signErp := "",""
			for _,v := range compareIndex{
				if fields[v] == "lastmodifytime"{
					//修改时间不参与签名
					continue
				}
				signMid += zeroNumberHandle(vMidSplit[v]);
				signErp += zeroNumberHandle(vErpsplit[v]);
			}
			signMid = strings.ReplaceAll(signMid,"\n","")
			signErp = strings.ReplaceAll(signErp,"\n","")

			if signErp != signMid{
				var linestr string = pk+"\t"
				for _,i := range compareIndex{
					linestr += vErpsplit[i]+"\t"+vMidSplit[i]+"\t"
				}
				diffDatas.PushBack(linestr)
			}
		}
	}

	//收集本地化校验库丢失的数据
	for pk,_ := range erpData{
		var vMid = midData[pk]
		if vMid == ""{
			midMissingDatas.PushBack(pk)
		}
	}
	targetFileName := fmt.Sprintf("%v#本地化校验库",tablename)
	write2File(dir+"/"+targetFileName+"丢失数据.txt",midMissingDatas)
	write2File(dir+"/"+targetFileName+"多余数据.txt",midMoreDatas)
	write2File(dir+"/"+targetFileName+"不一致数据.txt",diffDatas)
}

func main()(){
	compare()
}

/**
数字小数点后如果是0 统一处理成0 而不是0.0 0.00 等等
*/
func zeroNumberHandle(a string)(r string){
	if !strings.Contains(a,".0"){
		return a;
	}
	_,err:=strconv.ParseFloat(a,32)
	if err != nil{
		//不是个数字 原样返回
		return a
	}
	r = a
	for{
		if strings.Contains(r,".0") && strings.LastIndex(r,"0") == len(r)-1{
			r = r[0:len(r)-1]
		}else{
			break;
		}
	}
	if strings.LastIndex(r,".") == len(r)-1{
		r = r[0:len(r)-1]
	}
	return r
}
func getCompareIndex(all *[]string,compare *[]string)(r []int){
	r = []int{}
	for i, value := range *all {
		isCompare := false
		for _, _value := range *compare {
			if value == _value{
				isCompare = true
				break;
			}
		}
		if isCompare{
			r = append(r, i)
		}
	}
	return r
}

func load2Map(m *map[string]string,fpath string){
	defer wgitem.Done()
	f, err := os.Open(fpath)
	if err != nil {
		fmt.Printf("%v文件读取失败 退出..",fpath)
		os.Exit(1)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if line != ""{
			(*m)[getPk(line)] = line
		}
		if err != nil {
			if err == io.EOF { //读取结束，会报EOF
				break
			}
			break
		}

	}
}
/**
获取第一个字段 pk
*/
func getPk(line string)(pk string){
	var limits []string = strings.Split(line,",")
	return limits[0]
}

func write2File(filePath string,l *list.List)(){
	wf, error := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if error != nil{
		fmt.Printf("写入文件异常:%v",error)
	}
	defer wf.Close()
	writer := bufio.NewWriter(wf)
	for e:=l.Front();e!=nil; e= e.Next(){
		v := strings.ReplaceAll(fmt.Sprintf("%v",e.Value),"\n","")
		writer.WriteString(v+"\n")
	}
	writer.Flush()
}
