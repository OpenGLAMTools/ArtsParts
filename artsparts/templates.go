package main

import artsparts "github.com/OpenGLAMTools/ArtsParts"

type templateData struct {
	JSFiles  []string
	CSSFiles []string
	JQuery   bool
	VueJS    bool
	User     string
	Timeline []*artsparts.Artwork
}
