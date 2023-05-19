package main

import (
	"gobserver/internal/config"
	"gobserver/internal/storage/postgres_storage"
	"gobserver/internal/usecase"
	"gobserver/internal/watcher"
	"gobserver/pkg/closer"
	"gobserver/pkg/command_runner"
	"log"
	"os"
	"sync"
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

	var wg sync.WaitGroup
	wg.Add(len(appConfig.Targets))
	for _, target := range appConfig.Targets {

		targetWatcher, err := watcher.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		appCloser.Add(targetWatcher.Shutdown)

		appTargetInteractor := usecase.NewTargetInteractor(usecase.Target{
			Path:          target.Path,
			IncludeRegexp: target.IncludeRegexp,
			ExcludeRegexp: target.ExcludeRegexp,
			Commands:      target.Commands,
			LogFile:       target.LogFile,
		}, appStorage, targetWatcher, commandRunner)
		go appTargetInteractor.Run(&wg)
	}
	appCloser.Run()
	wg.Wait()
	os.Exit(0)
}
