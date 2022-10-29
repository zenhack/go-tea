//go:build !js

package vdom

type DomNode struct {
	Value any
}

type Event struct {
	Value any
}

func (ve VElem) ToDomNode() DomNode {
	return DomNode{Value: ve}
}

func (vt VText) ToDomNode() DomNode {
	return DomNode{Value: vt}
}

func (ModifyPatch) Patch(*DomNode) {
}
