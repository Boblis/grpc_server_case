
package main

import (
	"github.com/google/wire"
	"my_grpc/internal/conf"
	"my_grpc/internal/db"
)

func InitApp() (*db.App, error) {
	panic(wire.Build(conf.Provider, db.Provider, NewApp))
}
