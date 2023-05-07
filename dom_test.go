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
	cases := []struct {
		Name string
		Node VNode
	}{
		{
			Name: "Text node",
			Node: VText("some text"),
		},
		{
			Name: "Simple element",
			Node: &VElem{
				Tag: "a",
				Attrs: map[string]string{
					"href": "#",
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			// Make sure rendering and then parsing the node gives the
			// same result:
			this := c
			t.Parallel()
			dom := this.Node.ToDomNode()
			parsed := vdomFromJS(dom.Value)
			require.Equal(t, this.Node, parsed)
		})
	}
}
