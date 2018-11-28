package mysql

import (
	"fmt"
	"testing"
)

func Test_test(t *testing.T) {

	config := &Config{
		NodeName: "default",
		Host:     "localhost",
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
	tx, _ := dao.Begin()
	tx.Delete("delete from root where id > 1")
	tx.Rollback()
	ress, _ := dao.Select("select count(*) as count from root")
	fmt.Println(ress)

	tx, _ = dao.Begin()
	tx.Delete("delete from root where id > 1")
	tx.Commit()
	resss, _ := dao.Select("select count(*) as count from root")
	fmt.Println(resss)
}
