package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" //mysql dirver
	"reflect"
	"testing"
)

func Test_test(t *testing.T) {
	config := &Config{
		NodeName: "default",
		Host:     "127.0.0.1",
		Port:     3306,
		Database: "web",
		UserName: "root",
		Password: "",
		Charset:  "utf8",
	}

	dao, _ := GetDB(config)

	//res, _ := dao.Select("select * from contract limit 1")
	//fmt.Println(res)

	sqls := []string{"select * from web_mine_booking where id = 3", "select * from web_mine_booking where id = 2"}
	ret, err := dao.MultiSelect(sqls)
	fmt.Println(ret, err)

	fmt.Println(reflect.TypeOf(ret[0][0]["id"]))

	r, _ := dao.Select("select * from web_mine_booking where id = 3")
	fmt.Println(reflect.TypeOf(r[0]["id"]))

	//rows, _ := dao.Delete("delete from root where id = ?", 6)
	//fmt.Println(rows)

	//rowss, _ := dao.Update("update root set auth = '*'")
	//fmt.Println(rowss)
	//
	//lastID, _ := dao.Insert("insert into root (username, password, auth) values ('admin', 'admin', 'xxx')")
	//fmt.Println(lastID)
	//
	////trans
	//tx, _ := dao.Begin()
	//tx.Delete("delete from root where id > 1")
	//tx.Rollback()
	//ress, _ := dao.Select("select count(*) as count from root")
	//fmt.Println(ress)
	//
	//tx, _ = dao.Begin()
	//tx.Delete("delete from root where id > 1")
	//tx.Commit()
	//resss, _ := dao.Select("select count(*) as count from root")
	//fmt.Println(resss)
}
