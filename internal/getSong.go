package internal

import (
	"errors"
	"log"
)

func GetSong(videoId string) error {

	// Get song from youtube
	log.Println("Getting video from youtube: " + videoId)
	if err := GetVideo(videoId); err != nil {
		log.Println("AddSong(): " + err.Error())
		return errors.New("Couldn't get song")
	}

	// Convert to mp3
	log.Println("Converting video to mp3: " + videoId)
	if err := ConvertMp3(videoId); err != nil {
		log.Println("AddSong(): " + err.Error())
		return errors.New("Couldn't get audio")
	}

	// Convert to dca
	log.Println("Converting mp3 to dca: " + videoId)
	if err := ConvertDca(videoId); err != nil {
		log.Println("AddSong(): " + err.Error())
		return errors.New("Couldn't convert audio")
	}

	return nil
}