package main

import (
    "github.com/codegangsta/negroni"
    "github.com/phyber/negroni-gzip/gzip"
    "github.com/unrolled/secure"
)

func main() {

    router := NewRouter()

    secureMiddleware := secure.New(secure.Options{})

    n := negroni.Classic()
    n.Use(negroni.HandlerFunc(secureMiddleware.HandlerFuncWithNext))
    n.Use(gzip.Gzip(gzip.DefaultCompression))
    n.UseHandler(router)
    n.Run(":3000")

}