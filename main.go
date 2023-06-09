package main

import (
	"os"

	"github.com/abdou-1614/go-rest-api/common"
	"github.com/abdou-1614/go-rest-api/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	err := run()

	if err != nil {
		panic(err)
	}
}

func run() error {

	err := common.LoadEnv()

	if err != nil {
		return err
	}

	err = common.InitDB()

	if err != nil {
		return err
	}

	defer common.CloseDB()

	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(recover.New())

	router.AddTodoGroupe(app)

	var port string

	if port = os.Getenv("PORT"); port == "" {
		port = "3000"
	}

	app.Listen(":" + port)

	return nil

}
