package main

import (
	"ass1/pkg/models/mysql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"runtime/debug"
)

var templates *template.Template

type application struct {
	errorLog  *log.Logger
	infoLog   *log.Logger
	NewsModel *mysql.NewsModel
}

func main() {
	newsModel := &mysql.NewsModel{}
	err := newsModel.InitDB("temirkhan:temirkhan322@tcp(localhost:3306)/newsdb")
	if err != nil {
		log.Fatal(err)
	}

	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog:  errorLog,
		infoLog:   infoLog,
		NewsModel: newsModel,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	infoLog.Printf("Starting server on %s", *addr)
	err2 := srv.ListenAndServe()
	errorLog.Fatal(err2)
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
