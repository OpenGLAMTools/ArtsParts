package main

import (
	"reflect"
	"testing"

	artsparts "github.com/OpenGLAMTools/ArtsParts"
)

var inst1 = &artsparts.Institution{
	ID:     "inst1",
	Name:   "Institution1",
	Admins: []string{"alice", "bob"},
}

var inst2 = &artsparts.Institution{
	ID:     "inst2",
	Name:   "Institution2",
	Admins: []string{"alice", "cindy"},
}

var institutions = artsparts.Institutions{inst1, inst2}

func TestAdmin_Institutions(t *testing.T) {
	adm := NewAdmin(institutions)
	type args struct {
		twitterName string
	}
	tests := []struct {
		name string
		a    *Admin
		args args
		want artsparts.Institutions
	}{
		{
			"Return 2 institutions for alice",
			adm,
			args{"alice"},
			artsparts.Institutions{inst1, inst2},
		},
		{
			"Return 1 institution for bob",
			adm,
			args{"bob"},
			artsparts.Institutions{inst1},
		},
		{
			"Return nothing for myuser",
			adm,
			args{"myuser"},
			artsparts.Institutions{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Institutions(tt.args.twitterName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Admin.Institutions() = %v, want %v", got, tt.want)
			}
		})
	}
}
