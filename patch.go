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
	AddAttrs    map[string]string
	RemoveAttrs []string
	Events      map[string]EventHandler
	Children    []ChildPatch
}

type ChildPatch struct {
	Key   string
	Patch Patch
}
