package main

import (
	"alexedwards.net/snippetbox/pkg/models/postgresql"
	"database/sql"
	"flag"
	"github.com/golangcollege/sessions"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	news          *mysql.NewsModel
	templateCache map[string]*template.Template
	users         *mysql.UserModel
}

func main() {
	data := "user=postgres password=postgres dbname=news sslmode=disable host=localhost port=5431"
	addr := flag.String("addr", ":8790", "HTTP network address")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	db, err := InitDB(data)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		news:          &mysql.NewsModel{DB: db},
		templateCache: templateCache,
		users:         &mysql.UserModel{DB: db},
	}
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
func InitDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, err
}
