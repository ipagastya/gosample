package hello

import (
  "net/http"
  "log"
  "expvar"
  "encoding/json"

  "github.com/nsqio/go-nsq"
  logging "gopkg.in/tokopedia/logging.v1"
)

type ServerConfig struct {
  Name string
}
//FOR TRAINING
type NSQConfig struct {
  NSQD     string
  Lookupds string
}

type Config struct {
  Server ServerConfig
  NSQ    NSQConfig
}

type HelloWorldModule struct {
  cfg   *Config
  NSQ   *nsq.Producer
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

  producer, err := nsq.NewProducer(cfg.NSQ.NSQD, nsq.NewConfig())
  if err != nil {
	log.Fatalln("failed to create nsq producer. Error: ", err.Error())
  }

  // this message only shows up if app is run with -debug option, so its great for debugging
  logging.Debug.Println("hello init called",cfg.Server.Name)

  return &HelloWorldModule{
    cfg: &cfg,
    NSQ: producer,
    something: "John Doe",
    stats : expvar.NewInt("rpsStats"),
  }

}

func (hlm *HelloWorldModule) SayHelloWorld(w http.ResponseWriter, r *http.Request) {
  hlm.stats.Add(1)
  w.Write([]byte("Hello " + hlm.something))
}

//FOR TRAINING
func (hlm *HelloWorldModule) PublishNSQ(w http.ResponseWriter, r *http.Request) {
  hlm.stats.Add(1)
  
  message := r.FormValue("message")
  name := r.FormValue("name")

  result := "Push NSQ Success"
  data := map[string]string{
  	"name":    name,
  	"message": message,
  }
  nsqMessage, _ := json.Marshal(data) //<<--- ingatkan jangan sering2 pakai _

  err := hlm.NSQ.Publish("Training_NSQ_Mine", nsqMessage)
  if err != nil {
  	log.Println("Failed to publish NSQ message. Error: ", err)
  	result = "Push NSQ Success"
  }

  w.Write([]byte(result))
}