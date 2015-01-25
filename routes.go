package main

import "net/http"

type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

var routes = []Route{
    Route{
        "Index",
        "GET",
        "/",
        Index,
    },
    Route{
        "TodoIndex",
        "GET",
        "/todos",
        TodoIndex,
    },
    Route{
        "TodoShow",
        "GET",
        "/todos/{todoId}",
        TodoShow,
    },
    Route{
        "TodoCreate",
        "POST",
        "/todos",
        TodoCreate,
    },  
    Route{
        "TokenGet",
        "GET",
        "/token",
        TokenGet,
    },
}