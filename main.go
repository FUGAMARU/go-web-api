package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"io/ioutil"
)

func main() {
	//JSONファイルからDBサーバーの接続情報を取ってくる
	JSONFile, err := ioutil.ReadFile("./mysql-info.json")
	if err != nil {
		fmt.Println("JSONファイル読み込みエラー")
		log.Fatal(err)
	}

	err = json.Unmarshal(JSONFile, &mysql_info)
	if err != nil {
		fmt.Println("JSON変換エラー")
		log.Fatal(err)
	}
	fmt.Println("SQLサーバーの接続先情報を読み込みました")

	echo := echo.New()

	echo.Use(middleware.Logger())
	echo.Use(middleware.Recover())
	echo.Use(middleware.CORS())

	echo.GET("/get-json", getJSON)

	echo.Start(":9940")
}

var mysql_info MySQL_Info

type MySQL_Info struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type VoiceActor struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Office      string `json:"office"`
	Profile_URL string `json:"profile_url"`
}

func getJSON(c echo.Context) error {
	//DBサーバー接続
	db, err := sql.Open("mysql", mysql_info.User+":"+mysql_info.Password+"@tcp("+mysql_info.Host+":"+mysql_info.Port+")/"+mysql_info.Database)
	if err != nil {
		fmt.Println("MySQLサーバー 接続エラー")
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM voice_actors")
	if err != nil {
		fmt.Println("SQL実行エラー")
		log.Fatal(err)
	}
	defer rows.Close()

	var actors []VoiceActor //データーをまとめてJSONにするための配列
	for rows.Next() {
		actor := VoiceActor{}
		rows.Scan(&actor.Id, &actor.Name, &actor.Office, &actor.Profile_URL)
		fmt.Println(actor)
		actors = append(actors, actor)
	}

	res, err := json.Marshal(actors)
	if err != nil {
		fmt.Print("JSON変換エラー")
		log.Fatal(err)
	}

	return c.String(http.StatusOK, string(res))
}
