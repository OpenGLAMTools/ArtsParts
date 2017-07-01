package main

import artsparts "github.com/OpenGLAMTools/ArtsParts"

type templateData struct {
	JSFiles  []string
	CSSFiles []string
	JQuery   bool
	VueJS    bool
	Title    string
	User     string
	Admin    bool
	Timeline []*artsparts.Artwork
}
