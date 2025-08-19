package main

import "Crash-Auth-service/internal/container"

func main() {
	app := container.Build()

	app.Run()
}
