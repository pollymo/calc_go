package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	config *Config
}

func NewApplication() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

func (app *Application) Run() error {
	for {
		log.Println("input expression")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("error reading input")
		}
		text = strings.TrimSpace(text)
		if text == "exit" {
			log.Println("application was successfully exited")
		}

		result, err := calculation.Calc(text)
		if err != nil {
			log.Println(text, "error calculation with error", err)
		} else {
			log.Println(text, "=", result)
		}
	}
}

type Request struct {
	Expression string `json:"expression"`
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	request := new(Request)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := calculation.Calc(request.Expression)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	} else {
		fmt.Fprintf(w, "result: %f", result)
	}
	w.Header().Set("Content-Type", "application/json")
}

func (app *Application) RunServer() error {
	http.HandleFunc("/", CalcHandler)
	return http.ListenAndServe(":"+app.config.Addr, nil)
}
