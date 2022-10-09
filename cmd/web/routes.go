package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/blog", app.blog)
	mux.HandleFunc("/blog/view", app.blogView)
	mux.HandleFunc("/blog/create", app.blogCreate)
	// mux.HandleFunc("/user/create", app.userCreate)

	return mux
}