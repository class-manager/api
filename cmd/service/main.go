package main

import (
	"github.com/class-manager/api/pkg/db"
	"github.com/class-manager/api/pkg/util/env"
)

func main() {
	env.LoadEnv()

	db.Initialise()
}
