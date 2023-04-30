//go:build !tinygo.wasm

package S

import (
	"Yui/opq"
	"github.com/opq-osc/OPQBot/v2/apiBuilder"
	"github.com/opq-osc/OPQBot/v2/events"
)

func GetApi(event events.IEvent) apiBuilder.IMainFunc {
	return apiBuilder.New(opq.Url, event.GetCurrentQQ())
}
