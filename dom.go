//go:build js

package vdom

import "syscall/js"

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

func (n DomNode) addEventListener(name string, h EventHandler) {
	n.Value.Call("addEventListener", name, js.FuncOf(func(this js.Value, args []js.Value) any {
		return (*h)(Event{Value: args[0]})
	}))
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
		n.addEventListener(k, h)
	}
	for _, kid := range ve.Children {
		n.appendChild(kid.ToDomNode())
	}
	return n
}

func (vt VText) ToDomNode() DomNode {
	return DomNode{Value: document.Call("createTextNode", string(vt))}
}

func (p ModifyPatch) Patch(n *DomNode) {
	p.Attrs.patch(n)
	p.Events.patch(n)
	p.Children.patch(n)
}

func (AttrsPatch) patch(n *DomNode) {
}

func (EventsPatch) patch(n *DomNode) {
}

func (cp ChildPatch) patch(n *DomNode) {
	for i, p := range cp.Common {
		oldValue := n.Value.Get("children").Index(i)
		childNode := DomNode{Value: oldValue}
		p.Patch(&childNode)
		if !childNode.Value.Equal(oldValue) {
			n.Value.Call("replaceChild", childNode.Value, oldValue)
		}
	}
	for i := len(cp.Common); i < cp.Drop; i++ {
		child := n.Value.Get("children").Index(i)
		n.Value.Call("removeChild", child)
	}
	for _, child := range cp.Append {
		n.appendChild(child.ToDomNode())
	}
}

// An Updater manages updates to a dom node; create one with
// NewUpdater(), send updated VNode values with Update(), and shut
// it down with Close().
//
// The zero value is not meaningful.
type Updater struct {
	// Channel on which updates are sent. close this to shut
	// down the goroutine managing updates.
	updates chan VNode
}

// Close shuts down the updater.
func (up Updater) Close() error {
	close(up.updates)
	return nil
}

// Update updates the value of the node. The update will happen asynchronously,
// and multiple rapid updates may be coalesced.
func (up Updater) Update(vnode VNode) {
	up.updates <- vnode
}

// Create an Updater managing updates to the node, and return a handle to it.
func NewUpdater(node DomNode) Updater {
	updates := make(chan VNode)
	go func() {
		parent := node.Value.Get("parentNode")
		var (
			patch    Patch
			oldVNode VNode
		)

		for {
			// TODO(perf): use window.requestAnimationFrame to reduce
			// unnecessary updates.
			vnode, ok := <-updates
			if !ok {
				return
			}
			if oldVNode == nil {
				patch = ReplacePatch{Replacement: vnode}
			} else {
				patch = oldVNode.Diff(vnode)
			}
			oldNode := node
			patch.Patch(&node)
			if !node.Value.Equal(oldNode.Value) {
				parent.Call("replaceChild", node.Value, oldNode.Value)
			}
		}
	}()
	return Updater{updates: updates}
}
