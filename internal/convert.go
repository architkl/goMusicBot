package internal

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
)

func ConvertMp3(videoId string) error {
	// convert mp4 to mp3
	cmd := exec.Command("ffmpeg", "-i", "./songs/"+videoId+".mp4", "-vn", "-acodec", "libmp3lame", "-ac", "2", "-ab", "160k", "-ar", "48000", "./songs/"+videoId+".mp3")

	// run cmd and check for errors
	if err := cmd.Run(); err != nil {
		log.Println("Error in converting to mp3")
		log.Println(err)
		return err
	}

	// clean up the mp3 file
	if err := os.Remove("./songs/" + videoId + ".mp4"); err != nil {
		log.Println(err)
	}

	return nil
}

func ConvertDca(videoId string) error {
	// cmd1 converts mp3 to intermediate format and pipes the output to stdout
	cmd1 := exec.Command("ffmpeg", "-i", "./songs/"+videoId+".mp3", "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	// cmd2 reads from stdin and converts the output of cmd1 to proper dca format
	cmd2 := exec.Command("dca")

	// create dca file for the audio
	f, err := os.OpenFile("./songs/"+videoId+".dca", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println("error opening file: %v", err)
		return err
	}
	defer f.Close()

	// create a pipe to which cmd1 will write to and cmd2 will read from
	pipeRead, pipeWrite := io.Pipe()
	cmd1.Stdout = pipeWrite
	cmd2.Stdin = pipeRead

	// output of cmd2 will be the final converted dca
	var audioData bytes.Buffer
	cmd2.Stdout = &audioData

	// check if cmd1 starts successfully
	if err := cmd1.Start(); err != nil {
		log.Println(err)
		return err
	}

	// check if cmd2 starts successfully
	if err := cmd2.Start(); err != nil {
		log.Println(err)
		return err
	}

	// wait for cmd1 to write output to pipe
	if err := cmd1.Wait(); err != nil {
		log.Println(err)
		return err
	}

	// close writer end of pipe
	if err := pipeWrite.Close(); err != nil {
		log.Println(err)
		return err
	}

	// wait for cmd2 to read from the pipe
	if err := cmd2.Wait(); err != nil {
		log.Println(err)
		return err
	}

	// write output of cmd2 to the audio file
	if _, err = f.Write(audioData.Bytes()); err != nil {
		log.Println(err)
		return err
	}

	// clean up the mp3 file
	if err := os.Remove("./songs/" + videoId + ".mp3"); err != nil {
		log.Println(err)
	}

	return nil
}