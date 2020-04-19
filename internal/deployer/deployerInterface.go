package deployer

type DeployInterface interface {
	Prepare() error
	UpdateSource(git string) error
	RunPipe() error
	MakeLinks() error
	CleanUp(count int) error
}
