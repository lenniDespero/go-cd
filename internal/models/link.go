package models

type Link struct {
	From string `mapstructure:"from"`
	To   string `mapstructure:"to"`
}
type LinkError string

func (e LinkError) Error() string {
	return string(e)
}

//Errors
const (
	NoLinkFromError LinkError = "no link from"
	NoLinkToError   LinkError = "no link to"
)

func (link Link) checkArgsConfig() error {
	if link.From == "" {
		return NoLinkFromError
	}
	if link.To == "" {
		return NoLinkToError
	}
	return nil
}
