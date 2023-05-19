package usecase

import (
	"context"
	"gobserver/internal/entity"
	"log"
	"os"
	"sync"
	"time"
)

type Storage interface {
	AddFileChange(ctx context.Context, fileChange entity.FileChange) (int64, error)
	AddCommandExecution(ctx context.Context, commandExecution entity.CommandExecution) (int64, error)
}

type Watcher interface {
	Listen(path string, includeRegex, excludeRegex []string) (<-chan entity.Event, error)
}

type CommandRunner interface {
	RunCommand(command, dir string) (string, string, int, error)
}

type TargetInteractor struct {
	target        Target
	storage       Storage
	watcher       Watcher
	commandRunner CommandRunner
}

func NewTargetInteractor(target Target, storage Storage,
	watcher Watcher, commandRunner CommandRunner) *TargetInteractor {
	return &TargetInteractor{
		target:        target,
		storage:       storage,
		watcher:       watcher,
		commandRunner: commandRunner,
	}
}

func (ti *TargetInteractor) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	changesCh, err := ti.watcher.Listen(ti.target.Path,
		ti.target.IncludeRegexp,
		ti.target.ExcludeRegexp)
	if err != nil {
		log.Println(err)
		return
	}

	for change := range changesCh {

		_, err := ti.storage.AddFileChange(context.Background(),
			entity.FileChange{
				ProjectPath: ti.target.Path,
				FilePath:    change.Path,
				Operation:   change.Operation,
				Timestamp:   time.Now(),
			})
		if err != nil {
			log.Println(err)
			return
		}

		for _, command := range ti.target.Commands {

			_, stderr, exitCode, err := ti.commandRunner.RunCommand(command, ti.target.Path)
			if err != nil {
				log.Println(err)
				return
			}

			_, err = ti.storage.AddCommandExecution(context.Background(),
				entity.CommandExecution{
					ProjectPath: ti.target.Path,
					Command:     command,
					ExitCode:    exitCode,
					Timestamp:   time.Now(),
				})
			if err != nil {
				log.Println(err)
				break
			}

			err = AppendToFile(ti.target.LogFile, stderr)
			if err != nil {
				log.Println(err)
				return
			}

			if exitCode != 0 {
				break
			}
		}
	}
}

func AppendToFile(path string, text string) error {
	if path == "" {
		return nil
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	_, err = file.WriteString(text)
	if err != nil {
		return err
	}

	return nil
}
