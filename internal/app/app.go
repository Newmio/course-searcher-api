package app

import (
	"fmt"
	"searcher/internal/app/db/postgres"
	"searcher/internal/app/db/redis"
	"searcher/internal/course/handler"
	"searcher/internal/course/repository"
	"searcher/internal/course/service"

	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitProject() error {

	dbPsql, err := postgres.OpenDb()
	if err != nil {
		return err
	}

	dbRedis, err := redis.OpenDb()
	if err != nil {
		return err
	}

	e := echo.New()

	// userRepo := user.NewPsqlUserRepo(dbPsql)
	// userService := user.NewUserService(userRepo)
	// userHandler := user.NewHandler(userService)
	// userHandler.InitUserRoutes(e)

	managerCourseRepo := repository.NewManagerCourseRepo(dbPsql, dbRedis)
	courseService := service.NewCourseService(managerCourseRepo)
	courseHandler := handler.NewHandler(courseService)
	courseHandler.InitCourseRoutes(e)

	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "\n\n [ ${time_custom} ]  ${latency_human}  ${status}   ${method}   ${uri}",
		CustomTimeFormat: "2006/01/02 15:04:05",
		Output:           color.Output,
	}))

	e.GET("/", func(c echo.Context) error {
		return c.File("static/index.html")
	})

	color.New(color.BgMagenta, color.Bold).Println()
	for _, value := range e.Routes() {
		var customColor color.Attribute

		switch value.Method {
		case "GET":
			customColor = color.BgHiGreen

		case "POST":
			customColor = color.BgHiCyan

		case "PUT":
			customColor = color.BgHiMagenta
		}

		color.New(customColor, color.Bold).Print(fmt.Sprintf("\n\t%s : %s", value.Path, value.Method))
	}
	color.New(color.BgMagenta, color.Bold).Println()

	e.Logger.Fatal(e.Start(":8080"))

	return nil
}
