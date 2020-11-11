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
	// load songs from file system to program
	file, err := os.OpenFile("./docs/keys.txt", os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
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

func (songList *IdListHandler) UpdateSongs(videoId, title string) error {
	// update id map (program)
	songList.IdList[videoId] = title

	// update keys.txt (file system)
	file, err := os.OpenFile("./docs/keys.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return err
	}

	// write new id and title
	if _, err := file.WriteString(videoId + "," + title + "\n"); err != nil {
		log.Println(err)
	}

	// close the file
	if err := file.Close(); err != nil {
		log.Println(err)
	}

	return nil
}