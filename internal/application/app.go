package application

import (
	"github.com/jmoiron/sqlx"
)

type ApplicationContextBase struct {
	// Add fields as needed for your application context
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
