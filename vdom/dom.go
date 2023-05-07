//go:build js

package vdom

import (
	"syscall/js"
)

var document = js.Global().Get("document")

type DomNode struct {
	Value js.Value
}

type attributes struct {
	Value js.Value
}

func (a attributes) setNamed(k, v string) {
	attr := document.Call("createAttribute", k)
	attr.Set("value", v)
	a.Value.Call("setNamedItem", attr)
}

func (n DomNode) attributes() attributes {
	return attributes{Value: n.Value.Get("attributes")}
}

func (n DomNode) appendChild(child DomNode) {
	n.Value.Call("appendChild", child.Value)
}

func (n DomNode) setEventListener(name string, h EventHandler) {
	n.Value.Set("on"+name, js.FuncOf(func(this js.Value, args []js.Value) any {
		return (*h)(Event{Value: args[0]})
	}))
}

func (n DomNode) clearEventListener(name string) {
	v := n.Value.Get("on" + name)
	n.Value.Set("on"+name, js.Undefined())
	js.Func{Value: v}.Release()
}

type Event struct {
	Value js.Value
}

func (ve VElem) ToDomNode() DomNode {
	n := DomNode{Value: document.Call("createElement", ve.Tag)}
	if ve.Attrs != nil {
		attrs := n.attributes()
		for k, v := range ve.Attrs {
			attrs.setNamed(k, v)
		}
	}
	for k, h := range ve.Events {
		n.setEventListener(k, h)
	}
	for _, kid := range ve.Children {
		n.appendChild(kid.ToDomNode())
	}
	return n
}

func (vt VText) ToDomNode() DomNode {
	return DomNode{Value: document.Call("createTextNode", string(vt))}
}

func (p ReplacePatch) Patch(parent, orig DomNode) DomNode {
	newNode := p.Replacement.ToDomNode()
	parent.Value.Call("replaceChild", newNode.Value, orig.Value)
	return newNode
}

func (p ModifyPatch) Patch(parent, orig DomNode) DomNode {
	p.Attrs.patch(orig)
	p.Events.patch(orig)
	p.Children.patch(orig)
	return orig
}

func (p AttrsPatch) patch(n DomNode) {
	attrs := n.attributes()
	for k, v := range p.Add {
		attrs.setNamed(k, v)
	}
	for _, name := range p.Remove {
		n.Value.Call("removeAttribute", name)
	}
}

func (p EventsPatch) patch(n DomNode) {
	for _, v := range p.Remove {
		n.clearEventListener(v)
	}
	for k, v := range p.Add {
		n.setEventListener(k, v)
	}
}

func (cp ChildPatch) patch(n DomNode) {
	for i, p := range cp.Common {
		oldValue := n.Value.Get("childNodes").Index(i)
		childNode := DomNode{Value: oldValue}
		p.Patch(n, childNode)
	}
	firstDelete := len(cp.Common)
	for i := firstDelete; i < firstDelete+cp.Drop; i++ {
		child := n.Value.Get("childNodes").Index(firstDelete)
		js.Global().Get("console").Call("log", child)
		n.Value.Call("removeChild", child)
	}
	for _, child := range cp.Append {
		n.appendChild(child.ToDomNode())
	}
}
