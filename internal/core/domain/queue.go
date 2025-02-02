package domain

type Queue struct {
	Name     string
	Messages chan string
}
