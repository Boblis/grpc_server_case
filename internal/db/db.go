package db

import (
	"database/sql"
	"fmt"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"my_grpc/internal/conf"
)

var Provider = wire.NewSet(New)

type App struct {
	Db *sql.DB
}

type Dao struct {
	App App
}

type UserInfo struct {
	Name string
	Age int
	Address string
}

var GlobalApp *App

func New(cfg *conf.Config) (db *sql.DB, err error) {
	db, err = sql.Open("mysql", cfg.Database.Dsn)
	if err != nil {
		return
	}
	if err = db.Ping(); err != nil {
		return
	}
	return
}

func NewDao() (*Dao, error) {
	if GlobalApp == nil {
		return nil, errors.New("the global app is nil")
	}
	return &Dao{
		App: *GlobalApp,
	}, nil
}


func (d *Dao) QueryUserInfo(id string) (*UserInfo, error) {
	queryCmd := fmt.Sprintf(`select name, age, address from user where id = ?`)
	row := d.App.Db.QueryRow(queryCmd, id)
	userInfo := &UserInfo{}
	if err := row.Scan(
		userInfo.Name,
		userInfo.Age,
		userInfo.Address,
		); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(err, "no record match")
		} else {
			return nil, errors.Wrapf(err, "query user info error id value: %v", id)
		}
	}
	return userInfo, nil
}
