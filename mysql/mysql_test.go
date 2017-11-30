package mysql

import (
	"fmt"
	"testing"
)

func Test_test(t *testing.T) {

	config := &Config{
		NodeName: "default",
		Host:     "139.129.36.196",
		Port:     3306,
		Database: "gameinfo",
		UserName: "test",
		Password: "123qwe",
		Charset:  "utf8",
	}

	dao, _ := GetDB(config)

	res, _ := dao.Select("select * from root where id = ?", 1)
	fmt.Println(res)

	rows, _ := dao.Delete("delete from root where id = ?", 6)
	fmt.Println(rows)

	rowss, _ := dao.Update("update root set auth = '*'")
	fmt.Println(rowss)

	lastID, _ := dao.Insert("insert into root (username, password, auth) values ('admin', 'admin', 'xxx')")
	fmt.Println(lastID)

	//trans
	dao.Begin()
	dao.Delete("delete from root where id > 1")
	dao.Rockback()
	ress, _ := dao.Select("select count(*) as count from root")
	fmt.Println(ress)

	dao.Begin()
	dao.Delete("delete from root where id > 1")
	dao.Commit()
	resss, _ := dao.Select("select count(*) as count from root")
	fmt.Println(resss)
}