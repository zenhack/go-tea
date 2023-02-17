package builder

import (
	"strconv"

	"zenhack.net/go/vdom"
)

type A map[string]string
type E map[string]vdom.EventHandler

func H(tag string, attrs A, events E, children ...vdom.VNode) vdom.VNode {
	for i, v := range children {
		if v == nil {
			panic("Child #" + strconv.Itoa(i) + " is nil")
		}
	}
	return &vdom.VElem{
		Tag:      tag,
		Attrs:    map[string]string(attrs),
		Events:   map[string]vdom.EventHandler(events),
		Children: children,
	}
}

func T(text string) vdom.VNode {
	return vdom.VText(text)
}
