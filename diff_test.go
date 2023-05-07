package vdom

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type DiffCase struct {
	Name          string
	Before, After VNode
	Patch         Patch
}

var diffCases = []DiffCase{
	{
		Name:   "Same vtext noop",
		Before: VText("some text"),
		After:  VText("some text"),
		Patch:  NopPatch{},
	},
	{
		Name:   "Different vtext replace",
		Before: VText("some text"),
		After:  VText("some other text"),
		Patch: ReplacePatch{
			Replacement: VText("some other text"),
		},
	},
	{
		Name: "Different tags replace",
		Before: &VElem{
			Tag:   "a",
			Attrs: map[string]string{},
			Children: []VNode{
				VText("kid1"),
				VText("kid2"),
			},
		},
		After: &VElem{
			Tag:   "b",
			Attrs: map[string]string{},
			Children: []VNode{
				VText("kid1"),
				VText("kid2"),
			},
		},
		Patch: ReplacePatch{
			Replacement: &VElem{
				Tag:   "b",
				Attrs: map[string]string{},
				Children: []VNode{
					VText("kid1"),
					VText("kid2"),
				},
			},
		},
	},
	{
		Name: "Same tags recurse",
		Before: &VElem{
			Tag:   "a",
			Attrs: map[string]string{},
			Children: []VNode{
				VText("kid1"),
				VText("kid2"),
			},
		},
		After: &VElem{
			Tag:   "a",
			Attrs: map[string]string{},
			Children: []VNode{
				VText("kid1"),
				VText("kid3"),
			},
		},
		Patch: ModifyPatch{
			Attrs:  AttrsPatch{Add: map[string]string{}},
			Events: EventsPatch{Add: map[string]EventHandler{}},
			Children: ChildPatch{
				Common: []Patch{
					NopPatch{},
					ReplacePatch{
						Replacement: VText("kid3"),
					},
				},
				Append: []VNode{},
			},
		},
	},
	{
		Name: "Drop excess trailing children (1)",
		Before: &VElem{
			Tag:   "a",
			Attrs: map[string]string{},
			Children: []VNode{
				VText("a"),
				VText("b"),
			},
		},
		After: &VElem{
			Tag:   "a",
			Attrs: map[string]string{},
			Children: []VNode{
				VText("a"),
			},
		},
		Patch: ModifyPatch{
			Attrs:  AttrsPatch{Add: map[string]string{}},
			Events: EventsPatch{Add: map[string]EventHandler{}},
			Children: ChildPatch{
				Common: []Patch{
					NopPatch{},
				},
				Append: []VNode{},
				Drop:   1,
			},
		},
	},
	{
		Name: "Drop excess trailing children (3)",
		Before: &VElem{
			Tag:   "a",
			Attrs: map[string]string{},
			Children: []VNode{
				VText("a"),
				VText("b"),
				VText("c"),
				VText("d"),
			},
		},
		After: &VElem{
			Tag:   "a",
			Attrs: map[string]string{},
			Children: []VNode{
				VText("a"),
			},
		},
		Patch: ModifyPatch{
			Attrs:  AttrsPatch{Add: map[string]string{}},
			Events: EventsPatch{Add: map[string]EventHandler{}},
			Children: ChildPatch{
				Common: []Patch{
					NopPatch{},
				},
				Append: []VNode{},
				Drop:   3,
			},
		},
	},
}

func TestDiff(t *testing.T) {
	t.Parallel()
	for _, c := range diffCases {
		t.Run(c.Name, func(t *testing.T) {
			this := c
			t.Parallel()
			expected := this.Patch
			actual := this.Before.Diff(this.After)
			require.Equal(t, expected, actual)
		})
	}
}
