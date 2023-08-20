package main

import (
	sharedd "shared"

	app "analytics_service/internal/application"
)

func main() {
	sharedd.StoreSomething()
	app.Run()
}
