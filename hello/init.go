package hello

import (
  "net/http"
  "log"
  "fmt"
  "expvar"

  logging "gopkg.in/tokopedia/logging.v1"
  "github.com/garyburd/redigo/redis"
)

type ServerConfig struct {
  Name string
}
type RedisConfig struct {
  Connection string
}
type Config struct {
  Server ServerConfig
  Redis  RedisConfig
}

type HelloWorldModule struct {
  cfg   *Config
  Redis *redis.Pool
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

  redisPools := &redis.Pool{
	Dial: func() (redis.Conn, error) {
		conn, err := redis.Dial("tcp", cfg.Redis.Connection)
		if err != nil {
			return nil, err
		}
		return conn, err
	},
  }

  // this message only shows up if app is run with -debug option, so its great for debugging
  logging.Debug.Println("hello init called",cfg.Server.Name)

  return &HelloWorldModule{
    cfg: &cfg,
    Redis: redisPools,
    something: "John Doe",
    stats : expvar.NewInt("rpsStats"),
  }

}

func (hlm *HelloWorldModule) SayHelloWorld(w http.ResponseWriter, r *http.Request) {
  hlm.stats.Add(1)
  w.Write([]byte("Hello " + hlm.something))
}

func (hlm *HelloWorldModule) SetRedis(w http.ResponseWriter, r *http.Request) {
  hlm.stats.Add(1)

  key := r.FormValue("key")
  value := r.FormValue("value")

  result := "Set Redis Success"

  pool := hlm.Redis.Get()
  _, err := redis.String(pool.Do("SET", key, value))
  if err != nil {
  	log.Printf("Failed to Set key %s with value %s. Error: %s\n", key, value, err.Error())
  	result = "Set Redis Failed"
  }

  pool.Do("EXPIRE", key, 10)

  w.Write([]byte(result))
}

func (hlm *HelloWorldModule) GetRedis(w http.ResponseWriter, r *http.Request) {
  hlm.stats.Add(1)

  key := r.FormValue("key")

  result := "Get Redis Success"

  pool := hlm.Redis.Get()
  value, err := redis.String(pool.Do("GET", key))
  if err != nil {
  	log.Printf("Failed to Get key %s. Error: %s\n", key, err.Error())
  	result = "Get Redis Failed."
  }

  result = fmt.Sprintf("%s\nValue Redis with key %s is %s", result, key, value)

  w.Write([]byte(result))
}