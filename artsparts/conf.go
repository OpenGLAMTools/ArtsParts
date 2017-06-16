package main

type Conf struct {
	TwitterKey    string
	TwitterSecret string
}

func loadConf() Conf {
	return Conf{}
}
