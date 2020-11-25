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
		IdList map[string]Song
	}
)

func NewIdListHandler() *IdListHandler {
	return &IdListHandler{make(map[string]Song)}
}

func (songList *IdListHandler) LoadSongs() {
	// Load songs from file system to program
	file, err := os.OpenFile("./docs/keys.txt", os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, ";")

		songList.IdList[words[0]] = Song{
			Id:       words[0],
			Title:    words[1],
			Duration: words[2],
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err.Error())
	}

	err = file.Close()
	if err != nil {
		log.Println(err)
	}
}

func (songList *IdListHandler) UpdateSongs(song Song) error {
	// Update id map (program)
	songList.IdList[song.Id] = song

	// Update keys.txt (file system)
	file, err := os.OpenFile("./docs/keys.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return err
	}

	// Write new id and title
	if _, err := file.WriteString(song.Id + ";" + song.Title + ";" + song.Duration + "\n"); err != nil {
		log.Println(err)
		return err
	}

	// Close the file
	if err := file.Close(); err != nil {
		log.Println(err)
	}

	return nil
}