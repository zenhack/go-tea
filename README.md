[![GoDoc][godoc-img]][godoc]

This repository contains a simple Go library for working with Virtual
DOMs, as popularized by [React.JS][1]. When compiled to WASM it supports
patching the DOM via the `syscall/js` package.

The `./tea` package contains a higher-level interface for building
model-view-update apps, inspired by [Elm][2].

This library is experimental, but I am currently using it in
[Tempest][3].

[1]: https://react.dev/
[2]: https://elm-lang.org/
[3]: https://github.com/zenhack/tempest

[godoc]: https://pkg.go.dev/zenhack.net/go/vdom
[godoc-img]: https://pkg.go.dev/badge/zenhack.net/go/vdom
