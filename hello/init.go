package hello

import (
	"expvar"
	"log"
	"net/http"
	"text/template"

	logging "gopkg.in/tokopedia/logging.v1"
)

type ServerConfig struct {
	Name string
}

type Config struct {
	Server ServerConfig
}

type HelloWorldModule struct {
	cfg       *Config
	something string
	stats     *expvar.Int
	render    *template.Template  //FOR TRAINING
}

func NewHelloWorldModule() *HelloWorldModule {

	var cfg Config

	ok := logging.ReadModuleConfig(&cfg, "config", "hello") || logging.ReadModuleConfig(&cfg, "files/etc/gosample", "hello")
	if !ok {
		// when the app is run with -e switch, this message will automatically be redirected to the log file specified
		log.Fatalln("failed to read config")
	}

	// this message only shows up if app is run with -debug option, so its great for debugging
	logging.Debug.Println("hello init called", cfg.Server.Name)

	//FOR TRAINING
	engine := template.Must(template.ParseGlob("files/var/templates/*"))

	return &HelloWorldModule{
		cfg:       &cfg,
		something: "John Doe",
		stats:     expvar.NewInt("rpsStats"),
		render:    engine,    //FOR TRAINING
	}

}

func (hlm *HelloWorldModule) SayHelloWorld(w http.ResponseWriter, r *http.Request) {
	hlm.stats.Add(1)
	w.Write([]byte("Hello " + hlm.something))
}

//FOR TRAINING
type Animal struct {
	Name string
	Type string
	Legs int
}

func (hlm *HelloWorldModule) SayMyName(w http.ResponseWriter, r *http.Request) {
	hlm.stats.Add(1)

	thisMap := map[string]string{
		"bca":     "5270366793",
		"mandiri": "010002848657575",
		"bri":     "357575774",
		"bni":     "08348375756475",
	}
	
	thisArr := []string{"Apple", "Orange", "Banana"}

	thisStruct := Animal{
		Name: "Alvin",
		Type: "Anjing",
		Legs: 4,
	}

	data := map[string]interface{}{
		"name":   "My Name",
		"age":    "17",
		"banks":  thisMap,
		"fruits": thisArr,
		"dog":    thisStruct,
	}
	
	err := hlm.render.ExecuteTemplate(w, "home.html", data)
	if err != nil {
		log.Println("Gagal Render Template because: ", err.Error())
	}
}
