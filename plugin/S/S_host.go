//go:build !tinygo.wasm

package S

import (
	"github.com/opq-osc/OPQBot/v2/apiBuilder"
	"github.com/opq-osc/OPQBot/v2/events"
	"github.com/opq-osc/Yui/opq"
)

func GetApi(event events.IEvent) apiBuilder.IMainFunc {
	return apiBuilder.New(opq.Url, event.GetCurrentQQ())
}
