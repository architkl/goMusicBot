package framework

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Player struct {
	Queue       []Song
	idx         int
	Current     *Song
	voice       *discordgo.VoiceConnection
	IsConnected bool
	IsPaused    bool
	skip        bool       // skip current song
	audioMutex  sync.Mutex // lock for playing audio
	queueMutex  sync.Mutex // lock for accessing queue
}

// create new player object
func NewMediaPlayer() *Player {
	player := new(Player)
	player.IsPaused = true
	player.Queue = make([]Song, 0)
	player.audioMutex.Lock()
	return player
}

// Connect to voice channel
func (player *Player) ConnectVoice(session *discordgo.Session, guild *discordgo.Guild, authorID string) {

	// check if already connected
	if player.IsConnected {
		log.Println("Player.ConnectVoice(): Already connected")
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
	vc, err := session.ChannelVoiceJoin(guild.ID, voiceChannelID, false, true)
	if err != nil {
		log.Println(err)
		return
	}

	// assign voice channel to player
	player.voice = vc

	// set the player to connected
	player.IsConnected = true
}

// Disconnect from voice channel
func (player *Player) DisconnectVoice() {

	// Check if not connected
	if !player.IsConnected {
		log.Println("Player.DisconnectVoice(): Not connected")
		return
	}

	player.IsConnected = false

	if player.IsPaused {
		// Resume so that playSound can return
		player.Resume()
	}
}

// Add songs to the player
func (player *Player) AddSongs(songs ...Song) {
	player.Queue = append(player.Queue, songs...)
}

// Start playing audio
func (player *Player) StartPlaying(session *discordgo.Session, guild *discordgo.Guild, authorID string) {

	// connect to the voice channel if not connected
	if !player.IsConnected {
		player.ConnectVoice(session, guild, authorID)
		if player.voice == nil {
			return
		}
	}

	// check if index is at end of queue
	if player.idx == len(player.Queue) {
		player.idx = 0
	}

	// Start speaking.
	player.voice.Speaking(true)
	player.audioMutex.Unlock()
	player.IsPaused = false

	// loop through the queue
	for idx := &player.idx; *idx < len(player.Queue); *idx++ {
		player.queueMutex.Lock()
		song := player.Queue[*idx]
		player.queueMutex.Unlock()
		var buffer = make([][]byte, 0)

		// load song data from file system
		err := loadSound(song.Id, &buffer)
		if err != nil {
			log.Println("Player.StartPlaying(): Error loading sound: ", err)
			continue
		}

		player.playSound(buffer)

		// check if player still connected
		if !player.IsConnected {
			break
		}
	}

	// Stop speaking
	player.voice.Speaking(false)
	player.IsPaused = true
	player.audioMutex.Lock()

	// Disconnect from the provided voice channel.
	player.voice.Disconnect()
	player.IsConnected = false
}

// playSound plays the current buffer to the provided channel.
func (player *Player) playSound(buffer [][]byte) {

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// Send the buffer data.
	for _, buff := range buffer {
		player.audioMutex.Lock()
		if !player.IsConnected {
			// return due to disconnection
			player.audioMutex.Unlock()
			return
		}

		if player.skip == true {
			// skip current song
			player.skip = false
			player.audioMutex.Unlock()
			return
		}

		player.voice.OpusSend <- buff
		player.audioMutex.Unlock()
	}

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)
}

// Resume playing
func (player *Player) Resume() error {

	// Check if already playing
	if !player.IsPaused {
		log.Println("Player.Play(): Already playing")
		return nil
	}

	if player.idx == len(player.Queue) {
		log.Println("Player.Play(): At end of queue")
		return errors.New("No song to play!")
	}

	player.voice.Speaking(true)
	player.IsPaused = false
	player.audioMutex.Unlock()
	return nil
}

// Pause by locking mutex
func (player *Player) Pause() {

	// Check if already paused
	if player.IsPaused {
		log.Println("Player.Pause(): Already paused")
		return
	}

	player.audioMutex.Lock()
	player.voice.Speaking(false)
	player.IsPaused = true
}

// Skips current song
func (player *Player) Skip() {

	// Check if at end of queue
	if player.idx == len(player.Queue) {
		log.Println("Player.Skip(): At end of queue")
		return
	}

	// Resume if paused
	if player.IsPaused {
		player.Resume()
	}

	player.audioMutex.Lock()
	player.skip = true
	player.audioMutex.Unlock()
}

// Return index of current song
func (player *Player) CurPosition() (int, error) {

	// Check if empty queue
	if len(player.Queue) == 0 {
		log.Println("Player.CurPosition(): At end of queue")
		return 0, errors.New("No song to play!")
	}

	return player.idx, nil
}

// Shuffle the queue current song onwards
func (player *Player) Shuffle() error {

	// Check if at end of queue
	if player.idx == len(player.Queue) {
		log.Println("Player.Shuffle(): At end of queue")
		return errors.New("")
	}

	// Divide queue on current song
	q1 := player.Queue[:player.idx+1]
	q2 := player.Queue[player.idx+1:]

	// Shuffle the second part
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(q2), func(i, j int) { q2[i], q2[j] = q2[j], q2[i] })

	// Merge them back
	player.queueMutex.Lock()
	player.Queue = append(q1, q2...)
	player.queueMutex.Unlock()

	return nil
}
