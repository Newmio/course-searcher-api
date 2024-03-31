package app

import (
	"fmt"
	"searcher/internal/app/db/postgres"
	"searcher/internal/app/db/redis"
	handlerCourse "searcher/internal/course/handler"
	repoCourse "searcher/internal/course/repository"
	serviceCourse "searcher/internal/course/service"
	handlerMiddleware "searcher/internal/middleware/handler"
	serviceMiddleware "searcher/internal/middleware/service"
	handlerUser "searcher/internal/user/handler"
	repoUser "searcher/internal/user/repository"
	serviceUser "searcher/internal/user/service"

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
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           " [ ${time_custom} ]  ${latency_human}  ${status}   ${method}   ${uri}\n\n",
		CustomTimeFormat: "2006/01/02 15:04:05",
		Output:           color.Output,
	}))

	middlewareService := serviceMiddleware.NewMiddlewareService()
	middlewareHandler := handlerMiddleware.NewHandler(middlewareService)
	middlewares := map[string]echo.MiddlewareFunc{
		"api": middlewareHandler.ParseToken,
	}

	userRepo := repoUser.NewPsqlUserRepo(dbPsql)
	userService := serviceUser.NewUserService(userRepo)
	userHandler := handlerUser.NewHandler(userService)
	userHandler.InitUserRoutes(e, middlewares)

	managerCourseRepo := repoCourse.NewManagerCourseRepo(dbPsql, dbRedis)
	courseService := serviceCourse.NewCourseService(managerCourseRepo)
	courseHandler := handlerCourse.NewHandler(courseService)
	courseHandler.InitCourseRoutes(e, middlewares)

	e.Group("/api", middlewareHandler.ParseToken)

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

	e.Logger.Fatal(e.Start(":8088"))

	return nil
}
