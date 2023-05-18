package entity

import "time"

type CommandExecution struct {
	Id          int64
	ProjectPath string
	Command     string
	ExitCode    int
	Timestamp   time.Time
}

type FileChange struct {
	Id          int64
	ProjectPath string
	FilePath    string
	Operation   string
	Timestamp   time.Time
}
