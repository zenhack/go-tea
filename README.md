[![GoDoc][godoc-img]][godoc]

This contains a Go library for writing client-side web interfaces, in
the model-view-update style of the [Elm][elm] architecture ("TEA").

It contains a simple package `vdom` for working with Virtual DOMs, as
popularized by [React.JS][react]. When compiled to WASM, package `vdom`
supports patching the DOM via the `syscall/js` package. `vdom` can be
used stand-alone.

This library is experimental, but I am currently using it in
[Tempest][tempest].

There are some example applications in the `examples/` directory.

[godoc]: https://pkg.go.dev/zenhack.net/go/tea
[godoc-img]: https://pkg.go.dev/badge/zenhack.net/go/tea
[elm]: https://elm-lang.org/
[react]: https://react.dev/
[tempest]: https://github.com/zenhack/tempest
