package main

import (
	"database/sql"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"
	"time"

	//"github.com/jmoiron/sqlx"
)
func gintest() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, Geektutu, my name is Zhy" +
			"<head> --- </head>")
	})
	r.GET("/home", func(c *gin.Context) {
		c.String(http.StatusOK, "this is the home")
	})
	r.GET("/home/room", func(c *gin.Context) {
		c.String(http.StatusOK, "this is the room")
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
func redistest () {
	conn, err := redis.Dial("tcp", "10.64.33.42:6379")
	if err != nil {
		fmt.Println("connect redis error :", err)
		return
	}
	defer conn.Close()
	name, err := conn.Do("GET", "abc")
	if err != nil {
		fmt.Println("redis set error :", err)
	} else {
		fmt.Printf("Got name : %s \n", name)
	}
	name, err = redis.String(conn.Do("SET", "abc", "woyeainio"))
	if err != nil {
		fmt.Println("redis set error :", err)
	} else {
		fmt.Printf("Set name : %s \n", name)
	}
}
func mysqltest() (err error) {
	type user struct {
		id   int
		age  int
		name string
	}
	dsn := "root:woshishui@tcp(10.64.33.42:3306)/test_db"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()  // 注意这行代码要写在上面err判断的下面
	// 尝试与数据库建立连接（校验dsn是否正确）
	err = db.Ping()
	if err != nil {
		fmt.Printf("Ping failed, err:%v\n", err)
		return nil
	}
	fmt.Println("db.Ping Successed")
	sqlStr := "insert into user(name, age) values (?,?)"
	ret, err := db.Exec(sqlStr, "王五", 38)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	theID, err := ret.LastInsertId() // 新插入数据的id
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return
	}
	fmt.Printf("insert success, the id is %d.\n", theID)
	sqlStr = "select id, name, age from user where id=?"
	var u user
	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	err = db.QueryRow(sqlStr, 1).Scan(&u.id, &u.name, &u.age)
	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return
	}
	fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)
	return nil
}

const (
	PORT = ":8080"
)
func serveStatic(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static.html")
}
func serveDynamic(w http.ResponseWriter, r *http.Request) {
	response := "The time is now " + time.Now().String()
	fmt.Fprintln(w,response)
}
func pageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageID := vars["id"]
	fileName := "files/" + pageID + ".html"
	http.ServeFile(w,r,fileName)
}

func process(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Fprintln(w, r.PostForm)
	fmt.Fprintln(w, r.PostFormValue("hello"))
}
func main() {

	rtr := mux.NewRouter()
	rtr.HandleFunc("/pages/{id:[0-9]+}", pageHandler)
	rtr.HandleFunc("/process", process)
	//http.Handle("/", rtr)
	http.ListenAndServe(PORT, rtr)
	http.HandleFunc("/static",serveStatic)
	http.HandleFunc("/",serveDynamic)
	http.ListenAndServe(PORT,nil)
	gintest()
	mysqltest()
	redistest()
}
