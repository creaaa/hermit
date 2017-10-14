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

	"log"

	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	fmt.Println("init!!")
	var err error

	db, err = sql.Open("sqlite3", "./data.db")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func db_exec(db *sql.DB, q string) {
	var _, err = db.Exec(q)
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

func create(body string) {

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
	/*res*/
	_, err := db.Exec(
		`INSERT INTO memo (body) VALUES (?)`,
		body,
	)
	if err != nil {
		panic(err)
	}

	// 挿入処理の結果からIDを取得
	//id, err := res.LastInsertId()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(id)
}

func readAll() {
	// 複数レコード取得
	rows, err := db.Query(
		`SELECT * FROM memo`,
	)
	if err != nil {
		panic(err)
	}

	// 処理が終わったらカーソルを閉じる
	defer rows.Close()
	for rows.Next() {
		var id int
		var body string
		var created time.Time

		// カーソルから値を取得
		// ...なんかこう、C言語チックな「副作用前提の」コードバリバリ使うんやな。
		// これあんま好きじゃねぇな...
		// たった1節で、エラー処理とエラーなし時の処理を同時に書けるのがメリットなんだろう。
		// これをイヤと思うのはSwift脳だからだろうか...

		// このscanの中、定義したカラム文すべて引数取らないとエラーになる、回避策あるだろ
		if err := rows.Scan(&id, &body, &created); err != nil {
			log.Fatal("rows.Scan()", err)
			return
		}
		fmt.Printf("id: %d, title: %s, created: %v\n", id, body, created)
	}
}

func read(id int) {
	// 1件取得
	row := db.QueryRow(
		`SELECT * FROM memo WHERE ID=?`,
		id,
	)

	// var id int
	var body string
	var created time.Time

	err := row.Scan(&id, &body, &created)

	// エラー内容による分岐
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("Not found")
	case err != nil:
		panic(err)
	default:
		fmt.Printf("id: %d, title: %s, created: %v\n", id, body, created)
	}
}

func update(newValue string, id int) {
	// 更新
	res, err := db.Exec(
		`UPDATE memo SET body=? WHERE id=?`,
		newValue,
		id,
	)
	if err != nil {
		panic(err)
	}

	// 更新されたレコード数
	affect, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}

	fmt.Printf("affected by update: %d\n", affect)
}

func delete(id int) {
	// 削除
	res, err := db.Exec(
		`DELETE FROM memo WHERE ID=?`,
		id,
	)
	if err != nil {
		panic(err)
	}

	// 削除されたレコード数
	affect, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}

	fmt.Printf("affected by delete: %d\n", affect)
}

func main() {

	// setup()
	// create()
	// readAll()
	// read(5)
	// update("月曜おはす〜", 3)

	//delete(4)
	//readAll()
	//
	//create()
	//readAll()

	//create("マリリ")
	//readAll()

	//db.Close()
}
