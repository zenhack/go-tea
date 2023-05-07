//go:build js

package vdom

import (
	"strings"
	"syscall/js"
	"testing"

	"github.com/stretchr/testify/require"
)

func vdomFromJS(v js.Value) VNode {
	if v.Get("tagName").IsUndefined() {
		return VText(v.Get("textContent").String())
	} else {
		elem := &VElem{
			Tag:   strings.ToLower(v.Get("tagName").String()),
			Attrs: map[string]string{},
		}

		attrs := v.Get("attributes")
		for i := 0; i < attrs.Length(); i++ {
			attr := attrs.Index(0)
			elem.Attrs[attr.Get("name").String()] = attr.Get("value").String()
		}

		// TODO: events?

		childNodes := v.Get("childNodes")
		for i := 0; i < childNodes.Length(); i++ {
			elem.Children = append(elem.Children, vdomFromJS(childNodes.Index(i)))
		}

		return elem
	}
}

// Self tests for the vdomFromJS function.
func TestVdomFromJS(t *testing.T) {
	t.Parallel()
	for _, c := range diffCases {
		t.Run(c.Name, func(t *testing.T) {
			this := c
			t.Parallel()
			t.Run("Before", func(t *testing.T) {
				testVdomFromJS(t, this.Before)
			})
			t.Run("After", func(t *testing.T) {
				testVdomFromJS(t, this.After)
			})
		})
	}
}

func testVdomFromJS(t *testing.T, orig VNode) {
	t.Parallel()
	// Make sure rendering and then parsing the node gives the
	// same result:
	dom := orig.ToDomNode()
	parsed := vdomFromJS(dom.Value)
	require.Equal(t, orig, parsed)
}

func TestPatchAppliesCorrectly(t *testing.T) {
	t.Parallel()
	for _, c := range diffCases {
		t.Run(c.Name, func(t *testing.T) {
			this := c
			t.Parallel()
			dom := this.Before.ToDomNode()
			parentDiv := js.Global().Get("document").Call("createElement", "div")
			parentDiv.Call("appendChild", dom.Value)
			newDom := this.Patch.Patch(DomNode{Value: parentDiv}, dom)
			actualVdom := vdomFromJS(newDom.Value)
			require.Equal(t, this.After, actualVdom)
		})
	}
}
