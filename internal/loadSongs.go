package internal

/*
import (
	"bufio"
	"log"
	"strings"
)

func LoadSongs() map[string]string {
	file, err := os.OpenFile("./docs/keys.txt", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
		return
	}

	scanner := bufio.NewScanner(file)
	var list = make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, ",")

		list[words[0]] = words[1]
	}

	if err := scanner.Err(); err != nil {
		log.Println(err.Error())
	}

	err = file.Close()
	if err != nil {
		log.Println(err)
	}

	return list
}
*/