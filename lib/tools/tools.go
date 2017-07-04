package tools

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Database struct {
	Host     string `json:"Host"`
	Login    string `json:"Login"`
	Password string `json:"Password"`
	Database string `json:"Database"`
	Port     string `json:"Port"`
	Protocol string `json:"Protocol"`
}

type Bid struct {
	Sv_id                    int
	S_id										 int
	Bid_at                   time.Time
	Last_bid                 float64
	Calculations						 []Calculations
}

type Calculations struct {
	Type		string
	Value		float64
}

func (d *Database) DataSourceName() string {
	return d.Login + ":" + d.Password + "@" + d.Protocol + "(" + d.Host + ":" + d.Port + ")/" + d.Database
}

func (d *Database) InitConnect() (*sql.DB, error) {

	db, err := sql.Open("mysql", d.DataSourceName())
	if err != nil {
		return nil, errors.New("Bad Server configuration")
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
