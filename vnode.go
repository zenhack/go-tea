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

// N.B. the pointer indirection is so these can be compared.
type EventHandler = *func(Event) any

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

	return ModifyPatch{
		Attrs:    diffAttrs(ve.Attrs, dstElem.Attrs),
		Events:   diffEvents(ve.Events, dstElem.Events),
		Children: diffChildren(ve.Children, dstElem.Children),
	}
}

func diffAttrs(src, dst map[string]string) AttrsPatch {
	var patch AttrsPatch
	for k, _ := range src {
		if _, ok := dst[k]; !ok {
			patch.Remove = append(patch.Remove, k)
		}
	}
	patch.Add = make(map[string]string)
	for k, v := range dst {
		oldV := src[k]
		if v != oldV {
			patch.Add[k] = v
		}
	}
	return patch
}

func diffEvents(src, dst map[string]EventHandler) EventsPatch {
	var patch EventsPatch
	for k := range src {
		if _, ok := dst[k]; !ok {
			patch.Remove = append(patch.Remove, k)
		}
	}

	patch.Add = make(map[string]EventHandler)
	for k := range dst {
		if src[k] == dst[k] {
			continue
		}
		if _, ok := src[k]; ok {
			patch.Remove = append(patch.Remove, k)
		}
		patch.Add[k] = dst[k]
	}

	return patch
}

func diffChildren(src, dst []VNode) ChildPatch {
	var patch ChildPatch

	minLen := len(src)
	if len(dst) < minLen {
		minLen = len(dst)
	}

	for i := 0; i < minLen; i++ {
		patch.Common = append(patch.Common, src[i].Diff(dst[i]))
	}

	patch.Append = dst[minLen:]
	if len(src) > minLen {
		patch.Drop = len(src) - minLen
	}

	return patch
}

func (vt VText) Diff(dst VNode) Patch {
	dstText, ok := dst.(VText)
	if ok && dstText == vt {
		return NopPatch{}
	}
	return ReplacePatch{Replacement: dst}
}
