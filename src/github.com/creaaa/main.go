/*
DBは常にオープンしてないと、だめだ！！ランタイムエラーになる！

くそばまりポイント発見！
db, err = sql.Open("sqlite3", "./unko.db")

これ、通った「時点」でファイル作られると思うが違う！
実際は、テーブル定義を db.Exec(q) しないと作られない！
これ気付かず、単独で Open() だけしてると、一生ファイル作られないからマジ気をつけろ！！

これ、大いに役立った
http://kuroeveryday.blogspot.ca/2017/08/sqlite3-with-golang.html

ハマりポイント！！！
当たり前で、早く気付けよって感じだが、
current directory の出力: os.Getwd() は、ターミナルのカレントディレクトリの影響を受ける！！
「ファイルがどこに位置しているか」は関係ない！！
これは、プログラム内でファイルを相対パスで指定しているとき、バリバリ影響を受けるってこと！！
ターミナルから実行するとき、これは注意だ！
てかそれなら、たいていの場合、絶対パス指定のほうが望ましい...望ましいよね?

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

	"os/exec"

	"bufio"

	"net/http"

	"sort"

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

	q = "CREATE TABLE urls ("
	q += " id INTEGER PRIMARY KEY"
	q += ", alias VARCHAR(32) NOT NULL"
	q += ", desc VARCHAR(255)"
	q += ", url VARCHAR(255) NOT NULL UNIQUE"
	q += ", flag INTEGER"
	q += ")"

	db_exec(db, q)
}

// これ、スライス(参照型)渡してるから元来破壊的かと思ったが、
// なぜか ちゃんと結果を ([]string) 返さないとだめだった。なんで..
func argParse(args []string) []string {
	for idx, arg := range os.Args {
		if idx == 0 || idx == 1 {
			continue
		}

		if strings.HasPrefix(arg, "http") || strings.HasPrefix(arg, "https") {
			fmt.Println("this is url!!")
			args = append(args, arg)
		} else if i, err := strconv.Atoi(arg); err == nil {
			// intがまざってたら、URLに変換する処理を書く
			fmt.Println("int!!!")
			// DBから、IDをキーにURLを取得
			if url := readURL(i); url != "" {
				args = append(args, url)
				fmt.Println("おら", args)
			}
		} else {
			fmt.Println("エイリアスの可能性!")
			// エイリアスならURLに変換する処理を書く
			if url := readURL(arg); url != "" {
				args = append(args, url)
			}
		}
	}
	fmt.Println("最終結果: ", args)
	return args
}

// org -o
func openURL() {
	// current directory
	//dir, _ := os.Getwd()
	//fmt.Println(dir)

	// 外部コマンドの結果をターミナルに出したいなら、こうしてわざわざ変数に入れないといけない
	// res, _ := exec.Command("ls", "-la").Output()
	// fmt.Printf("%s", res)

	urls := []string{
		"-a", "Google Chrome", "-n",
		"--args", "--incognito",
	}

	urls = argParse(urls)
	fmt.Println("ks", urls)
	exec.Command("open", urls...).Run()
}

// org -a
func add() {

	args := os.Args

	// URLバリデーション
	if !strings.HasPrefix(args[2], "http") && !strings.HasPrefix(args[2], "https") {
		panic("invalid URL")
	}

	var (
		id    = getMinimumID()
		url   = args[2]
		alias = args[3]
		desc  string
		flag  int = 0
	)

	if len(args) >= 5 {
		desc = args[4]
	}

	_, err := db.Exec(
		`INSERT INTO urls (id, alias, desc, url, flag) VALUES (?, ?, ?, ?, ?)`,
		id, alias, desc, url, flag,
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
		`SELECT id, alias, desc, url, flag FROM urls`,
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
			url   string
			flag  int
		)

		// カーソルから値を取得
		// ...なんかこう、C言語チックな「副作用前提の」コードバリバリ使うんやな。
		// これあんま好きじゃねぇな...
		// たった1節で、エラー処理とエラーなし時の処理を同時に書けるのがメリットなんだろう。
		// これをイヤと思うのはSwift脳だからだろうか...

		// このscanの中、定義したカラム文すべて引数取らないとエラーになる、回避策あるだろ
		if err := rows.Scan(&id, &alias, &desc, &url, &flag); err != nil {
			log.Fatal("rows.Scan()", err)
			return
		}
		fmt.Printf("id: %d, alias: %s, desc: %s, url: %s, flag: %d\n", id, alias, desc, url, flag)
	}
}

// org -f
func fetch() {
	urls := getAllURLs()
	// 404を返したレコードに対しフラグをオン
	count := isResourceExist(urls)

	fmt.Printf("affected by update: %d\n", count)
}

// 全レコードのURLを抽出して返す
func getAllURLs() []string {
	var urls []string
	// 複数レコード取得
	rows, err := db.Query(
		`SELECT url FROM urls`,
	)
	if err != nil {
		panic(err)
	}
	// 処理が終わったらカーソルを閉じる
	defer rows.Close()
	for rows.Next() {
		var url string
		// このscanの中、定義したカラム文すべて引数取らないとエラーになる、回避策あるだろ
		if err := rows.Scan(&url); err != nil {
			panic(err)
		}
		urls = append(urls, url)
	}
	return urls
}

// 全リソースに対しHEADメソッドを投げて404かチェック。
// 404だったURLのレコード数を返す
func isResourceExist(urls []string) int {

	var res *http.Response
	count := 0

	for _, url := range urls {
		res, _ = http.Head(url)
		fmt.Println("status code: ", res.StatusCode)
		// もし404かつフラグが立ってないならば、flagをon
		if res.StatusCode == 404 && getFlag(url) == 0 {
			fmt.Println("404!!")
			updateFlag(url)
			count += 1
		}
	}
	return count
}

func getFlag(url string) (flag int) {

	row := db.QueryRow(`SELECT flag FROM urls WHERE url=?`, url)
	err := row.Scan(&flag)

	if err != nil {
		panic(err)
	}
	return
}

func updateFlag(targetURL string) {
	_, err := db.Exec(`UPDATE urls SET flag=1 WHERE url=?`, targetURL)
	if err != nil {
		panic(err)
	}
}

// org -d
func delete() {
	args := argParse([]string{})

	var (
		res sql.Result
		err error
	)

	if strings.HasPrefix(args[0], "http") || strings.HasPrefix(args[0], "https") {
		fmt.Println("this is url!!")
		res, err = db.Exec(
			`DELETE FROM urls WHERE url=?`,
			args[0], // 1個しか消せないようにした
		)
	} else if _, err := strconv.Atoi(args[0]); err == nil {
		// intがまざってたら、URLに変換する処理を書く
		fmt.Println("int!!!")
		res, err = db.Exec(
			`DELETE FROM urls WHERE ID=?`,
			args[0], // 1個しか消せないようにした
		)
	} else {
		fmt.Println("エイリアスの可能性!")
		// エイリアスならURLに変換する処理を書く
		res, err = db.Exec(
			`DELETE FROM urls WHERE alias=?`,
			args[0], // 1個しか消せないようにした
		)
	}

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

func deleteOnlyFlagOn() {
	res, err := db.Exec(`DELETE FROM urls WHERE flag=1`)
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

// org -da
func deleteAll() {

	var (
		res sql.Result
		err error
	)

	fmt.Println("Do you really delete all records? [Y/n]")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Text() == "y" || scanner.Text() == "Y" {
		res, err = db.Exec(`DELETE FROM urls`)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("did nothing...")
	}

	// 削除されたレコード数
	affect, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}

	fmt.Printf("affected by delete: %d\n", affect)
}

//////////////////////////////////////////////////

func create(body string) {

	_, err := db.Exec(
		`INSERT INTO memo (body) VALUES (?)`,
		body,
	)
	if err != nil {
		panic(err)
	}
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
		fmt.Println("ここってことやな！", key)
		// エイリアスが同じでも、QueryRowなら最初の1行だけ返すだけ!
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

func parse() {

	if !isEqualOrGreaterThanMinArgs(2) {
		fmt.Println("Invalid argument. exit...")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "-o":
		fmt.Println("openURL!")
		openURL()
	case "add", "-a", "--add":
		fmt.Println("add!")
		if isEqualOrGreaterThanMinArgs(4) {
			add()
		} else {
			panic("invalid argument: you need to add at least URL & alias")
		}
	//case "-n":
	//	fmt.Println("openURL!")
	//case "-e":
	//	fmt.Println("openURL!")
	//case "-u":
	//	fmt.Println("openURL!")
	case "delete", "-d", "--delete":
		if isEqualOrGreaterThanMinArgs(3) {
			// -f をつけると、flagがonのものだけ消す
			if os.Args[2] == "-f" || os.Args[2] == "--flag" {
				fmt.Println("delete only flag!")
				deleteOnlyFlagOn()
			} else {
				delete()
			}
		} else {
			panic("invalid argument: you need to add at least URL or alias or ID")
		}
	case "deleteall", "-da", "--deleteall":
		fmt.Println("deleteAll!")
		deleteAll()
	case "list", "-l", "--list":
		fmt.Println("list!")
		list()
	case "fetch", "-f", "--fetch":
		fmt.Println("fetch!!")
		fetch()
	default:
		fmt.Println("no such command. exit...")
		os.Exit(1)
	}
}

func isEqualOrGreaterThanMinArgs(minimum int) bool {
	if len(os.Args) >= minimum {
		return true
	}
	return false
}

// 空いているIDの最小値を求める
func getMinimumID() int {

	rows, err := db.Query(`SELECT id FROM urls`)
	if err != nil {
		panic(err)
	}

	ids := []int{}

	for rows.Next() {
		var id int
		rows.Scan(&id)
		ids = append(ids, id)
	}

	sort.Ints(ids)

	fmt.Println("ソート済: ", ids)

	// ソート完了したので。。
	return subRoutine(ids, 1)
}

func subRoutine(ids []int, inspector int) int {
	fmt.Println("調査開始！", ids)
	for _, id := range ids {
		if id == inspector {
			// あったのでまだ調査
			fmt.Println("あったのでまだ調査: 次は", ids, inspector+1)
			return subRoutine(ids, inspector+1)
		} else {
			fmt.Println("現在のID: ", id, "調査対象: ", inspector)
			fmt.Println("違った")
		}
	}
	// ないので終了
	fmt.Println("ないので終了")
	return inspector
}

func main() {
	//setup()
	parse()
	//db.Close()

	fmt.Println(getMinimumID())
}
