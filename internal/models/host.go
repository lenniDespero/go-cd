package models

type Host struct {
	Host string `mapstructure:"host"`
	User string `mapstructure:"user"`
}
