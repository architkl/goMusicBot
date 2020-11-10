package framework

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type (
	Song struct {
		Id       string
		Title    string
		Duration string
	}

	IdListHandler struct {
		IdList map[string]string
	}
)

func NewIdListHandler() *IdListHandler {
	return &IdListHandler{make(map[string]string)}
}

func (songList *IdListHandler) LoadSongs() {
	file, err := os.OpenFile("./docs/keys.txt", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, ",")

		songList.IdList[words[0]] = words[1]
	}

	if err := scanner.Err(); err != nil {
		log.Println(err.Error())
	}

	err = file.Close()
	if err != nil {
		log.Println(err)
	}
}