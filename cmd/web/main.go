package main

import (
	"database/sql"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/amirhasanpour/golang-subscription-service/data"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "8080"

func main() {
	// create loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	wg := sync.WaitGroup{}

	// setup the application config (partially)
	app := Config{
		InfoLog: infoLog,
		ErrorLog: errorLog,
		Wait: &wg,
	}

	// define flags directly onto the Config struct
	flag.StringVar(&app.DSN, "DSN", "host=localhost port=5431 user=postgres password=password dbname=subscriptionservice sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	flag.StringVar(&app.REDIS, "REDIS", "127.0.0.1:6379", "Redis connection string")
	flag.Parse()

	// connect to the database with correct DSN
	db := connectToDB(app.DSN)
	if db == nil {
		log.Panic("can't connect to database")
	}
	app.DB = db
	app.Models = data.New(db)

	// initialize sessions with correct Redis address
	app.Session = initSessions(app.REDIS)

	// listen for shutdown
	go app.listenForShutdown()

	// start web server
	app.serve()
}

func (app *Config) serve() {
	// start http server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	app.InfoLog.Println("Starting web server...")
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectToDB(dsn string) *sql.DB {
	counts := 0

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("postgres not yet ready...")
		} else {
			log.Println("connected to database!")
			return connection
		}

		if counts > 10 {
			return nil
		}

		log.Println("Backing off for 1 second")
		time.Sleep(1 * time.Second)
		counts++
	}
}


func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return  db, nil
}

func initSessions(redisAddr string) *scs.SessionManager {
	gob.Register(data.User{})

	session := scs.New()
	session.Store = redisstore.New(initRedis(redisAddr))
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	return session
}

func initRedis(redisAddr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisAddr)
		},
	}
}

func (app *Config) listenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.shutdown()
	os.Exit(0)
}

func (app *Config) shutdown() {
	// perform any cleanup tasks
	app.InfoLog.Println("would run cleanup tasks...")

	// block until waitgroup is empty
	app.Wait.Wait()

	app.InfoLog.Println("closing channels and shutting down application...")
}