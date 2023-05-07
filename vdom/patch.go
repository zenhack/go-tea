package vdom

type Patch interface {
	// Apply the patch to a DOM node (orig). The first argument
	// is the parent of the dom node being patched. The return
	// value is the new node, which may or may not be the same
	// as the old node.
	Patch(parent, orig DomNode) DomNode
}

type NopPatch struct{}

func (NopPatch) Patch(parent, orig DomNode) DomNode {
	return orig
}

type ReplacePatch struct {
	Replacement VNode
}

type ModifyPatch struct {
	Attrs    AttrsPatch
	Events   EventsPatch
	Children ChildPatch
}

type EventsPatch struct {
	Add    map[string]EventHandler
	Remove []string
}

type AttrsPatch struct {
	Add    map[string]string
	Remove []string
}

type ChildPatch struct {
	Common []Patch
	Drop   int
	Append []VNode
}
