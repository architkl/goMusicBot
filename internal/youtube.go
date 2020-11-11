package internal

import (
	"io"
	"log"
	"os"

	"github.com/kkdai/youtube"
)

func Get(videoId string) error {
	client := youtube.Client{}

	video, err := client.GetVideo(videoId)
	if err != nil {
		log.Println(err)
		return err
	}

	resp, err := client.GetStream(video, &video.Formats[0])
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create("./songs/" + videoId + ".mp4")
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}