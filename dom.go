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
	// TODO: event handlers
	for _, kid := range ve.Children {
		e.Call("appendChild", kid.Node.ToDomNode().Value)
	}
	return DomNode{Value: e}
}

func (vt VText) ToDomNode() DomNode {
	return DomNode{Value: document.Call("createTextNode", string(vt))}
}

func (ModifyPatch) Patch(n *DomNode) {
	panic("TODO")
}
