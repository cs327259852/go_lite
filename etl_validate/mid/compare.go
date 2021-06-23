package main

import (
	"bufio"
	"container/list"
	"flag"
	"fmt"
	"io"
	"os"
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
	//需要对比的表列名
	var compareFields [][]string = [][]string{[]string{"branchid","prodid","invbalqty", "invbalamt", "storeid", "deleteflag",  "version"},
										[]string{"version", "branchid", "storeid", "prodid", "notavailableqty", "preassignedqty" },
										[]string{"version", "branchid", "storeid", "prodid", "quantity"}}
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
	//中间库多余 和中间库丢失数据列表
	var midMissingDatas = list.New()
	var midMoreDatas = list.New()
	var diffDatas = list.New()
	var fieldtitle string = "左ERP\t右中间库\t"
	for _,v := range compareIndex{
		fieldtitle += fields[v]+"\t"+fields[v]+"\t"
	}
	diffDatas.PushBack(fieldtitle)

	var fieldLen = len(fields)
	for pk,vMid := range midData{
		var vErp = erpData[pk]
		if vErp == ""{
			//中间库pk不在erp库中 计入中间库多余数据
			midMoreDatas.PushBack(pk)
		}else{
			vMidSplit := strings.Split(vMid, fieldSeperator)
			vErpsplit := strings.Split(vErp, fieldSeperator)
			if len(vMidSplit) != fieldLen || len(vErpsplit) != fieldLen{
				fmt.Printf("数据格式错误:%v",vMid)
				continue
			}

			//计算数据签名（需要比较的字段拼接)
			signMid,signErp := "",""
			for _,v := range compareIndex{
				signMid += vMidSplit[v]
				signErp += vErpsplit[v]
			}
			signMid = strings.ReplaceAll(signMid,"\n","")
			signErp = strings.ReplaceAll(signErp,"\n","")

			if signErp != signMid{
				var linestr string = "左ERP\t右中间库\t"
				for _,i := range compareIndex{
					linestr += vErpsplit[i]+"\t"+vMidSplit[i]+"\t"
				}
				diffDatas.PushBack(linestr)
			}
		}
	}

	//收集中间库丢失的数据
	for pk,_ := range erpData{
		var vMid = midData[pk]
		if vMid == ""{
			midMissingDatas.PushBack(pk)
		}
	}
	targetFileName := fmt.Sprintf("%v#中间库",tablename)
	write2File(dir+"/"+targetFileName+"丢失数据.txt",midMissingDatas)
	write2File(dir+"/"+targetFileName+"多余数据.txt",midMoreDatas)
	write2File(dir+"/"+targetFileName+"不一致数据.txt",diffDatas)
}

func main()(){
	compare()
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
		print("文件读取失败 退出..")
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
