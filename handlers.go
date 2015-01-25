package main

import (
    "encoding/json"
    "fmt"
    "time"
    "net/http"

    "github.com/gorilla/mux"
     jwt "github.com/dgrijalva/jwt-go"
    "io/ioutil" 
)

var (
    privateKey []byte
    publicKey []byte
)
 
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
  
func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Welcome!")
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
    RepoGetAll()
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(todos); err != nil {
        panic(err)
    }
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    todoId := vars["todoId"]
    fmt.Fprintln(w, "Todo show:", todoId)
}

func TodoCreate(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    var t Todo   
    err := decoder.Decode(&t)
    if err != nil {
        panic(err)
    }
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)

    nt, err := RepoCreateTodo(t)

    if err != nil {
        panic(err)
    }

    if err := json.NewEncoder(w).Encode(nt); err != nil {
        panic(err)
    }
}


func TokenGet(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    username := vars["username"]
    password := vars["password"]

    userId, _ := TestCredentials(username, password)

    // Sign and get the complete encoded token as a string
    tokenString, err := CreateTokenString(userId)
  
    if err != nil {
        panic(err)
    }

    w.WriteHeader(http.StatusOK)
    // w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    
    if err := json.NewEncoder(w).Encode(tokenString); err != nil {
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