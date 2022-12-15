package web

import (
	"bytes"
	"context"
	"html/template"

	"io/fs"
)

type TemplateEngine interface {
	// Render 渲染页面
	// data 是渲染页面所需要的数据
	Render(ctx context.Context, tplName string, data any) ([]byte, error)
}

type GoTemplateEngine struct {
	T *template.Template
}

func (g *GoTemplateEngine) Render(ctx context.Context, tplName string, data any) ([]byte, error) {
	res := &bytes.Buffer{}
	err := g.T.ExecuteTemplate(res, tplName, data)
	return res.Bytes(), err
}

func (g *GoTemplateEngine) LoadFromGlob(pattern string) error {
	var err error
	g.T, err = template.ParseGlob(pattern)
	return err
}

func (g *GoTemplateEngine) LoadFromFiles(files ...string) error {
	var err error
	g.T, err = template.ParseFiles(files...)
	return err
}

func (g *GoTemplateEngine) LoadFromFS(fs fs.FS, pattern ...string) error {
	var err error
	g.T, err = template.ParseFS(fs, pattern...)
	return err
}
