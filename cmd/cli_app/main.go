package main

import (
	"gobserver/internal/config"
	"gobserver/internal/storage/postgres_storage"
	"gobserver/internal/usecase"
	"gobserver/internal/watcher_wrap"
	"gobserver/pkg/command_runner"
	"log"
	"time"
)

const (
	configFilePath = "./"
	configFileName = "config"
)

func main() {

	appConfig, err := config.NewConfig(configFilePath, configFileName)
	if err != nil {
		log.Fatal(err)
	}
	appStorage, err := postgres_storage.NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}
	commandRunner := command_runner.NewCommandRunner()
	for _, target := range appConfig.Targets {
		watcher, err := watcher_wrap.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		appTargetInteractor := usecase.NewTargetInteractor(usecase.Target{
			Path:         target.Path,
			IncludeRegex: target.IncludeRegexp,
			ExcludeRegex: target.ExcludeRegexp,
			Commands:     target.Commands,
			LogFile:      target.LogFile,
		}, appStorage, watcher, commandRunner)
		go appTargetInteractor.Run()
	}
	time.Sleep(10 * time.Hour)
}
