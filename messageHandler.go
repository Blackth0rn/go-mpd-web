package main

import (
	"github.com/fhs/gompd/mpd"
	"log"
)

type Response struct {
	CurrentSong string
	Volume      int
	IsPlaying   bool
}

func mpdMessageHandle(c *mpd.Client, m []byte) error {
	var err error
	switch string(m) {
	case "play":
		err = c.Play(-1)
		log.Print("play")
	case "stop":
		err = c.Stop()
		log.Print("stop")
	}
	return err
}
