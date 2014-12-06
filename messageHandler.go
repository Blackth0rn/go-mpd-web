package main

import (
	"encoding/json"
	"github.com/fhs/gompd/mpd"
	"log"
	"strconv"
)

type cmdInput struct {
	Cmd   string
	Data  string
	Token int
}

type attrReturn struct {
	Cmd  string
	Attr mpd.Attrs
}

type initReturn struct {
	Cmd  string
	Attr mpd.Attrs
}

func mpdMessageHandle(c *mpd.Client, m []byte) ([]byte, error, int) {
	var err error
	var jsonReturn []byte
	var input cmdInput
	var broadcast int

	if err := json.Unmarshal(m, &input); err != nil {
		return m, err, input.Token
	} else {
		broadcast = input.Token
		switch input.Cmd {
		case "play":
			if err = c.Play(-1); err == nil {
				var attrs mpd.Attrs
				attrs, err = c.Status()
				if err == nil {
					jsonReturn, err = json.Marshal(attrReturn{input.Cmd, attrs})
					broadcast = 0
				}
			}
		case "pause":
			if err = c.Pause(true); err == nil {
				var attrs mpd.Attrs
				attrs, err = c.Status()
				if err == nil {
					jsonReturn, err = json.Marshal(attrReturn{input.Cmd, attrs})
					broadcast = 0
				}
			}
		case "stop":
			if err = c.Stop(); err == nil {
				var attrs mpd.Attrs
				attrs, err = c.Status()
				if err == nil {
					jsonReturn, err = json.Marshal(attrReturn{input.Cmd, attrs})
					broadcast = 0
				}
			}
		case "init":
			var attrs mpd.Attrs
			attrs, err = c.Status()
			if err == nil {
				jsonReturn, err = json.Marshal(initReturn{input.Cmd, attrs})
			}
		case "setVolume":
			if volume, err := strconv.ParseInt(input.Data, 0, 0); err == nil {
				if err = c.SetVolume(int(volume)); err == nil {
					var attrs mpd.Attrs
					attrs, err = c.Status()
					if err == nil {
						jsonReturn, err = json.Marshal(attrReturn{input.Cmd, attrs})
						broadcast = 0
					}
				}
			}
		}
	}
	log.Print("jsonReturn:", string(jsonReturn), " err:", err, " bcast:", broadcast)
	return jsonReturn, err, broadcast
}
