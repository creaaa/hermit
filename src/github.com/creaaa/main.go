/*
DBは常にオープンしてないと、だめだ！！ランタイムエラーになる！

くそばまりポイント発見！
db, err = sql.Open("sqlite3", "./unko.db")

これ、通った「時点」でファイル作られると思うが違う！
実際は、テーブル定義を db.Exec(q) しないと作られない！
これ気付かず、単独で Open() だけしてると、一生ファイル作られないからマジ気をつけろ！！

*/

package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	fmt.Println("はいうんこ")
	var err error

	db, err = sql.Open("sqlite3", "./data.db")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setup() {

	var q = ""

	q = "CREATE TABLE memo ("
	q += " id INTEGER PRIMARY KEY AUTOINCREMENT"
	q += ", body VARCHAR(255) NOT NULL"
	q += ", created_at TIMESTAMP DEFAULT (DATETIME('now','localtime'))"
	q += ")"
	db_exec(db, q)
}

func create() {

	// まごころこめるとこう

	//var q = ""
	//
	//q = "INSERT INTO memo "
	//q += " (body)"
	//q += " VALUES"
	//q += " ('body1')"
	//q += ",('body2')"
	//q += ",('body3')"
	//db_exec(db, q)

	// ライブラリの恩恵に授かるとこう
	res, err := db.Exec(
		`INSERT INTO memo (body) VALUES (?)`,
		"body4",
	)
	if err != nil {
		panic(err)
	}

	// 挿入処理の結果からIDを取得
	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	fmt.Println(id)

}

func db_exec(db *sql.DB, q string) {
	var _, err = db.Exec(q)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {

	// setup()
	create()

	//db.Close()
}
