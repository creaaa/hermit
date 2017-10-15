/*
DBは常にオープンしてないと、だめだ！！ランタイムエラーになる！

くそばまりポイント発見！
db, err = sql.Open("sqlite3", "./unko.db")

これ、通った「時点」でファイル作られると思うが違う！
実際は、テーブル定義を db.Exec(q) しないと作られない！
これ気付かず、単独で Open() だけしてると、一生ファイル作られないからマジ気をつけろ！！

これ、大いに役立った
http://kuroeveryday.blogspot.ca/2017/08/sqlite3-with-golang.html

*/

package main

import (
	"database/sql"
	"fmt"
	"os"

	"log"

	"time"

	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	fmt.Println("init!!")

	//// current directory
	dir, err := os.Getwd()
	fmt.Println(dir)

	//var err error

	db, err = sql.Open("sqlite3", "../../../data.db")

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

	//q = "CREATE TABLE memo ("
	//q += " id INTEGER PRIMARY KEY AUTOINCREMENT"
	//q += ", body VARCHAR(255) NOT NULL"
	//q += ", created_at TIMESTAMP DEFAULT (DATETIME('now','localtime'))"
	//q += ")"

	q = "CREATE TABLE urls ("
	q += " id INTEGER PRIMARY KEY AUTOINCREMENT"
	q += ", alias VARCHAR(32) NOT NULL"
	q += ", desc VARCHAR(255)"
	q += ", url VARCHAR(255) NOT NULL"
	q += ", flag INTEGER"
	q += ")"

	db_exec(db, q)
}

// org -o
func openURL() {
	// current directory
	//dir, _ := os.Getwd()
	//fmt.Println(dir)

	// 外部コマンドの結果をターミナルに出したいなら、こうしてわざわざ変数に入れないといけない
	// res, _ := exec.Command("ls", "-la").Output()
	// fmt.Printf("%s", res)

	urls := []string{}

	for idx, arg := range os.Args {
		if idx == 0 {
			continue
		}

		if strings.HasPrefix(arg, "http") || strings.HasPrefix(arg, "https") {
			fmt.Println("this is url!!")
			urls = append(urls, arg)
		} else if i, err := strconv.Atoi(arg); err == nil {
			// intがまざってたら、URLに変換する処理を書く
			fmt.Println("int!!!")
			// DBから、IDをキーにURLを取得
			if url := readURL(i); url != "" {
				urls = append(urls, url)
			}
		} else {
			fmt.Println("エイリアスの可能性!")
			// エイリアスならURLに変換する処理を書く
			if url := readURL(arg); url != "" {
				urls = append(urls, url)
			}
		}
	}

	fmt.Println(urls)

	//exec.Command("open", "-a", "Google Chrome", "-n",
	//	"--args", "--incognito", "http://www.yahoo.co.jp", "https://www.google.ca/").Run()

}

// org -a
func add(url, alias, desc string) {

	if url == "" {
		panic("no url!")
	} else if alias == "" {
		panic("no alias!")
	}

	var (
		flag int = 0
	)

	_, err := db.Exec(
		`INSERT INTO urls (alias, desc, url, flag) VALUES (?, ?, ?, ?)`,
		alias, desc, url, flag,
	)
	if err != nil {
		panic(err)
	}
}

// org -n
func name() {
}

// org -e
func explain() {
}

// org -l
func list() {
	// 複数レコード取得
	rows, err := db.Query(
		`SELECT id, alias, desc, flag FROM urls`,
	)
	if err != nil {
		panic(err)
	}

	// 処理が終わったらカーソルを閉じる
	defer rows.Close()
	for rows.Next() {
		var (
			id    int
			alias string
			desc  string
			flag  int
		)

		// カーソルから値を取得
		// ...なんかこう、C言語チックな「副作用前提の」コードバリバリ使うんやな。
		// これあんま好きじゃねぇな...
		// たった1節で、エラー処理とエラーなし時の処理を同時に書けるのがメリットなんだろう。
		// これをイヤと思うのはSwift脳だからだろうか...

		// このscanの中、定義したカラム文すべて引数取らないとエラーになる、回避策あるだろ
		if err := rows.Scan(&id, &alias, &desc, &flag); err != nil {
			log.Fatal("rows.Scan()", err)
			return
		}
		fmt.Printf("id: %d, alias: %s, desc: %s, flag: %d\n", id, alias, desc, flag)
	}
}

// org -f
func fetch() {
}

// org -u
func update() {
}

// org -d
func delete() {
}

// org -da
func deleteAll() {
}

//////////////////////////////////////////////////

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

func readURL(key interface{}) string {
	var row *sql.Row

	switch key.(type) {
	case int:
		row = db.QueryRow(`SELECT url FROM urls WHERE ID=?`, key)
	case string:
		row = db.QueryRow(`SELECT url FROM urls WHERE ALIAS=?`, key)
	}

	var url string

	err := row.Scan(&url)

	// エラー内容による分岐
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("Not found")
	case err != nil:
		panic(err)
	}

	return url // 該当行がなかった場合は、""(空文字)が返る点に注意
}

//func update(newValue string, id int) {
//	// 更新
//	res, err := db.Exec(
//		`UPDATE memo SET body=? WHERE id=?`,
//		newValue,
//		id,
//	)
//	if err != nil {
//		panic(err)
//	}
//
//	// 更新されたレコード数
//	affect, err := res.RowsAffected()
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Printf("affected by update: %d\n", affect)
//}
//
//func delete(id int) {
//	// 削除
//	res, err := db.Exec(
//		`DELETE FROM memo WHERE ID=?`,
//		id,
//	)
//	if err != nil {
//		panic(err)
//	}
//
//	// 削除されたレコード数
//	affect, err := res.RowsAffected()
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Printf("affected by delete: %d\n", affect)
//}

func main() {

	//setup()
	//add("https://www.youtube.com/watch?v=87wf45zW5NA",
	//	"sayoemo",
	//	"さよならはエモーション")

	//list()

	openURL()

	//create()
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
