package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"searcher/internal/app/db/postgres"
	"searcher/internal/app/db/redis"
	handlerCourse "searcher/internal/course/handler"
	repoCourse "searcher/internal/course/repository"
	serviceCourse "searcher/internal/course/service"
	handlerFile "searcher/internal/file/handler"
	repoFile "searcher/internal/file/repository"
	serviceFile "searcher/internal/file/service"
	handlerMiddleware "searcher/internal/middleware/handler"
	serviceMiddleware "searcher/internal/middleware/service"
	handlerUser "searcher/internal/user/handler"
	repoUser "searcher/internal/user/repository"
	serviceUser "searcher/internal/user/service"
	handlerView "searcher/internal/view/handler"
	repoView "searcher/internal/view/repository"
	serviceView "searcher/internal/view/service"
	"syscall"
	"time"

	redisClient "github.com/go-redis/redis/v8"

	"github.com/fatih/color"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type App struct {
	Psql  *sqlx.DB
	Redis *redisClient.Client
	Log   *zap.Logger
	Echo  *echo.Echo
}

func InitProject() *App {

	dbPsql, err := postgres.OpenDb()
	if err != nil {
		panic(err)
	}

	dbRedis, err := redis.OpenDb()
	if err != nil {
		panic(err)
	}

	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	return &App{
		Psql:  dbPsql,
		Redis: dbRedis,
		Log:   log,
		Echo:  echo.New(),
	}
}

func (app *App) Run() {
	app.initService()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM)
	signal.Notify(quit, syscall.SIGINT)

	go func() {
		<-quit
		app.Log.Info("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		app.Psql.Close()
		app.Redis.Close()
		app.Echo.Shutdown(ctx)
	}()

	app.Echo.Logger.Fatal(app.Echo.Start(":8088"))
}

func (app *App) initService() {
	app.Echo.Static("/template", "template")

	app.Echo.Use(middleware.Recover())
	app.Echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           " [ ${time_custom} ]  ${latency_human}  ${status}   ${method}   ${uri}\n\n",
		CustomTimeFormat: "2006/01/02 15:04:05",
		Output:           color.Output,
	}))
	app.Echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
	}))
	app.Echo.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", "no-cache, no-store")
			return next(c)
		}
	})

	middlewareService := serviceMiddleware.NewMiddlewareService()
	middlewareHandler := handlerMiddleware.NewHandler(middlewareService)
	middlewares := map[string]echo.MiddlewareFunc{
		"api": middlewareHandler.ParseToken,
	}

	managerUserRepo := repoUser.NewManagerUserRepo(app.Psql)
	userService := serviceUser.NewUserService(managerUserRepo)
	userHandler := handlerUser.NewHandler(userService)
	userHandler.InitUserRoutes(app.Echo, middlewares)

	managerCourseRepo := repoCourse.NewManagerCourseRepo(app.Psql, app.Redis)
	courseService := serviceCourse.NewCourseService(managerCourseRepo)
	courseHandler := handlerCourse.NewHandler(courseService)
	courseHandler.InitCourseRoutes(app.Echo, middlewares)

	managerFileRepo := repoFile.NewManagerFileRepo(app.Psql)
	FileService := serviceFile.NewFileService(managerFileRepo, managerUserRepo)
	FileHandler := handlerFile.NewHandler(FileService)
	FileHandler.InitFileRoutes(app.Echo, middlewares)

	viewRepo := repoView.NewDiskViewRepo()
	viewService := serviceView.NewViewService(viewRepo, managerUserRepo)
	viewHandler := handlerView.NewHandler(viewService)
	viewHandler.InitViewRoutes(app.Echo, middlewares)

	printRoutes(app.Echo)
}

func printRoutes(e *echo.Echo) {
	color.New(color.BgHiBlack, color.Bold).Println("                                                    ")

	for _, value := range e.Routes() {
		var customColor color.Attribute

		switch value.Method {
		case "GET":
			customColor = color.FgGreen

		case "POST":
			customColor = color.FgHiYellow

		case "PUT":
			customColor = color.FgBlue

		case "PATCH":
			customColor = color.FgMagenta

		case "DELETE":
			customColor = color.FgHiRed

		default:
			continue
		}

		color.New(customColor, color.Bold).Print(fmt.Sprintf("\n\t%s : %s", value.Path, value.Method))
	}

	fmt.Println()
	color.New(color.BgHiBlack, color.Bold).Println("                                                    ")
}
