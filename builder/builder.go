// Package builder provides helpers for constructing VNodes
//
// A common pattern is to define shorthands for these in the package
// that does your app's rendering. For example:
//
// ```
// import (
// 	"zenhack.net/go/vdom"
//	"zenhack.net/go/vdom/builder"
// )
//
// type (
//	a = builder.A
//	e = builder.E
// )
//
// var (
// 	h = builder.H
//	t = builder.T
// )
//
// func view() vdom.VNode {
//	return h("a", a{"href": "/"}, nil, t("home"))
// }
// ```
package builder

import (
	"strconv"

	"zenhack.net/go/vdom"
)

// An A is a map of attribute names to values.
type A map[string]string

// An E is a map of event names to event handlers.
type E map[string]vdom.EventHandler

// H constructs a vdom.VElem node with the given tag, attributes, event handlers,
// and child nodes. Either attrs or events or both may be nil.
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

// T constructs a vdom.VText node with the given text.
func T(text string) vdom.VNode {
	return vdom.VText(text)
}
