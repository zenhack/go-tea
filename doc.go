// Package vdom implements a Virtual DOM, as popularized by React.JS.
//
// Virtual DOM nodes can be constructed with VText and VElem, and patches
// can be generated with the VNode.Diff method. When run in the browser,
// patches may be applied to the real DOM to update with minimal changes.
//
// See the builder package for convenient helper functions, and the tea
// package for a higher-level framework for building apps.
package vdom
