package main

import (
	"github.com/codegangsta/negroni"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/rs/cors"
	"github.com/unrolled/secure"
)

func main() {

	router := NewRouter()
	secureMiddleware := secure.New(secure.Options{})
	c := cors.Default()
	n := negroni.Classic()
	n.Use(negroni.HandlerFunc(secureMiddleware.HandlerFuncWithNext))
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.Use(c)
	n.UseHandler(router)
	n.Run(":3000")

}
