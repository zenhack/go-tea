//go:build js

package vdom

import "syscall/js"

var document = js.Global().Get("document")

type DomNode struct {
	Value js.Value
}

type Event struct {
	Value js.Value
}

func (ve VElem) ToDomNode() DomNode {
	e := document.Call("createElement", ve.Tag)
	if ve.Attrs != nil {
		eAttrs := e.Get("attributes")
		for k, v := range ve.Attrs {
			attr := document.Call("createAttribute", k)
			attr.Set("value", v)
			eAttrs.Call("setNamedItem", attr)
		}
	}
	for k, h := range ve.Events {
		e.Call("addEventListener", k, js.FuncOf(func(this js.Value, args []js.Value) any {
			return h(Event{Value: args[0]})
		}))
	}
	for _, kid := range ve.Children {
		e.Call("appendChild", kid.ToDomNode().Value)
	}
	return DomNode{Value: e}
}

func (vt VText) ToDomNode() DomNode {
	return DomNode{Value: document.Call("createTextNode", string(vt))}
}

func (ModifyPatch) Patch(n *DomNode) {
	panic("TODO")
}
