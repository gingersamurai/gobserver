package usecase

type Target struct {
	Path          string
	IncludeRegexp []string
	ExcludeRegexp []string
	Commands      []string
	LogFile       string
}
