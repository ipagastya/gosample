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

func (hlm *HelloWorldModule) GetSingleDataFromDatabase(w http.ResponseWriter, r *http.Request) {
  hlm.stats.Add(1)

  test := ZakiTest{}
  query := "SELECT id, name FROM zaki_test LIMIT 1"
  err := hlm.DB.Get(&test, query)
  if err != nil {
  	log.Println("Error Query Database. Error: ", err.Error())
  }

  result := fmt.Sprintf("Hello User ID %d with Name %s", test.ID, test.FullName)

  w.Write([]byte(result))
}

func (hlm *HelloWorldModule) GetMultiDataFromDatabase(w http.ResponseWriter, r *http.Request) {
  hlm.stats.Add(1)

  test := []ZakiTest{}
  query := "SELECT id, name FROM zaki_test LIMIT 10"
  err := hlm.DB.Select(&test, query)
  if err != nil {
  	log.Println("Error Query Database. Error: ", err.Error())
  }

  result := "List:\n"
  for _, v := range test {
  	result = fmt.Sprintf("%sHello User ID %d with Name %s\n", result, v.ID, v.FullName)
  }

  w.Write([]byte(result))
}

func (hlm *HelloWorldModule) SearchDataFromDatabase(w http.ResponseWriter, r *http.Request) {
  hlm.stats.Add(1)

  name := r.FormValue("name")

  test := []ZakiTest{}
  query := "SELECT id, name FROM zaki_test WHERE lower(name) = lower($1)"
  err := hlm.DB.Select(&test, query, name)
  if err != nil {
  	log.Println("Error Query Database. Error: ", err.Error())
  }

  result := fmt.Sprintf("List User with Name %s:\n", name)
  for _, v := range test {
  	result = fmt.Sprintf("%s Hello User ID %d with Name %s\n", result, v.ID, v.FullName)
  }

  w.Write([]byte(result))
}