package main

import (
	"context"
	"strconv"
	"syscall/js"

	"zenhack.net/go/tea"
	"zenhack.net/go/tea/vdom"
	"zenhack.net/go/tea/vdom/builder"
)

type (
	a = builder.A
	e = builder.E
)

var (
	h = builder.H
	t = builder.T
)

type Model int

type Cmd = func(context.Context, func(tea.Message[Model]))

func (m Model) View(ms tea.MessageSender[Model]) vdom.VNode {
	return h("div", nil, nil,
		h("button", nil, e{"click": ms.Event(Increment{})}, t("+")),
		t(strconv.Itoa(int(m))),
		h("button", nil, e{"click": ms.Event(Decrement{})}, t("-")),
	)
}

type Increment struct{}
type Decrement struct{}

func (msg Increment) Update(m *Model) Cmd {
	*m = *m + 1
	return nil
}
func (msg Decrement) Update(m *Model) Cmd {
	*m = *m - 1
	return nil
}

func main() {
	app := tea.NewApp[Model](0)
	elem := js.Global().Get("document").Call("getElementById", "app")
	app.Run(context.Background(), vdom.DomNode{Value: elem})
}
