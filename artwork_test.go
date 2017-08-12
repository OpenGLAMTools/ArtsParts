package artsparts

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewArtworkEnsureConf(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "artsparts_test")
	defer os.RemoveAll(tmpDir)
	if err != nil {
		t.Error("can not create tmp dir:", err)
	}
	artw, err := NewArtwork(tmpDir, nil)
	if err != nil {
		t.Error("error creating artwork:", err)
	}
	filename := artw.dataFilePath()
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Error("data file does not exist")
	}
}

func TestArtworkPath(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "artsparts_test")
	defer os.RemoveAll(tmpDir)
	if err != nil {
		t.Error("can not create tmp dir:", err)
	}
	artw, err := NewArtwork(tmpDir, nil)
	if err != nil {
		t.Error("error creating artwork:", err)
	}
	if artw.Path() != tmpDir {
		t.Error("Artwork need to return the path")
	}
}

func TestArtworkImgFile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "artsparts_test")
	if err != nil {
		t.Error("can not create tmp dir:", err)
	}
	defer os.RemoveAll(tmpDir)
	b := []byte{}
	ioutil.WriteFile(filepath.Join(tmpDir, "f1.txt"), b, 0777)
	ioutil.WriteFile(filepath.Join(tmpDir, "img.jpg"), b, 0777)
	artw, err := NewArtwork(tmpDir, nil)
	if err != nil {
		t.Error("error creating artwork:", err)
	}
	imgFile, err := artw.ImgFile()
	if err != nil {
		t.Error(err)
	}
	if imgFile != "img.jpg" {
		t.Error("Expect img.jpg got:", imgFile)
	}
}

func TestArtworkIsAdminUser(t *testing.T) {
	app, _ := NewApp(filepath.Join("test"))
	artw, _ := app.GetArtwork("inst1", "coll1", "pic1")
	tests := []struct {
		user   string
		expect bool
	}{
		{"user1", true},
		{"abc", false},
	}
	for _, tt := range tests {
		got := artw.IsAdminUser(tt.user)
		if got != tt.expect {
			t.Errorf("IsAdminUser returns wrong value\nInput:%s\nExp: %v  Got: %v", tt.user, tt.expect, got)
		}
	}
}

func TestCopyArtwork(t *testing.T) {
	app, _ := NewApp(filepath.Join("test"))
	pic1, _ := app.GetArtwork("inst1", "coll1", "pic1")
	pic2, _ := app.GetArtwork("inst1", "coll1", "pic2")
	type args struct {
		dst    *Artwork
		src    *Artwork
		tagval string
	}
	type exp struct {
		Timestamp     string
		TweetIDString string
		Shortlink     string
		Meta          MetaData
	}
	tests := []struct {
		name string
		args args
		exp  exp
	}{
		{
			"t",
			args{
				&Artwork{},
				pic1,
				"admin",
			},
			exp{
				"201708011155",
				"",
				"",
				MetaData{Artist: "Artistval", Date: "Dateval", Link: "http://test.com"},
			},
		},
		{
			"t",
			args{
				pic2,
				pic1,
				"admin",
			},
			exp{
				"201708011155",
				"pic2",
				"/s/2",
				MetaData{Artist: "Artistval", Date: "Dateval", Link: "http://test.com"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CopyArtwork(tt.args.dst, tt.args.src, tt.args.tagval)
			// Check important fields
			dst := tt.args.dst
			exp := tt.exp
			assert.Equal(t, exp.Timestamp, dst.Timestamp)
			assert.Equal(t, exp.TweetIDString, dst.TweetIDString)
			assert.Equal(t, exp.Shortlink, dst.ShortLink)
			assert.Equal(t, exp.Meta, dst.Meta)

		})
	}
}
