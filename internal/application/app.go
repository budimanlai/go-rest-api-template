package application

import (
	"github.com/jmoiron/sqlx"
)

type ApplicationContextBase struct {
	// Database connection
	Db *sqlx.DB
}

var (
	// AppContext is the global application context
	AppContext *ApplicationContextBase
)

func init() {
	AppContext = &ApplicationContextBase{
		Db: nil, // This will be set when the database connection is opened
	}
}
