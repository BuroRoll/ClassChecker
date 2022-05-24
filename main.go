package main

import (
	"ClassChecker/repository"
	"ClassChecker/service"
	"fmt"
	"github.com/robfig/cron/v3"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	db := repository.GetDB()
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	fmt.Println("Начало процесса каждые 3 часа")

	jakartaTime, _ := time.LoadLocation("Europe/Moscow")
	scheduler := cron.New(cron.WithLocation(jakartaTime))

	defer scheduler.Stop()

	scheduler.AddFunc("*/1 * * * *", services.CheckClassEnd)
	scheduler.AddFunc("*/1 * * * *", services.CheckClasses)
	//scheduler.AddFunc("0 */3 * * *", services.CheckClasses)
	//scheduler.AddFunc("0 */3 * * *", services.CheckClassEnd)

	go scheduler.Start()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
