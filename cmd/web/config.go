package main

import (
	"database/sql"
	"log"
	"sync"

	"github.com/alexedwards/scs/v2"
	"github.com/amirhasanpour/golang-subscription-service/data"
)

type Config struct {
	Session  	  *scs.SessionManager
	DB 		 	  *sql.DB
	InfoLog  	  *log.Logger
	ErrorLog 	  *log.Logger
	Wait     	  *sync.WaitGroup
	Models	 	  data.Models
	DSN		 	  string
	REDIS         string
	Mailer 	 	  Mail
	ErrorChan 	  chan error
	ErrorChanDone chan bool
}