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
	"study/etl_validate/common"
	"sync"
)

var wgbig sync.WaitGroup
var branchStoreidMap = make(map[string]string)

func main() {
	var dir string
	flag.StringVar(&dir, "d", "目录名", "目录名")
	flag.Parse()
	//初始化分公司仓库map
	wgbig.Add(7)
	initBranchStoreidMap()

	//库存计算的相关表
	s1 := make(map[string]string)
	//加载1库存数
	loadStorage1ItemMap(&s1, dir+"/tb_cen_account_o_storeinven")

	s2 := make(map[string]string)
	loadStorage2ItemMap(&s2, dir+"/tb_cen_storenotavailableqty")

	s3 := make(map[string]string)
	loadStorage3ItemMap(&s3, dir+"/tb_gos_stock_stockpreemption")

	s4 := make(map[string]string)
	loadStorage4ItemMap(&s4, dir+"/tb_common_productreserves")

	storageMap := make(map[string]string)
	loadStorageItemMap(&storageMap, dir+"/tb_merchandise_storage")
	//加载商品编码和商品id字典
	var prodidprodnoMap map[string]string = make(map[string]string)
	loadProdnoProdIdMap(&prodidprodnoMap, dir+"/vw_common_prod")

	wgbig.Wait()

	var diffDatas = list.New()
	for prodno, actualStorage := range storageMap {
		prodid := prodidprodnoMap[prodno]
		if prodid == "" {
			fmt.Println("未找到商品编码和id映射，跳过..")
			continue
		}
		a, err1 := strconv.ParseFloat(s1[prodid], 32)
		b, err2 := strconv.ParseFloat(s2[prodid], 32)
		c, err3 := strconv.ParseFloat(s3[prodid], 32)
		d, err4 := strconv.ParseFloat(s4[prodid], 32)
		ac, err5 := strconv.ParseFloat(actualStorage, 32)
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			fmt.Println("各个库存转换数字异常:%v,跳过...", prodid)
			continue
		}
		if err5 != nil {
			fmt.Println("实际库存转换数字异常:%v,跳过...", prodno)
			continue
		}
		_ac := a - b - c - d
		if ac != _ac {
			diffDatas.PushBack(fmt.Sprintf("%v,%v,%v", prodno, _ac, ac))
		}

	}

	common.Write2File(dir+"/库存不相等数据.txt", diffDatas)
	fmt.Println("完成对比库存...")
}

func initBranchStoreidMap() {
	defer wgbig.Done()
	branchStoreidMap["FDG"] = ""
	branchStoreidMap["FDW"] = ""
	branchStoreidMap["FDC"] = ""
}

func loadProdnoProdIdMap(m *map[string]string, fpath string) {
	defer wgbig.Done()
	f, err := os.Open(fpath)
	if err != nil {
		fmt.Printf("%v文件读取失败 退出..", fpath)
		os.Exit(1)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if line != "" {
			arr := strings.Split(line, ",")
			if len(arr) < 4 {
				fmt.Print("商品信息不全，跳过...")
				continue
			}
			branchid := arr[1]
			prodid := arr[2]
			prodno := arr[3]
			(*m)[branchid+"-"+prodno] = branchid + "-" + prodid
		}
		if err != nil {
			if err == io.EOF { //读取结束，会报EOF
				break
			}
			break
		}

	}
}

func loadStorage1ItemMap(m *map[string]string, fpath string) {
	defer wgbig.Done()
	f, err := os.Open(fpath)
	if err != nil {
		fmt.Printf("%v文件读取失败 退出..", fpath)
		os.Exit(1)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if line != "" {
			arr := strings.Split(line, ",")
			if len(arr) < 8 {
				continue
			}
			qty := arr[7]
			branchid := arr[5]
			prodid := arr[6]
			if branchStoreidMap[branchid] != arr[9] {
				//仓库id不一样
				continue
			}
			(*m)[branchid+"-"+prodid] = qty

		}
		if err != nil {
			if err == io.EOF { //读取结束，会报EOF
				break
			}
			break
		}

	}
}

func loadStorage2ItemMap(m *map[string]string, fpath string) {
	defer wgbig.Done()
	f, err := os.Open(fpath)
	if err != nil {
		fmt.Printf("%v文件读取失败 退出..", fpath)
		os.Exit(1)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if line != "" {
			arr := strings.Split(line, ",")
			if len(arr) < 8 {
				continue
			}
			branchid := arr[4]
			if branchStoreidMap[branchid] != arr[5] {
				//仓库id不一样
				continue
			}
			qty1, err := strconv.ParseFloat(arr[6], 32)
			if err != nil {
				continue
			}
			qty2, err := strconv.ParseFloat(arr[7], 32)
			if err != nil {
				continue
			}
			qty := fmt.Sprintf("%f", qty1+qty2)
			prodid := arr[6]
			(*m)[branchid+"-"+prodid] = qty

		}
		if err != nil {
			if err == io.EOF { //读取结束，会报EOF
				break
			}
			break
		}

	}
}

func loadStorage3ItemMap(m *map[string]string, fpath string) {
	defer wgbig.Done()
	f, err := os.Open(fpath)
	if err != nil {
		fmt.Printf("%v文件读取失败 退出..", fpath)
		os.Exit(1)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if line != "" {
			arr := strings.Split(line, ",")
			if len(arr) < 17 {
				continue
			}
			if arr[7] == "1" {
				//deleteflag = 1
				continue
			}
			branchid := arr[6]
			if branchStoreidMap[branchid] != arr[16] {
				//仓库id不一样
				continue
			}

			qty := arr[12]
			prodid := arr[10]
			oldqty := (*m)[branchid+"-"+prodid]
			if oldqty != "" {
				oldqtyf, err := strconv.ParseFloat(oldqty, 32)
				if err != nil {
					qtyf, err := strconv.ParseFloat(qty, 32)
					if err != nil {
						qty = fmt.Sprintf("%f", qtyf+oldqtyf)
					} else {
						fmt.Println("预占库存数量计算异常:%v-%v", branchid, prodid)
					}
				}
			}
			(*m)[branchid+"-"+prodid] = qty
		}
		if err != nil {
			if err == io.EOF { //读取结束，会报EOF
				break
			}
			break
		}

	}
}

func loadStorage4ItemMap(m *map[string]string, fpath string) {
	defer wgbig.Done()
	f, err := os.Open(fpath)
	if err != nil {
		fmt.Printf("%v文件读取失败 退出..", fpath)
		os.Exit(1)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if line != "" {
			arr := strings.Split(line, ",")
			if len(arr) < 8 {
				continue
			}
			branchid := arr[2]
			if branchStoreidMap[branchid] != arr[7] {
				//仓库id不一样
				continue
			}

			prodid := arr[4]
			(*m)[branchid+"-"+prodid] = arr[5]
		}
		if err != nil {
			if err == io.EOF { //读取结束，会报EOF
				break
			}
			break
		}

	}
}

func loadStorageItemMap(m *map[string]string, fpath string) {
	defer wgbig.Done()
	f, err := os.Open(fpath)
	if err != nil {
		fmt.Printf("%v文件读取失败 退出..", fpath)
		os.Exit(1)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if line != "" {
			arr := strings.Split(line, ",")
			if len(arr) < 4 {
				fmt.Println("业务库实际库存信息不全，跳过...")
				continue
			}
			branchid := arr[0]

			prodno := arr[1]
			(*m)[branchid+"-"+prodno] = arr[2]
		}
		if err != nil {
			if err == io.EOF { //读取结束，会报EOF
				break
			}
			break
		}

	}
}