package hello

import (
  "net/http"
  "log"
  "fmt"
  "expvar"
  logging "gopkg.in/tokopedia/logging.v1"
)

//FOR TRAINING
type Config struct {
	Server        ServerConfig
	Training      TrainingConfig  
	ThisUpperCase UpperCaseConfig
	ThisLowerCase LowerCaseConfig
	DifferentName HelloConfig     `gcfg:"Hello"`
}

type ServerConfig struct {
	Name string
}
type TrainingConfig struct {
	UpperCase string
	LowerCase int64
	Alayers   bool
}
type UpperCaseConfig struct {
	CustomVariable float64
}
type LowerCaseConfig struct {
	Goal int64
}
type HelloConfig struct {
	Name string
}

type HelloWorldModule struct {
  cfg   *Config
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

  fmt.Printf("%+v", cfg)

  // this message only shows up if app is run with -debug option, so its great for debugging
  logging.Debug.Println("hello init called",cfg.Server.Name)

  return &HelloWorldModule{
    cfg: &cfg,
    something: "John Doe",
    stats : expvar.NewInt("rpsStats"),
  }

}

func (hlm *HelloWorldModule) SayHelloWorld(w http.ResponseWriter, r *http.Request) {
  hlm.stats.Add(1)
  w.Write([]byte("Hello " + hlm.something))
}
