package internal

import (
	"../framework"
	"errors"
)

func CheckVoice(ctx framework.Context) error {
	Logging(ctx)

	// Find the channel that the message came from.
	c, err := ctx.Discord.State.Channel(ctx.Message.ChannelID)
	if err != nil {
		// Could not find channel.
		return err
	}

	// Find the guild for that channel.
	g, err := ctx.Discord.State.Guild(c.GuildID)
	if err != nil {
		// Could not find guild.
		return err
	}

	// Look for the message sender in that guild's current voice states.
	for _, vs := range g.VoiceStates {
		if vs.UserID == ctx.Message.Author.ID {
			return nil
		}
	}

	return errors.New("User not in voice channel")
}