package main

import (
	"appengine"
	"appengine/datastore"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"time"
)
type users struct {
	UserName       string
	Password       string
	DateRegistered time.Time
}
func init() {

	router := mux.NewRouter().StrictSlash(true)

	router.Handle("/static/CSS/", http.StripPrefix("/static/CSS/", http.FileServer(http.Dir("static/CSS"))))
	router.HandleFunc("/main", mainHandler)
	router.HandleFunc("/", rootHandler)

	http.Handle("/", router)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	us := users{
		UserName:       r.FormValue("username"),
		Password:       r.FormValue("password"),
		DateRegistered: time.Now(),
	}

	key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "tblUsers", nil), &us)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, "")
		return
	}

	var e2 users
	if err = datastore.Get(c, key, &e2); err != nil {
		errorHandler(w, r, http.StatusInternalServerError, "")
		return
	}

	//fmt.Fprintf(w, "Stored and retrieved the Employee named %q", us.UserName)

	if r.URL.Path != "/main" {
		errorHandler(w, r, http.StatusNotFound, "")
		return
	}
	page := template.Must(template.ParseFiles(
		"static/_base.gtpl",
		"static/main.gtpl",
	))

	if err := page.Execute(w, us); err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound, "")
		return
	}

	page := template.Must(template.ParseFiles(
		"static/_base.gtpl",
		"static/index.gtpl",
	))
	if err := page.Execute(w, nil); err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

}

func errorHandler(w http.ResponseWriter, r *http.Request, status int, err string) {
	w.WriteHeader(status)
	switch status {

	case http.StatusNotFound:
		page := template.Must(template.ParseFiles(
			"static/_base.gtpl",
			"static/404.gtpl",
		))
		if err := page.Execute(w, nil); err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err.Error())
			return
		}

	case http.StatusInternalServerError:
		page := template.Must(template.ParseFiles(
			"static/_base.gtpl",
			"static/500.gtpl",
		))
		if err := page.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
