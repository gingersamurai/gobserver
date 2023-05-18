package usecase

import (
	"gobserver/internal/entity"
	"os"
	"time"
)

type Storage interface {
	AddFileChange(fileChange entity.FileChange) (int64, error)
	AddCommandExecution(commandExecution entity.CommandExecution) (int64, error)
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

func (ti *TargetInteractor) Run() error {
	changesCh, err := ti.watcher.Listen(ti.target.Path, ti.target.ExcludeRegex, ti.target.ExcludeRegex)
	if err != nil {
		return err
	}

	for change := range changesCh {
		_, err := ti.storage.AddFileChange(entity.FileChange{
			ProjectPath: ti.target.Path,
			FilePath:    change.Path,
			Operation:   change.Operation,
			Timestamp:   time.Now(),
		})
		if err != nil {
			return err
		}
		for _, command := range ti.target.Commands {

			_, stderr, exitCode, err := ti.commandRunner.RunCommand(command, ti.target.Path)
			if err != nil {
				break
			}
			_, err = ti.storage.AddCommandExecution(entity.CommandExecution{
				ProjectPath: ti.target.Path,
				Command:     command,
				ExitCode:    exitCode,
				Timestamp:   time.Now(),
			})
			if err != nil {
				break
			}
			err = os.WriteFile(ti.target.LogFile, []byte(stderr), 0666)
			if err != nil {
				break
			}
			if exitCode != 0 {
				break
			}
		}
	}
	return nil
}
