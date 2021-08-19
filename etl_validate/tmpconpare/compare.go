package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"study/etl_validate/common2"
	"time"
)

func main() {
	var dir,tablefile string
	flag.StringVar(&dir, "d", "目录名", "目录名")
	flag.StringVar(&tablefile, "t", "表文件", "表文件")
	flag.Parse()
	now := time.Now()
	//对比的库源和目标库前缀
	var srouceDestPair [][]string = [][]string{[]string{"erp_", "mid_"},
	}
	//对比的表名
	var tables *[]string = getAllTables(tablefile)
	var fields []string = []string{"pk","lastmodifytime"}
	var compareFields []string = []string{"lastmodifytime"}
	common2.CommonCompare2(&dir, &srouceDestPair, tables, fields, compareFields)
	fmt.Printf("完成对比，耗时:%v\n", time.Since(now))
}

func getAllTables( fpath string) *[]string{
	//defer wgitem.Done()
	f, err := os.Open(fpath)
	if err != nil {
		fmt.Printf("%v表文件读取失败 exit..", fpath)
		os.Exit(-1)
	}
	var tables []string = make([]string,500)
	defer f.Close()
	buf := bufio.NewReader(f)
	i := 0
	for {
		line, err := buf.ReadString('\n')
		if line != "" {
			line = strings.TrimSpace(strings.ReplaceAll(line, "\n", ""))
			tables[i] = line
			i++
		}
		if err != nil {
			if err == io.EOF { //读取结束，会报EOF
				break
			}
			break
		}

	}

	var n []string = tables[0:i]
	return &n
}