package website

import (
	"log"
	"net/http"
	"text/template"
	"database/sql"
    "github.com/lib/pq"
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
	render    *template.Template  //FOR TRAINING
	db 		  *sqlx.DB
	pool 	  *redis.Pool
}

func NewWebsiteModule() *WebsiteModule {
	db, err := sqlx.Connect("postgres", "postgres://da161205:123Toped456@devel-postgre.tkpd/tokopedia-user?sslmode=disable")
	if err != nil {
        log.Fatalln(err)
    }

	//FOR TRAINING
	engine := template.Must(template.ParseGlob("files/var/templates/*"))

	return &WebsiteModule{
		cfg:       	&cfg,
		render:    	engine,
		db:			db,
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
			        },
	}
}

func (nwm *WebsiteModule) RenderWebpage(w http.ResponseWriter, r *http.Request) {
	conn := nwm.pool.Get()
	visitorCount, err := conn.Do("INCR", "visitors")
	if err != nil {
		log.Println(err.Error())
	}

	user := []User{}
	user_name := "%"+r.FormValue("q")+"%"
	if user_name != "%%" {
		nwm.db.Select(&user, "SELECT user_id, COALESCE(user_name,'-'), COALESCE(msisdn,'-'), email, COALESCE(birth_date,'-'), COALESCE(create_time, date_trunc('second', now()::timestamp)), COALESCE(update_time, '-'), COALESCE(EXTRACT(YEAR from AGE(birth_date)),0) AS user_age FROM WS_USER;")
	} else {
		nwm.db.Select(&user, "SELECT user_id, COALESCE(user_name,'-'), COALESCE(msisdn,'-'), email, COALESCE(birth_date,'-'), COALESCE(create_time, date_trunc('second', now()::timestamp)), COALESCE(update_time, '-'), COALESCE(EXTRACT(YEAR from AGE(birth_date)),0) AS user_age FROM WS_USER WHERE user_name LIKE $1 ORDER BY user_name ASC LIMIT 10;", user_name)
	}
	
	calculation := []string{}
	const constant = 125.25
	for _, usr := range user {
		calculation = append(calculation, fmt.Sprintf("%.1f", (float64(usr.ID) * constant)))
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