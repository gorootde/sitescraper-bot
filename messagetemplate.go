package main

import (
	"bytes"
	"html/template"
	"path"

	"github.com/sirupsen/logrus"
)

type MessageTemplate struct {
	templateFile string
}

type context struct {
	Fields map[string]string
}

func NewMessageTemplate(templateFile string) *MessageTemplate {
	return &MessageTemplate{templateFile}
}

func (t *MessageTemplate) render(fields map[string]string) string {
	context := context{fields}

	tpl := template.Must(template.New(path.Base(t.templateFile)).ParseFiles(t.templateFile))

	result := new(bytes.Buffer)
	err := tpl.Execute(result, context)
	if err != nil {
		logrus.Panicf("Error rendering template: %s", err)
	}
	return result.String()
}
