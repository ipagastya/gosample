package hello

import (
	"expvar"
	"log"
	"net/http"
	"text/template"
	logging "gopkg.in/tokopedia/logging.v1"
	"database/sql"
    _ "github.com/lib/pq"
    "fmt"
    "github.com/jmoiron/sqlx"
    "github.com/garyburd/redigo/redis"
    "time"
)

type User struct {
	ID int `db:"user_id"`
	Name sql.NullString `db:"user_name"`
	MSISDN sql.NullString
	Email string `db:"user_email"`
	BirthDate pq.NullTime `db:"birth_date"`
	CreatedTime time.Time `db:"create_time"`
	UpdateTime pq.NullTime `db:"update_time"`
	UserAge int `db:"-"`
}

type ServerConfig struct {
	Name string
}

type Config struct {
	Server ServerConfig
}

type WebsiteModule struct {
	cfg       *Config
	render    *template.Template  //FOR TRAINING
	db 		  *sqlx.DB
	pool 	  *redis.Pool
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

	db, err := sqlx.Connect("postgres", "postgres://da161205:123Toped456@devel-postgre.tkpd/tokopedia-user?sslmode=disable")
	if err != nil {
        log.Fatalln(err)
    }

	//FOR TRAINING
	engine := template.Must(template.ParseGlob("files/var/templates/*"))

	return &WebsiteModule{
		cfg:       	&cfg,
		render:    	engine,
		db:			&db,
		pool: 		&redis.Pool{
			            MaxIdle:     3,
			            IdleTimeout: 240 * time.Second,
			            Dial: func() (redis.Conn, error) {
			                conn, err := redis.Dial("tcp", "devel-redis.tkpd:6379")
			                if err != nil {
			                    return nil, err
			                }
			                return conn, err
			          	},
			        }
	}
}

func (nwm *NewWebsiteModule) RenderWebpage(w http.ResponseWriter, r *http.Request) {
	conn := nwm.pool.Get()
	visitorCount, err := conn.Do("INCR", "visitors")
	if err != nil {
		log.Println(err.Error())
	}

	user := []User{}
	user_name := "%"+r.FormValue("q")+"%"
	var query string
	if user_name != "%%" {
		query = "SELECT user_id, COALESCE(user_name,'-'), COALESCE(msisdn,'-'), email, COALESCE(brith_date,'-'), COALESCE(create_time, date_trunc('second', now()::timestamp)), COALESCE(update_time, '-'), COALESCE(EXTRACT(YEAR from AGE(birth_date)),'0') AS user_age FROM WS_USER")
	} else {
		query = "SELECT user_id, COALESCE(user_name,'-'), COALESCE(msisdn,'-'), email, COALESCE(birth_date,'-'), COALESCE(create_time, date_trunc('second', now()::timestamp)), COALESCE(update_time, '-'), COALESCE(EXTRACT(YEAR from AGE(birth_date)),'0') AS user_age FROM WS_USER WHERE user_name LIKE $1 ORDER BY user_name ASC LIMIT 10;", user_name)
	}
	nwm.db.Select(&user, query)
	
	calculation := []String{}
	for _, usr := range user {
		calculation = append(fmt.Sprintf("%.1f", (usr.ID * 25.25)))
	}

	data := map[string]interface{}{
		"user": user,
		"visitorCount": visitorCount,
		"calculation": calculation,
	}

	err = nwm.render.ExecuteTemplate(w, "home.html", data)
	if err != nil {
		log.Println("Gagal Render Template because: ", err.Error())
	}
}