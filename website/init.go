package hello

import (
	"expvar"
	"log"
	"net/http"
	"text/template"
	logging "gopkg.in/tokopedia/logging.v1"
	"database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "github.com/jmoiron/sqlx"
    "github.com/garyburd/redigo/redis"
)

type User struct {
	ID int `db:"user_id"`
	Name string `db:"user_name"`
	MSISDN string
	Email string
	birth_date time.Time
	created_time time.Time
	update_time time.Time
	user_age int `db:"-"`
}

func init() {
	*db, err = sqlx.Connect("postgres", "postgres://da161205:123Toped456@devel-postgre.tkpd/tokopedia-user?sslmode=disable")
	if err != nil {
        log.Fatalln(err)
    }
    *c, err = redis.Dial("tcp", "devel-redis.tkpd:6379")
    c.Do("SET", "visitors", 0)
}

type ServerConfig struct {
	Name string
}

type Config struct {
	Server ServerConfig
}

type WebsiteModule struct {
	cfg       *Config
	visitorCount string
	stats     *expvar.Int
	render    *template.Template  //FOR TRAINING
	db 		  *DB
	c 		  Conn
}

func NewWebsiteModule() *WebsiteModule {

	var cfg Config

	ok := logging.ReadModuleConfig(&cfg, "config", "website") || logging.ReadModuleConfig(&cfg, "files/etc/gosample", "website")
	if !ok {
		// when the app is run with -e switch, this message will automatically be redirected to the log file specified
		log.Fatalln("failed to read config")
	}

	// this message only shows up if app is run with -debug option, so its great for debugging
	logging.Debug.Println("hello init called", cfg.Server.Name)

	//FOR TRAINING
	engine := template.Must(template.ParseGlob("files/var/templates/*"))

	return &WebsiteModule{
		cfg:       &cfg,
		visitorCount: c.Do("GET", "visitors"),
		render:    engine,
	}
}

func (nwm *NewWebsiteModule) RenderWebpage(w http.ResponseWriter, r *http.Request, data) {
	visitorCount := c.Do("INCR", "visitors")

	user := []User{}
	err = db.Select(&user, "SELECT user_id, user_name, msisdn, email, birth_date, created_time, update_time, COALESCE(EXTRACT(YEAR from AGE(birth_date)),"0") AS user_age FROM WS_USER WHERE user_name LIKE $1 ORDER BY user_name ASC LIMIT 10", r.formValue())
	data := map[string]interface{}{
		"user": user,
		"visitorCount": visitorCount,
	}

	err := nwm.render.ExecuteTemplate(w, "home.html", data)
	if err != nil {
		log.Println("Gagal Render Template because: ", err.Error())
	}
}