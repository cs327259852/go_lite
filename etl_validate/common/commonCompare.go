package common

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

var wgbig sync.WaitGroup  //定义一个同步等待的组
var wgitem sync.WaitGroup //定义一个同步等待的组
var fieldSeperator = ","

func CommonCompare(dir *string, sdPare *[][]string, tables *[]string, fields *[][]string, compareFields *[][]string) {
	for idx, t := range *tables {
		for _, sdv := range *sdPare {
			if len(sdv) < 2 {
				fmt.Printf("源目标库对比数量错误,%v,%v\n", sdv, t)
				continue
			}
			//源文件
			sfileName := *dir + "/" + sdv[0] + t
			//目标文件
			dfileName := *dir + "/" + sdv[1] + t
			if fileExists(sfileName) && fileExists(dfileName) {
				wgbig.Add(1)
				go compareInner(*dir, sfileName, dfileName, strings.Replace(sdv[0]+"2"+sdv[1]+"-"+t, "_", "", 2), (*fields)[idx], (*compareFields)[idx])
			} else {
				fmt.Printf("对比文件不存在%v-%v，跳过..\n", sfileName, dfileName)
			}
		}

	}
	wgbig.Wait()
	println("finished")

}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

// 比较
func compareInner(dir string, sfileName string, dfileName string, resultname string, fields []string, compareFields []string) {
	defer wgbig.Done()
	wgitem.Add(2)
	//var filename string = "tb_cen_account_o_storeinven"
	var destData map[string]string = make(map[string]string)
	var sourceData map[string]string = make(map[string]string)
	var compareIndex []int = getCompareIndex(&fields, &compareFields)
	//文件内容加载到map
	go load2Map(&sourceData, sfileName)
	go load2Map(&destData, dfileName)
	wgitem.Wait()
	println("数据准备完毕，开始对比。。")
	//目标库多余 和中间库丢失数据列表
	var destMissingDatas = list.New()
	var destMoreDatas = list.New()
	var diffDatas = list.New()
	var fieldtitle string = "pk\t"
	for _, v := range compareIndex {
		fieldtitle += fields[v] + "\t" + fields[v] + "\t"
	}
	diffDatas.PushBack(fieldtitle)

	var fieldLen = len(fields)
	for pk, vMid := range destData {
		var vErp = sourceData[pk]
		if vErp == "" {
			//目标库pk不在源库中 计入目标多余数据
			destMoreDatas.PushBack(pk)
		} else {
			vMidSplit := strings.Split(vMid, fieldSeperator)
			vErpsplit := strings.Split(vErp, fieldSeperator)
			if len(vMidSplit) != fieldLen || len(vErpsplit) != fieldLen {
				fmt.Printf("数据格式错误,table:%v,pk:%v,跳过...", resultname, pk)
				continue
			}

			//计算数据签名（需要比较的字段拼接)
			signMid, signErp := "", ""
			for _, v := range compareIndex {
				if fields[v] == "lastmodifytime" {
					//修改时间不参与签名
					continue
				}
				m := vMidSplit[v]
				e := vErpsplit[v]

				m = strings.TrimSpace(strings.ReplaceAll(m, "\n", ""))
				e = strings.TrimSpace(strings.ReplaceAll(e, "\n", ""))
				signMid += zeroNumberHandle(m)
				signErp += zeroNumberHandle(e)
			}
			if signErp != signMid {
				var linestr string = pk + "\t"
				for _, i := range compareIndex {
					m := vMidSplit[i]
					e := vErpsplit[i]
					m = strings.TrimSpace(strings.ReplaceAll(m, "\n", ""))
					e = strings.TrimSpace(strings.ReplaceAll(e, "\n", ""))
					linestr += e + "\t" + m + "\t"
				}
				diffDatas.PushBack(linestr)
			}
		}
	}

	//收集中间库丢失的数据
	for pk, _ := range sourceData {
		var vMid = destData[pk]
		if vMid == "" {
			destMissingDatas.PushBack(pk)
		}
	}
	targetFileName := fmt.Sprintf("%v#", resultname)
	Write2File(dir+"/"+targetFileName+"丢失数据.crt", destMissingDatas)
	Write2File(dir+"/"+targetFileName+"多余数据.crt", destMoreDatas)
	Write2File(dir+"/"+targetFileName+"不一致数据.crt", diffDatas)
}

/**
数字小数点后如果是0 统一处理成0 而不是0.0 0.00 等等
*/
func zeroNumberHandle(a string) (r string) {
	if !strings.Contains(a, ".") {
		return a
	}
	_, err := strconv.ParseFloat(a, 32)
	if err != nil {
		//不是个数字 原样返回
		return a
	}
	hasDotMinus := strings.Index(a, "-.") == 0
	if hasDotMinus {
		a = strings.Replace(a, "-.", "-0.", -1)
	}
	r = a
	for {
		if strings.LastIndex(r, "0") == len(r)-1 {
			r = r[0 : len(r)-1]
		} else {
			break
		}
	}
	if strings.LastIndex(r, ".") == len(r)-1 {
		r = r[0 : len(r)-1]
	}
	if strings.LastIndex(r, ".") == 0 {
		r = "0" + r
	}
	return r
}
func getCompareIndex(all *[]string, compare *[]string) (r []int) {
	r = []int{}
	for i, value := range *all {
		isCompare := false
		for _, _value := range *compare {
			if value == _value {
				isCompare = true
				break
			}
		}
		if isCompare {
			r = append(r, i)
		}
	}
	return r
}

func load2Map(m *map[string]string, fpath string) {
	defer wgitem.Done()
	f, err := os.Open(fpath)
	if err != nil {
		fmt.Printf("%v文件读取失败 跳过..", fpath)
		return
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if line != "" {
			line = strings.TrimSpace(strings.ReplaceAll(line, "\n", ""))
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
func getPk(line string) (pk string) {
	var limits []string = strings.Split(line, ",")
	return limits[0]
}

func Write2File(filePath string, l *list.List) {
	wf, error := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if error != nil {
		fmt.Printf("写入文件异常:%v", error)
	}
	defer wf.Close()
	writer := bufio.NewWriter(wf)
	for e := l.Front(); e != nil; e = e.Next() {
		v := strings.ReplaceAll(fmt.Sprintf("%v", e.Value), "\n", "")
		writer.WriteString(v + "\n")
	}
	writer.Flush()
}
