package vdom

type Patch interface {
	Patch(*DomNode)
}

type NopPatch struct{}

func (NopPatch) Patch(*DomNode) {
}

type ReplacePatch struct {
	Replacement VNode
}

func (p ReplacePatch) Patch(n *DomNode) {
	*n = p.Replacement.ToDomNode()
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
