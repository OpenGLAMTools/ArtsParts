package main

import "github.com/gorilla/mux"

var conf Conf

func init() {
	conf = loadConf()
}

func main() {

}

func makeRoutes() *mux.Router {
	return mux.NewRouter()
}
