package artsparts

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestApp_GetTimeline(t *testing.T) {
	app, _ := NewApp(filepath.Join("test"))
	pic1, _ := app.GetArtwork("inst1", "coll1", "pic1")
	pic2, _ := app.GetArtwork("inst1", "coll1", "pic2")
	type args struct {
		filter string
	}
	tests := []struct {
		name    string
		a       *App
		args    args
		want    []*Artwork
		wantErr bool
	}{
		{
			"find all artworks",
			app,
			args{""},
			[]*Artwork{pic1, pic2},
			false,
		},
		{
			"find one artworks",
			app,
			args{"/*/pic2"},
			[]*Artwork{pic2},
			false,
		},
		{
			"wrong pattern",
			app,
			args{"*/pic2"},
			[]*Artwork{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.GetTimeline(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("App.GetTimeline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("App.GetTimeline() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestApp_AdminInstitutions(t *testing.T) {
	app, _ := NewApp(filepath.Join("test"))
	inst1, _ := app.GetInstitution("inst1")
	type args struct {
		twitterName string
	}
	tests := []struct {
		name string
		a    *App
		args args
		want Institutions
	}{
		{
			"find inst1 for user1",
			app,
			args{"user1"},
			Institutions{inst1},
		},
		{
			"find nothing for user3",
			app,
			args{"user3"},
			Institutions{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.AdminInstitutions(tt.args.twitterName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("App.AdminInstitutions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewInstitution(t *testing.T) {
	_, err := NewInstitution("NotExistPath")
	if err == nil {
		t.Error("Expect error for NotExistPath")
	}
	inst1, _ := NewInstitution(filepath.Join("test", "inst1"))

	if inst1.ID != "inst1" {
		t.Error("institution needs an ID")
	}
	coll1, _ := NewCollection(filepath.Join("test", "inst1", "coll1"), inst1)
	coll2, _ := NewCollection(filepath.Join("test", "inst1", "coll2"), inst1)
	tests := []struct {
		coll *Collection
	}{
		{coll1},
		{coll2},
	}
	for _, tt := range tests {
		if !reflect.DeepEqual(*inst1.Collections[tt.coll.ID], *tt.coll) {
			t.Errorf(
				"collection not loaded correct:\nExp:\n%#v\nGot:\n%#v\n",
				*tt.coll,
				*inst1.Collections[tt.coll.ID],
			)
		}
	}
}

func TestNewCollection(t *testing.T) {
	_, err := NewCollection("NotExistingPath", nil)
	if err == nil {
		t.Error("Expect error for NotExistingPath")
	}
	coll1, _ := NewCollection(filepath.Join("test", "inst1", "coll1"), nil)

	pic1, err := NewArtwork(filepath.Join("test", "inst1", "coll1", "pic1"), coll1)
	if err != nil {
		t.Error("NewArtwork returns an error: ", err)
	}
	if !reflect.DeepEqual(*coll1.Artworks[0], *pic1) {
		t.Errorf(
			"Artwork is not loaded into collection\nExp:\n%#v\nGot:\n%#v",
			*pic1,
			*coll1.Artworks[0],
		)
	}

}

func TestNewCollectionID(t *testing.T) {
	coll2, err := NewCollection(filepath.Join("test", "inst1", "coll2"), nil)
	// conf file does not have a id definition so the
	// id should be created from the path
	if err != nil {
		t.Error("should not return an error", err)
	}
	if coll2.ID != "coll2" {
		t.Errorf("should create an error\n Exp: coll2\nGot: %s", coll2.ID)
	}
}

func TestCollection_GetArtwork(t *testing.T) {
	coll1, _ := NewCollection(filepath.Join("test", "inst1", "coll1"), nil)
	type args struct {
		artwID string
	}
	tests := []struct {
		name  string
		coll  *Collection
		args  args
		want  *Artwork
		want1 bool
	}{
		{
			"search pic1",
			coll1,
			args{"pic1"},
			coll1.Artworks[0],
			true,
		},
		{
			"search not existing artwork",
			coll1,
			args{"picNotExist"},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.coll.GetArtwork(tt.args.artwID)
			if got != tt.want {
				t.Error("Pointers are not the same")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collection.GetArtwork() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Collection.GetArtwork() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

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
func Test_loadConf(t *testing.T) {
	type args struct {
		filePath string
		out      interface{}
	}
	tests := []struct {
		name    string
		target  string
		args    args
		wantErr bool
		expect  interface{}
	}{
		{
			"loading default conf in institution",
			"Institution",
			args{
				"default.conf.yml",
				&Institution{},
			},
			false,
			Institution{
				ID:          "new_id",
				Name:        "The name which is displayed",
				Description: "A short description about everything",
				License:     "MIT Open Source License",
				Order:       10,
				Admins:      []string{"user1", "user2"},
			},
		},
		{
			"loading default conf in collection",
			"Collection",
			args{
				"default.conf.yml",
				&Collection{},
			},
			false,
			Collection{
				ID:          "new_id",
				Name:        "The name which is displayed",
				Description: "A short description about everything",
				License:     "MIT Open Source License",
				Order:       10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := loadConf(tt.args.filePath, tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf("loadConf() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.target == "Institution" {
				got := tt.args.out.(*Institution)
				if !reflect.DeepEqual(*got, tt.expect) {
					t.Errorf("Expect:\n%#v\nGot:\n%#v", tt.expect, *got)
				}
			}
			if tt.target == "Collection" {
				got := tt.args.out.(*Collection)
				if !reflect.DeepEqual(*got, tt.expect) {
					t.Errorf("Expect:\n%#v\nGot:\n%#v", tt.expect, *got)
				}
			}

		})
	}
}
