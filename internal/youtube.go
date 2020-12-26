package internal

import (
	"io"
	"log"
	"os"
	"regexp"

	"github.com/kkdai/youtube/v2"
)

// Get videoID, title and duration from url or id
func GetMetaData(query string) (string, string, string) {

	log.Println("Getting data for: " + query)
	var videoId string
	// Check if a url is entered
	if matched, err := regexp.Match(`^((?:https?:)?\/\/)?((?:www|m)\.)?((?:youtube\.com|youtu.be))(\/(?:[\w\-]+\?v=|embed\/|v\/)?)([\w\-]+)(\S+)?$`, []byte(query)); err != nil {
		log.Println("GetMetaData(): " + err.Error())
		return "", "", ""
	} else if !matched {
		// search for the song online
		videoId = Search(query)
	} else {
		// store url as id
		videoId = query
	}

	client := youtube.Client{}

	video, err := client.GetVideo(videoId)
	if err != nil {
		log.Println("GetMetaData(): " + err.Error())
		return "", "", ""
	}

	return video.ID, video.Title, video.Duration.String()
}

// Get the video
func GetVideo(videoId string) error {

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