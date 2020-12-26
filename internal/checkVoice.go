package internal

import (
	"architkl/goMusicBot/framework"
	"errors"
)

func CheckVoice(ctx framework.Context) error {

	guild := ctx.Guild

	// Look for the message sender in that guild's current voice states.
	for _, vs := range guild.VoiceStates {
		if vs.UserID == ctx.Message.Author.ID {
			return nil
		}
	}

	return errors.New("User not in voice channel")
}