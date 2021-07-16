package main

import (
	"flag"
	"fmt"
	"study/etl_validate/common"
	"time"
)

func main() {
	var dir string
	flag.StringVar(&dir, "d", "目录名", "目录名")
	flag.Parse()
	now := time.Now()
	//对比的表名
	var tables []string = []string{"tb_gos_stock_stockpreemption", "vw_common_prod"}
	//对比的表列名
	var fields [][]string = [][]string{
		[]string{"pk", "fk", "createtime", "lastmodifytime", "version", "lineid", "branchid", "deleteflag", "note", "preemptionpreemption", "prodid", "lotno", "quantity", "rowguid", "billid", "whseid", "storeid", "billguid", "opid", "custid", "custno", "custname"},
		[]string{"pk", "branchid", "prodid", "prodno", "deleteflag", "version"}}
	//需要对比的表列名 对比的列名才会打印出来
	var compareFields [][]string = [][]string{
		[]string{"version", "branchid", "storeid", "prodid", "quantity", "lastmodifytime", "deleteflag"},
		[]string{"branchid", "prodid", "prodno", "deleteflag", "version"}}
	common.CommonCompare(&dir, &tables, &fields, &compareFields)
	fmt.Printf("完成对比，耗时:%v\n", time.Since(now))
}
