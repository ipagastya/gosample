package hello

import (
	"expvar"
	"log"
	"net/http"
	"encoding/json"

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

	return &HelloWorldModule{
		cfg:       &cfg,
		something: "John Doe",
		stats:     expvar.NewInt("rpsStats"),
	}

}

func (hlm *HelloWorldModule) SayHelloWorld(w http.ResponseWriter, r *http.Request) {
	hlm.stats.Add(1)
	w.Write([]byte("Hello " + hlm.something))
}

//FOR TRAINING
type apiResponse struct {
	Status  int         `json:"status"`
	Message []string    `json:"message"`
	Data    interface{} `json:"data"`
}

type Book struct {
	ISBN       string
	Title      string  `json:"book_title"`
	Author     string  `json:"-"`
	PageCount  int     `json:"page_count"`
	Price      float64 `json:"book_price"`
	Translated bool    `json:"translated"`
	OmitEmpty  string  `json:"empty,omitempty"`
	NotRender  string  `json:"-"`
	private    string  `json:"please_render_me"`
}

func (hlm *HelloWorldModule) BookInfoHandler(w http.ResponseWriter, r *http.Request) {
	hlm.stats.Add(1)
	
	myBook := Book{
		ISBN:       "IND-08187457364",
		Title:      "How to be a great Nakama",
		Author:     "Leontinus Tanuwijaya",
		PageCount:  230,
		Price:      120.5,
		Translated: true,
		NotRender:  "This won't be rendered because you opt to not render it.",
		private:    "Whatever you put here wont be render because it's private",
	}

	response := apiResponse{
		Status: 200,
		Data:   myBook,
	}
	//result, err := json.Marshal(response)
	result, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		log.Println("Cannot Marshal JSON because ", err.Error())
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}
