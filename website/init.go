package website

import (
	"log"
	"net/http"
	"text/template"
	"github.com/lib/pq"
    "fmt"
    "github.com/jmoiron/sqlx"
    "github.com/garyburd/redigo/redis"
    "time"
)

type User struct {
	ID int `db:"user_id"`
	Name string `db:"full_name"`
	MSISDN string
	Email string `db:"user_email"`
	BirthTime pq.NullTime `db:"birth_date"`
	BirthDate string
	CreateTime time.Time `db:"create_time"`
	CreatedTime string 
	UpdatedTime pq.NullTime `db:"update_time"`
	UpdateTime string
	UserAge int `db:"user_age"`
	Calculation string `db:"-"`
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

	users := []User{}
	userName := "%"+r.FormValue("q")+"%"
	if userName == "%%" {
		err = nwm.db.Select(&users, "SELECT user_id, full_name, msisdn, user_email, birth_date, create_time, update_time, COALESCE(EXTRACT(YEAR from AGE(birth_date)),0) AS user_age FROM WS_USER LIMIT 10;")
	} else {
		err = nwm.db.Select(&users, "SELECT user_id, full_name, msisdn, user_email, birth_date, create_time, update_time, COALESCE(EXTRACT(YEAR from AGE(birth_date)),0) AS user_age FROM WS_USER WHERE full_name LIKE $1 ORDER BY full_name ASC LIMIT 10;", userName)
	}
	if err != nil {
		log.Println(err.Error())
	}

	const constant = 125.25
	for idx, _ := range users {
		// = users[idx].BirthTime.Format("2006/01/02 15:04:05")
		users[idx].CreatedTime = users[idx].CreateTime.Format("2006/01/02 15:04:05")
		users[idx].UpdateTime = "-"
		k, _ := users[idx].UpdatedTime.Value()
		if k != nil {
			users[idx].UpdateTime = k.(time.Time).Format("2006/01/02 15:04:05")
		}

		users[idx].BirthDate = "-"
		k, _ = users[idx].BirthTime.Value()
		if k != nil {
			users[idx].BirthDate = k.(time.Time).Format("2006/01/02 15:04:05")
		}

		users[idx].Calculation = fmt.Sprintf("%.1f", (float64(users[idx].ID) * constant))
	}

	data := map[string]interface{}{
		"users": users,
		"visitorCount": visitorCount,
	}

	err = nwm.render.ExecuteTemplate(w, "home.html", data)
	if err != nil {
		log.Println("Gagal Render Template because: ", err.Error())
	}
}