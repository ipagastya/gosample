package hello

import (
  "net/http"
  "log"
  "fmt"
  "expvar"

  _"github.com/lib/pq"
  logging "gopkg.in/tokopedia/logging.v1"
  "github.com/tokopedia/sqlt"
)

type ServerConfig struct {
  Name string
}
type DatabaseConfig struct {
  Type       string
  Connection string
}

type Config struct {
  Server   ServerConfig
  Database DatabaseConfig
}

type HelloWorldModule struct {
  cfg   *Config
  DB    *sqlt.DB
  something string
  stats   *expvar.Int
}

func NewHelloWorldModule() *HelloWorldModule {

  var cfg Config

  ok := logging.ReadModuleConfig(&cfg, "config", "hello") || logging.ReadModuleConfig(&cfg, "files/etc/gosample", "hello")
  if !ok {
    // when the app is run with -e switch, this message will automatically be redirected to the log file specified
    log.Fatalln("failed to read config")
  }

  masterDB := cfg.Database.Connection
  slaveDB := cfg.Database.Connection
  dbConnection := fmt.Sprintf("%s;%s", masterDB, slaveDB)

  db, err := sqlt.Open(cfg.Database.Type, dbConnection)
  if err != nil {
  	log.Fatalln("Failed to connect database. Error: ", err.Error())
  }

  // this message only shows up if app is run with -debug option, so its great for debugging
  logging.Debug.Println("hello init called",cfg.Server.Name)

  return &HelloWorldModule{
    cfg: &cfg,
    DB: db,
    something: "John Doe",
    stats : expvar.NewInt("rpsStats"),
  }
}

func (hlm *HelloWorldModule) SayHelloWorld(w http.ResponseWriter, r *http.Request) {
  hlm.stats.Add(1)
  w.Write([]byte("Hello " + hlm.something))
}

//FOR TRAINING
/*
CREATE TABLE "public"."zaki_test" (
	"id" int4 NOT NULL,
	"name" varchar(255) COLLATE "default",
	CONSTRAINT "zaki-test_pkey" PRIMARY KEY ("id") NOT DEFERRABLE INITIALLY IMMEDIATE
)
WITH (OIDS=FALSE);
ALTER TABLE "public"."zaki_test" OWNER TO "tokopedia";
*/
type ZakiTest struct {
	ID       int64 
	FullName string `db:"name"`
}

func (hlm *HelloWorldModule) InsertDataIntoDatabase(w http.ResponseWriter, r *http.Request) {
  hlm.stats.Add(1)

  id := r.FormValue("id")
  name := r.FormValue("name")
  result := "Insert Success"

  query := "INSERT INTO zaki_test(id, name) VALUES($1, $2)"
  _, err := hlm.DB.Exec(query, id, name)
  if err != nil {
  	log.Println("Error Insert to Database. Error: ", err.Error())
  	result = "Insert Failed"
  }

  w.Write([]byte(result))
}

func (hlm *HelloWorldModule) UpdateDataIntoDatabase(w http.ResponseWriter, r *http.Request) {
  hlm.stats.Add(1)

  id := r.FormValue("id")
  name := r.FormValue("name")
  result := "Update Success"

  query := "UPDATE zaki_test SET name = $1 WHERE id = $2"
  _, err := hlm.DB.Exec(query, name, id)
  if err != nil {
  	log.Println("Error Update to Database. Error: ", err.Error())
  	result = "Update Failed"
  }

  w.Write([]byte(result))
}