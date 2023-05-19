package main

import (
	"gobserver/internal/config"
	"gobserver/internal/storage/postgres_storage"
	"gobserver/internal/usecase"
	"gobserver/internal/watcher_wrap"
	"gobserver/pkg/closer"
	"gobserver/pkg/command_runner"
	"log"
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

	appCloser := closer.NewCloser(appConfig.ShutdownTimeout)

	appStorage, err := postgres_storage.NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}
	appCloser.Add(appStorage.Shutdown)
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
	appCloser.Run()
}
