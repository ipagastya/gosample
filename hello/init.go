package hello

import (
  "net/http"
  "log"
  "expvar"
  "context"
  "time"
  "fmt"

  logging "gopkg.in/tokopedia/logging.v1"
  "github.com/opentracing/opentracing-go"
)

type ServerConfig struct {
  Name string
}

type Config struct {
  Server ServerConfig
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

  // this message only shows up if app is run with -debug option, so its great for debugging
  logging.Debug.Println("hello init called",cfg.Server.Name)

  return &HelloWorldModule{
    cfg: &cfg,
    something: "John Doe",
    stats : expvar.NewInt("rpsStats"),
  }

}

func (hlm *HelloWorldModule) SayHelloWorld(w http.ResponseWriter, r *http.Request) {
  span, ctx := opentracing.StartSpanFromContext(r.Context(), r.URL.Path)
  defer span.Finish()

  ctx = context.WithValue(ctx, "color", r.FormValue("color"))

  hlm.stats.Add(1)
  hlm.LogValueOfContext(ctx)
  hlm.someSlowFuncWeWantToTrace(ctx, w)
}

func (hlm *HelloWorldModule) LogValueOfContext(ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LogValueOfContext")
	defer span.Finish()

	color := ctx.Value("color")
	if color != nil {
		fmt.Printf("%+v\n", color)
	}
}

func (hlm *HelloWorldModule) someSlowFuncWeWantToTrace(ctx context.Context, w http.ResponseWriter) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "someSlowFuncWeWantToTrace")
	defer span.Finish()

	ctx = context.WithValue(ctx, "color", "red")

	time.Sleep(3 * time.Second)

	w.Write([]byte("Hello " + hlm.something))
}