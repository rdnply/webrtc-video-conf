package app

import (
	"github.com/rdnply/webrtc-video-conf/app/rtc"
	"html/template"
)

type App struct {
	connections *rtc.Connections
	templates   *templates
}

func New() *App {
	return &App{
		connections: rtc.New(),
		templates: readTemplates(),
	}
}

type templates struct {
	login *template.Template
	index *template.Template
}

func readTemplates() *templates {
	return &templates{
		login: readTemplate("login"),
		index: readTemplate("index"),
	}
}

func readTemplate(name string) *template.Template {
	path := "./static/html/"

	t := template.Must(template.New(name + ".html").
		ParseGlob(path + name + ".html"))

	return t
}
