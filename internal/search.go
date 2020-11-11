package internal

import (
	"../pkg"
	"log"

	"google.golang.org/api/youtube/v3"
)

func Search(query string) (string, string) {

	client := pkg.GetClient()

	service, err := youtube.New(client)
	if err != nil {
		log.Println(err)
		return "", ""
	}

	// Make the API call to YouTube.
	call := service.Search.List([]string{"id,snippet"}).
		Q(query).
		MaxResults(5)
	response, err := call.Do()
	if err != nil {
		log.Println(err)
		return "", ""
	}

	// get the first video id
	var videoId, title string
	for _, item := range response.Items {
		if item.Id.Kind == "youtube#video" {
			videoId = item.Id.VideoId
			title = item.Snippet.Title
			break
		}
	}

	return videoId, title
}