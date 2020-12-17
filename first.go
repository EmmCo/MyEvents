package main

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"time"
    mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
type DatabaseHandler interface {
	AddEvent(Event) ([] byte, error)
	FindEvent([]byte) (Event, error)
	FindEventByName(string) (Event, error)
	FindAllAvailableEvents()( []Event, error)
}


func (mgolayer *MongoDBLayer) getFreshSession() *mgo.Session {
	return mgolayer.session.Copy()
}

func (mgolayer *MongoDBLayer) FindEvent(id []byte) (Event, error)  {
	s := mgolayer.getFreshSession()
	defer s.Close()
	e := Event{}
	err := s.DB(DB).C(EVENTS).FindId(bson.ObjectId(id)).One(&e)
	return e, err
}

func (mgolayer *MongoDBLayer) FindAllAvailableEvents() ([]Event, error)  {
	s := mgolayer.getFreshSession()
	defer s.Close()
	es := []Event{}
	err := s.DB(DB).C(EVENTS).Find(nil ).All(&es)
	return es, err
}

func (mgolayer *MongoDBLayer) FindEventByName(name string) (Event, error)  {
	s := mgolayer.getFreshSession()
	defer s.Close()
	e := Event{}
	err := s.DB(DB).C(EVENTS).Find(bson.M{"name": name }).One(&e)
	return e, err
}

func (mgolayer *MongoDBLayer) AddEvent(e Event) ([]byte, error)  {
	s := mgolayer.getFreshSession()
	defer s.Close()
	if !e.ID.Valid() {
		e.ID = bson.NewObjectId()
	}
	if !e.Location.ID.Valid() {
		e.Location.ID = bson.NewObjectId()
	}
	return []byte(e.ID), s.DB(DB).C(EVENTS).Insert(e)
}


type Event struct {
	ID bson.ObjectId `bson:"_id"`
	Name string
	Duration int
	StartDate int64
	EndDate int64
	Location Location
}

type Location struct {
	ID bson.ObjectId `bson:"_id"`
	Name string
	Address string
	Country string
	OpenTime int
	CloseTime int
	Halls []Hall
}

type Hall struct {
	Name string `json:"name"`
	Location string `json:location,omitempty`
	Capacity int `json:"capacity"`
}

type eventServiceHandler struct {
	dbhandler DatabaseHandler
}

func newEventHandler(databasehandler DatabaseHandler) *eventServiceHandler {
	return &eventServiceHandler {
		dbhandler: databasehandler,
	}
}

func (eh * eventServiceHandler) findEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	criteria, ok := vars["SearchCriteria"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, `{error: No Search criteria found}`)
		return
	}
	searchkey, ok := vars["Search"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, `{error: No Search keys found}`)
		return
	}
	var event Event
	var err error
	switch strings.ToLower(criteria) {
	case "name":
		event, err = eh.dbhandler.FindEventByName(searchkey)
	case "id":
		id, err := hex.DecodeString(searchkey)
		if err == nil {
			event, err = eh.dbhandler.FindEvent(id)
		}
	}
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "{error occured %s}", err)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	json.NewEncoder(w).Encode(&event)
}

func (eh * eventServiceHandler) allEventHandler(w http.ResponseWriter, r *http.Request) {
	events, err := eh.dbhandler.FindAllAvailableEvents()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{error: %s}", err)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	json.NewEncoder(w).Encode(&events)
}

func (eh * eventServiceHandler) newEventHandler(w http.ResponseWriter, r *http.Request) {
	event := Event{}
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{error : %s}", err)
		return
	}
	id, err := eh.dbhandler.AddEvent(event)
	if nil != err {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{id: %d, error: %s}", id, err)
		return
	}
}

const (
	DB = "myevents"
	USERS = "users"
	EVENTS = "events"
)
type MongoDBLayer struct {
	session *mgo.Session
}

func NewMongoDBLayer(connection string) (*MongoDBLayer, error) {
	s, err := mgo.Dial(connection)
	if err != nil {
		return nil, err
	}
	return &MongoDBLayer {
		session: s,
	}, nil
}

func ServeAPI(endpoint string, dbHandler DatabaseHandler) error {
	handler := newEventHandler(dbHandler)
	rtr := mux.NewRouter()
	rtr.HandleFunc("/pages/{id:[0-9]+}", pageHandler)
	rtr.HandleFunc("/process", process)
	eventsrouter := rtr.PathPrefix("/events").Subrouter()
	eventsrouter.Methods("GET").
		Path("/{SearchCriteria}/{search}").HandlerFunc(handler.findEventHandler)
	eventsrouter.Methods("GET").Path("").HandlerFunc(handler.allEventHandler)
	eventsrouter.Methods("POST").Path("").HandlerFunc(handler.allEventHandler)
}

func main() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/pages/{id:[0-9]+}", pageHandler)
	rtr.HandleFunc("/process", process)
	handler := &eventServiceHandler{}
	eventsrouter := rtr.PathPrefix("/events").Subrouter()
	eventsrouter.Methods("GET").
		Path("/{SearchCriteria}/{search}").HandlerFunc(handler.findEventHandler)
	eventsrouter.Methods("GET").Path("").HandlerFunc(handler.allEventHandler)
	eventsrouter.Methods("POST").Path("").HandlerFunc(handler.allEventHandler)

	//http.Handle("/", rtr)
	http.ListenAndServe(PORT, rtr)
	http.HandleFunc("/static",serveStatic)
	http.HandleFunc("/",serveDynamic)
	http.ListenAndServe(PORT,nil)
	gintest()
	mysqltest()
	redistest()
}
