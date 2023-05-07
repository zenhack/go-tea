package vdom

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiff(t *testing.T) {
	t.Parallel()
	cases := []struct {
		Name          string
		Before, After VNode
		Patch         Patch
	}{
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
				Tag: "a",
				Children: []VNode{
					VText("kid1"),
					VText("kid2"),
				},
			},
			After: &VElem{
				Tag: "b",
				Children: []VNode{
					VText("kid1"),
					VText("kid2"),
				},
			},
			Patch: ReplacePatch{
				Replacement: &VElem{
					Tag: "b",
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
				Tag: "a",
				Children: []VNode{
					VText("kid1"),
					VText("kid2"),
				},
			},
			After: &VElem{
				Tag: "a",
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
			Name: "Drop excess trailing children",
			Before: &VElem{
				Tag: "a",
				Children: []VNode{
					VText("a"),
					VText("b"),
				},
			},
			After: &VElem{
				Tag: "a",
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
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			this := c
			t.Parallel()
			expected := this.Patch
			actual := this.Before.Diff(this.After)
			require.Equal(t, expected, actual)
		})
	}
}
