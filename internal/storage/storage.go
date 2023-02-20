package storage

import (
	"bytes"
	"fmt"
	"github.com/tcolgate/mp3"
	"io"
	"log"
	"os"
)

type File struct {
	Name     string
	Path     string
	Duration int
	Buffer   *bytes.Buffer
}

func NewFile(name string) File {
	dir := "files"
	pathToSave := fmt.Sprintf("%s/%s", dir, name)
	return File{
		Name:     name,
		Path:     pathToSave,
		Duration: 0,
		Buffer:   &bytes.Buffer{},
	}
}

func (f File) CheckDuration() {
	var t float64
	open, err := os.Open(f.Path)
	if err != nil {
		log.Fatal(err)
	}
	d := mp3.NewDecoder(open)
	var s mp3.Frame
	var skipped int

	for {
		if err := d.Decode(&s, &skipped); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			return
		}
		t += s.Duration().Seconds()
	}
	f.Duration = int(t)
}
