package framework

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// PlaySound plays the current buffer to the provided channel.
func playSound(vc *discordgo.VoiceConnection, buffer [][]byte) {

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(1000 * time.Millisecond)

	// Start speaking.
	vc.Speaking(true)

	// Send the buffer data.
	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(1000 * time.Millisecond)
}