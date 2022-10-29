package vdom

type VNode interface {
	Diff(other VNode) Patch
	ToDomNode() DomNode
}

type VElem struct {
	Tag      string
	Attrs    map[string]string
	Events   EventHandler
	Children []Child
}

type Child struct {
	Key  string
	Node VNode
}

type VText string

type EventHandler = func(Event)

func (ve *VElem) Diff(other VNode) Patch {
	otherElem, ok := other.(*VElem)
	if !ok {
		return ReplacePatch{Replacement: other}
	}
	if ve == otherElem {
		return NopPatch{}
	}
	if ve.Tag != otherElem.Tag {
		return ReplacePatch{Replacement: other}
	}

	patch := ModifyPatch{}
	for k, _ := range ve.Attrs {
		if _, ok := otherElem.Attrs[k]; !ok {
			patch.RemoveAttrs = append(patch.RemoveAttrs, k)
		}
	}
	patch.AddAttrs = make(map[string]string)
	for k, v := range otherElem.Attrs {
		oldV := ve.Attrs[k]
		if v != oldV {
			patch.AddAttrs[k] = v
		}
	}
	// TODO: events
	// TODO: kids
	return patch
}

func (vt VText) Diff(other VNode) Patch {
	otherText, ok := other.(VText)
	if ok && otherText == vt {
		return NopPatch{}
	}
	return ReplacePatch{Replacement: other}
}
