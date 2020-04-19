package pkg

type ArgsInterface interface {
	CheckConfig() error
	ExecuteOnLocal() error
	GetRemoteCommand() string
}
