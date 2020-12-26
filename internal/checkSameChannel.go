package internal

import (
	"architkl/goMusicBot/framework"
	"errors"
)

// Check if user and bot are in the same voice channel
func CheckSameChannel(ctx framework.Context) error {

	Logging(ctx)

	// Check if connected
	if !ctx.MediaPlayer.IsConnected {
		ctx.ReplyEmbed("Oops!", "Bot not connected!", 0xEB5160)
		return errors.New("Player not connected")
	}

	// Look for the message sender in that guild's current voice states.
	var userID, botID, userVoiceChannelID, botVoiceChannelID string
	userID = ctx.Message.Author.ID

	if usr, err := ctx.Discord.User("@me"); err != nil {
		ctx.ReplyEmbed("Oops!", "Error encountered", 0xEB5160)
		return errors.New("Error obtaining account details")
	} else {
		botID = usr.ID
	}

	// Get voice channels for user and bot
	for _, vs := range ctx.Guild.VoiceStates {
		if vs.UserID == userID {
			userVoiceChannelID = vs.ChannelID
		}

		if vs.UserID == botID {
			botVoiceChannelID = vs.ChannelID
		}
	}

	if userVoiceChannelID != botVoiceChannelID {
		ctx.ReplyEmbed("Oops!", "You're not in my voice channel!", 0xEB5160)
		return errors.New("User and Bot not in same channel")
	}

	return nil
}