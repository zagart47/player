package music

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
)

type Song struct {
	Name     string
	Duration int
}

type Playlist []Song

func TrimDuration(s interface{}) int {
	reg := regexp.MustCompile(".*: ")
	res := reg.ReplaceAllString(fmt.Sprint(s), "${1}")
	d, err := strconv.Atoi(res)
	if err != nil {
		log.Fatal("duration extract error")
	}
	return d
}
