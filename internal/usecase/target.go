package usecase

type Target struct {
	Path         string
	IncludeRegex []string
	ExcludeRegex []string
	Commands     []string
	LogFile      string
}
