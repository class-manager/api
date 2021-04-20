package main

import (
	"github.com/class-manager/api/pkg/db"
	http_server "github.com/class-manager/api/pkg/http"
	"github.com/class-manager/api/pkg/util/env"
)

func main() {
	env.LoadEnv()

	db.Initialise()
	http_server.Start()
}
