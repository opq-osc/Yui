package cron

import "github.com/robfig/cron/v3"

var (
	C = cron.New()
)

func init() {
	C.Start()
}
