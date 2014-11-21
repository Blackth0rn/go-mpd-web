package main

import (
	"encoding/json"
	"github.com/fhs/gompd/mpd"
	"log"
)

type Response struct {
	CurrentSong string
	Volume      int
	IsPlaying   bool
}

type Command struct {
	Type string
	Data string
	Err  []byte
}

func mpdMessageHandle(c *mpd.Client, m []byte) (Command, error) {
	var err error
	var cmd Command
	if err := json.Unmarshal(m, &cmd); err != nil {
		return cmd, err
	} else {
		log.Print(cmd.Type)
		var attrs mpd.Attrs
		switch string(cmd.Type) {
		case "play":
			err = c.Play(-1)
		case "stop":
			err = c.Stop()
		case "init":
			attrs, err = c.Status()
			if err == nil {
				var tmpJson []byte
				tmpJson, err = json.Marshal(&attrs)
				cmd.Data = string(tmpJson)
			}
		}

	}
	log.Print(string(cmd.Data))
	return cmd, err
}
