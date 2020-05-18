package pkg

// ArgsInterface for pipe types
type ArgsInterface interface {
	CheckConfig() error
	ExecuteOnLocal() error
	GetRemoteCommand() string
}
