package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func main() {
	var dir string
	flag.StringVar(&dir, "d", "目录名", "目录名")
	flag.Parse()
	now := time.Now()
	//对比的库源和目标库前缀
	var srouceDestPair [][]string = [][]string{[]string{"erp_", "mid_"},
	}
	//dir="/home/peter/tmp"

	//对比的表名
	var tables *[]string = getAllTables("/home/peter/go/src/study/etl_validate/tmpconpare/tables.txt")
	var fields []string = []string{"pk","lastmodifytime"}
	var compareFields []string = []string{"lastmodifytime"}
	commonCompare(&dir, &srouceDestPair, tables, fields, compareFields)
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