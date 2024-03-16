package app

import (
	"searcher/internal/app/db/postgres"
	"searcher/internal/app/db/redis"
	"searcher/internal/course"
	"searcher/internal/user"

	"github.com/labstack/echo/v4"
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

	userRepo := user.NewPsqlUserRepo(dbPsql)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewHandler(userService)
	e = userHandler.InitUserRoutes(e)

	courseRepo := course.NewCourseRepo(dbPsql, dbRedis)
	courseService := course.NewCourseService(courseRepo)
	courseHandler := course.NewHandler(courseService)
	e = courseHandler.InitCourseRoutes(e)

	return e.Start(":8080")
}
