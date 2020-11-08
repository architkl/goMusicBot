package framework

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type Player struct {
	Queue     []Song
	Current   *Song
	IsRunning bool
}

func NewMediaPlayer() *Player {
	player := new(Player)
	player.Queue = make([]Song, 0)
	return player
}

func (player *Player) AddSongs(songs ...Song) {
	player.Queue = append(player.Queue, songs...)
}

func (player *Player) StartPlaying(session *discordgo.Session, guildID, authorID string) {

	// Find the guild for that channel.
	guild, err := session.State.Guild(guildID)
	if err != nil {
		// Could not find guild.
		log.Println("Could not find guild: ", err)
		return
	}

	// Look for the message sender in that guild's current voice states.
	var voiceChannelID string
	for _, vs := range guild.VoiceStates {
		if vs.UserID == authorID {
			voiceChannelID = vs.ChannelID
			break
		}
	}

	// Join the provided voice channel.
	vc, err := session.ChannelVoiceJoin(guildID, voiceChannelID, false, true)
	if err != nil {
		log.Println(err)
		return
	}

	player.IsRunning = true

	for idx := 0; idx < len(player.Queue); idx++ {
		song := player.Queue[idx]
		var buffer = make([][]byte, 0)
		err := loadSound(song.Id, &buffer)
		if err != nil {
			log.Println("Error loading sound: ", err)
			continue
		}

		playSound(vc, buffer)
	}

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	player.IsRunning = false
}