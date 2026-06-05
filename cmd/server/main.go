package main

import (
	"log"

	"edu-admin/internal/app/bootstrap"
)

func main() {
	app, err := bootstrap.NewApplication()
	if err != nil {
		log.Fatal(err)
	}

	if err = app.Run(); err != nil {
		log.Fatal(err)
	}
}
