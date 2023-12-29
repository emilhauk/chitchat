package templates

import (
	"github.com/emilhauk/chitchat/config"
	"html/template"
)

var (
	err       error
	Templates *template.Template
)

func init() {
	Templates, err = template.ParseGlob("./templates/**")
	if err != nil {
		config.Logger.Fatal().Err(err).Msg("Failed to parse templates")
	}
}
