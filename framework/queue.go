package framework

type Queue struct {
	list    []Song
	current *Song
	running bool
}