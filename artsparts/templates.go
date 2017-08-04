package main

import (
	artsparts "github.com/OpenGLAMTools/ArtsParts"
)

// TemplateData defines the default values for the templates
type TemplateData struct {
	JSFiles  []string
	CSSFiles []string
	JQuery   bool
	VueJS    bool
	Title    string
	User     string
	Admin    bool
	Vars     map[string]string
	Pages    []Page
	Timeline []*artsparts.Artwork
	Pagenr   int
	Session  map[string]string
}

// AddJS adds a string to the JSFiles
func (td *TemplateData) AddJS(s string) {
	td.JSFiles = append(td.JSFiles, s)
}

// AddCSS adds a string to the CSSFiles
func (td *TemplateData) AddCSS(s string) {
	td.CSSFiles = append(td.CSSFiles, s)
}
