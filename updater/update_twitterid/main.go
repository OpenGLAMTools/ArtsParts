package main

import (
	"fmt"
	"os"
	"path/filepath"

	"flag"

	"io/ioutil"

	"encoding/json"

	artsparts "github.com/OpenGLAMTools/ArtsParts"
	"gopkg.in/yaml.v2"
)

var confFile = flag.String("conf", ".conf.yml", "conf file")

type Conf struct {
	SourceFolder string `yaml:"source_folder,omitempty"`
}

func main() {
	flag.Parse()
	b, err := ioutil.ReadFile(*confFile)
	if err != nil {
		fmt.Println("error loading conf file", err)
	}
	c := Conf{}
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		fmt.Println("error unmarshaling conf", err)
	}
	err = filepath.Walk(c.SourceFolder, fwalk)
	if err != nil {
		fmt.Println("error with filewalk", err)
	}
}

func fwalk(p string, info os.FileInfo, err error) error {
	if filepath.Base(p) == artsparts.DataFileName {
		b, err := ioutil.ReadFile(p)
		if err != nil {
			fmt.Println("error reading file", err)
		}
		artw := &artsparts.Artwork{}
		err = json.Unmarshal(b, artw)
		if err != nil {
			fmt.Println("error unmarshaling of ", p, err)
		}

		artw.TweetIDString = fmt.Sprintf("%d", artw.TweetID)
		if artw.TweetIDString == "0" {
			artw.TweetIDString = ""
		}

		for _, part := range artw.Parts {
			part.TweetIDString = fmt.Sprintf("%d", part.TweetID)
			part.MediaIDString = fmt.Sprintf("%d", part.MediaID)
			fmt.Println(part.TweetID, part.TweetIDString)
		}
		b, err = artw.Marshal()
		if err != nil {
			fmt.Println("error marshaling artw", p, err)
		}
		err = ioutil.WriteFile(p, b, 0777)

	}
	return err
}
