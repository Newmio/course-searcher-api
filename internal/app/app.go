package app

import (
	"fmt"
	"searcher/internal/app/db/postgres"
	"searcher/internal/app/db/redis"
	"searcher/internal/course"
	"searcher/internal/user"

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

	e.Use(middleware.Recover())

	userRepo := user.NewPsqlUserRepo(dbPsql)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewHandler(userService)
	e = userHandler.InitUserRoutes(e)

	managerCourseRepo := course.NewManagerCourseRepo(dbPsql, dbRedis)
	courseService := course.NewCourseService(managerCourseRepo)
	courseHandler := course.NewHandler(courseService)
	e = courseHandler.InitCourseRoutes(e)

	for _, value := range e.Routes() {
		var customColor color.Attribute

		switch value.Method{
		case "GET":
			customColor = color.BgHiGreen

		case "POST":
			customColor = color.BgHiCyan

		case "PUT":
			customColor = color.BgHiMagenta
		}

		color.New(customColor, color.Bold).Println(fmt.Sprintf("\n\t%s : %s", value.Path, value.Method))
	}

	e.Logger.Fatal(e.Start(":8080"))

	return nil
}
