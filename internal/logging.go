package internal

import (
	"architkl/goMusicBot/framework"
	"log"
)

func Logging(ctx framework.Context) error {
	log.Println(ctx.Message.Author.Username + " : " + ctx.Message.Content)

	return nil
}