package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	_ "project/docs"
	"project/helper"
	"project/infra"
	"project/routes"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
)

// @title Ecommerce Dashboard API
// @version 1.0
// @description Nothing.
// @termsOfService http://example.com/terms/
// @contact.name Team-1
// @contact.url https://academy.lumoshive.com/contact-us
// @contact.email lumoshive.academy@gmail.com
// @license.name Lumoshive Academy
// @license.url https://academy.lumoshive.com
// @host localhost:8080
// @schemes http
// @BasePath /
// @securityDefinitions.apikey token
// @in header
// @name token
func main() {
	migrateDb := flag.Bool("m", false, "use this flag to migrate database")
	seedDb := flag.Bool("s", false, "use this flag to seed database")
	flag.Parse()

	crn := cron.New()

	ctx, err := infra.NewServiceContext(*migrateDb, *seedDb)
	if err != nil {
		log.Fatal("can't init service context %w", err)
	}
	crn.AddFunc("* * * * *", helper.CronExcel(*migrateDb, *seedDb))
	crn.Start()

	if !shouldLaunchServer(*migrateDb, *seedDb) {
		return
	}

	srv := routes.NewRoutes(*ctx)

	go func() {
		// service connections
		log.Println("Listening and serving HTTP on", ctx.Cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	appContext, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(appContext); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching appContext.Done(). timeout of 5 seconds.
	select {
	case <-appContext.Done():
		log.Println("timeout of 3 seconds.")
	}
	log.Println("Server exiting")

}

func shouldLaunchServer(migrateDb bool, seedDb bool) bool {
	if migrateDb {
		return false
	}

	if seedDb {
		return false
	}

	return true
}
