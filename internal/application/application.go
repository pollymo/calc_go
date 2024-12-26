package application

import (
	"bufio"
	"calc_go/pkg/calculation"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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

type Response struct {
	Result string `json:"result"`
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Wrong Method"}`, http.StatusMethodNotAllowed)
		return
	}

	request := new(Request)

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("Invalid response: %v", err)
		http.Error(w, `{"error":"Invalid Body"}`, http.StatusBadRequest)
		return
	}

	result, err := calculation.Calc(request.Expression)
	if err != nil {
		var errorMsg string
		statusCode := http.StatusUnprocessableEntity

		switch err {
		case calculation.ErrInvalidExpression:
			errorMsg = "Error calculation"
		case calculation.ErrDivisionByZero:
			errorMsg = "Division by zero"
		case calculation.ErrInvalidNumber:
			errorMsg = "Invalid number"
		case calculation.ErrUnexpectedToken:
			errorMsg = "Unexpected token"
		case calculation.ErrNotEnoughValues:
			errorMsg = "Not enough values"
		case calculation.ErrInvalidOperator:
			errorMsg = "Invalid operator"
		case calculation.ErrEmptyInput:
			errorMsg = "Empty input"
		default:
			errorMsg = "Error calculation"
			statusCode = http.StatusUnprocessableEntity
		}

		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, errorMsg), statusCode)
		return
	}
	res := strconv.FormatFloat(result, 'f', -1, 64)
	resp := Response{Result: res}
	jsonResp, _ := json.Marshal(resp)
	if err != nil {
		log.Printf("Error while marshaling response: %v", err)
		http.Error(w, `{"error": "Unknown error occurred"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResp)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (app *Application) RunServer() error {
	http.HandleFunc("/api/v1/calculate", CalcHandler)
	return http.ListenAndServe(":"+app.config.Addr, nil)
}
