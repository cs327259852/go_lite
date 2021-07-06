package main

import (
	"flag"
	"study/etl_validate/common"
)

func main()(){
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
		[]string{"version", "branchid", "storeid", "prodid", "quantity","lastmodifytime","deleteflag",}}
	common.CommonCompare(&dir,&tables,&fields,&compareFields)
}