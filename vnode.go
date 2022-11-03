package vdom

type VNode interface {
	Diff(dst VNode) Patch
	ToDomNode() DomNode
}

type VElem struct {
	Tag      string
	Attrs    map[string]string
	Events   map[string]EventHandler
	Children []VNode
}

type VText string

type EventHandler = func(Event) any

func (ve *VElem) Diff(dst VNode) Patch {
	dstElem, ok := dst.(*VElem)
	if !ok {
		return ReplacePatch{Replacement: dst}
	}
	if ve == dstElem {
		return NopPatch{}
	}
	if ve.Tag != dstElem.Tag {
		return ReplacePatch{Replacement: dst}
	}

	patch := ModifyPatch{}
	for k, _ := range ve.Attrs {
		if _, ok := dstElem.Attrs[k]; !ok {
			patch.RemoveAttrs = append(patch.RemoveAttrs, k)
		}
	}
	patch.AddAttrs = make(map[string]string)
	for k, v := range dstElem.Attrs {
		oldV := ve.Attrs[k]
		if v != oldV {
			patch.AddAttrs[k] = v
		}
	}
	// TODO: events
	// TODO: kids
	return patch
}

func (vt VText) Diff(dst VNode) Patch {
	dstText, ok := dst.(VText)
	if ok && dstText == vt {
		return NopPatch{}
	}
	return ReplacePatch{Replacement: dst}
}
