package internal

import (
	"../framework"
	"errors"
)

// Check if bot can join the user's voice channel
func CanPlay(ctx framework.Context) error {

	Logging(ctx)

	// Look for the message sender in that guild's current voice states.
	var userVoiceChannelID string
	for _, vs := range ctx.Guild.VoiceStates {
		if vs.UserID == ctx.Message.Author.ID {
			userVoiceChannelID = vs.ChannelID
			break
		}
	}

	// User not in a voice channel
	if userVoiceChannelID == "" {
		ctx.Reply("Please join a voice channel!")
		return errors.New("User not in voice channel")
	}

	// If not connected then connect to the user's channel
	if !ctx.MediaPlayer.IsConnected {
		return nil
	}

	// Look for the message sender in that guild's current voice states.
	var botID, botVoiceChannelID string
	if usr, err := ctx.Discord.User("@me"); err != nil {
		ctx.Reply("Error")
		return errors.New("Error obtaining account details")
	} else {
		botID = usr.ID
	}

	// Get voice channel for bot
	for _, vs := range ctx.Guild.VoiceStates {
		if vs.UserID == botID {
			botVoiceChannelID = vs.ChannelID
			break
		}
	}

	if userVoiceChannelID != botVoiceChannelID {
		ctx.Reply("You're not in my voice channel!")
		return errors.New("User and Bot not in same channel")
	}

	// Same voice channel
	return nil
}