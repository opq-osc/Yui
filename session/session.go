package session

import (
	"github.com/charmbracelet/log"
	"github.com/mcoo/OPQBot/session"
	_ "github.com/mcoo/OPQBot/session/provider"
)

var (
	S *session.Manager
)

func init() {
	var err error
	S, err = session.NewManager("qq", 600)
	if err != nil {
		log.Fatal(err)
	}
	S.GC()
}
