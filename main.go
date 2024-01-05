package main

import (
	"fmt"
	"sgnacos/conf"
	"sgnacos/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	conf.InitConf()

	app := fiber.New()
	app.Use(logger.New())

	router.Http(app)

	err := app.Listen(fmt.Sprintf("0.0.0.0:%d", conf.BaseConf.Server.Port))
	if err != nil {
		return
	}
}
