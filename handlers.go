package main

import (
	"encoding/json"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"io/ioutil"
)

var (
	privateKey []byte
	publicKey  []byte
)

type Token struct {
	Expiry time.Time `json:"expiry"`
	Value  string    `json:"value"`
}

type ClientError struct {
	Message string `json:message`
}

type objectDelegate func(*json.Encoder) error

func init() {
	/*
	   openssl genrsa -out demo.rsa 2056 # the 1024 is the size of the key we are generating
	   openssl rsa -in demo.rsa -pubout > demo.rsa.pub
	*/
	var e error
	privateKey, e = ioutil.ReadFile("test/demo.rsa")
	if e != nil {
		panic(e.Error())
	}
	publicKey, _ = ioutil.ReadFile("test/demo.rsa.pub")
}

func send(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func sendJson(w http.ResponseWriter, status int, fn objectDelegate) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	send(w, status)
	if err := fn(json.NewEncoder(w)); err != nil {
		panic(err)
	}
}

func sendText(w http.ResponseWriter, status int, text string) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	send(w, status)
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
	RepoGetAll()

	sendJson(w, http.StatusOK, func(x *json.Encoder) error {
		return x.Encode(todos)
	})
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	todoId := mux.Vars(r)["todoId"]

	if todo, err := RepoFindTodo(todoId); err != nil {
		sendText(w, http.StatusNotFound, "")
	} else {
		sendJson(w, http.StatusOK, func(x *json.Encoder) error { return x.Encode(todo) })
	}
}

func TodoDelete(w http.ResponseWriter, r *http.Request) {
	todoId := mux.Vars(r)["todoId"]

	if err := RepoDestroyTodo(todoId); err != nil {
		sendText(w, http.StatusNotFound, "")
	} else {
		sendText(w, http.StatusNoContent, "")
	}
}

func TodoCreate(w http.ResponseWriter, r *http.Request) {
	var t Todo

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&t); err != nil {
		sendJson(w, http.StatusBadRequest, func(x *json.Encoder) error {
			return x.Encode(ClientError{
				Message: err,
			})
		})
	}

	if nt, err := RepoCreateTodo(t); err != nil {
		panic(err)
	} else {
		sendJson(w, http.StatusOK, func(x *json.Encoder) error { return x.Encode(nt) })
	}
}

func TokenGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	password := vars["password"]

	userId, _ := TestCredentials(username, password)

	// Sign and get the complete encoded token as a string
	tokenString, err := CreateTokenString(userId)

	to := Token{
		Value:  tokenString,
		Expiry: time.Now(),
	}

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// fmt.Fprintln(w, tokenString)

	if err := json.NewEncoder(w).Encode(to); err != nil {
		panic(err)
	}
}

func TestCredentials(username string, password string) (string, error) {
	if len(username) != 0 {
		return "sdfsdfsdfsfsdfsf", nil
	}
	return "", nil
}

func CreateTokenString(userId string) (string, error) {
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	token.Claims["userId"] = userId
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	// Sign and get the complete encoded token as a string
	return token.SignedString(privateKey)

}

func authRequest(w http.ResponseWriter, r *http.Request) {
	token, _ := jwt.ParseFromRequest(r, func(t *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if token.Valid {
		//YAY!
	} else {
		//Someone is being funny
	}
}
