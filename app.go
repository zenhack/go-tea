// Package tea implements a high-level framework for building web user
// interfaces, using the Elm architecture.
//
// The user interface consists of an App, which manages a model representing
// application state. The app responds to messages by updating its model, and
// re-renders via a virtual dom when the model changes. Messages are processed
// sequentially in a single goroutine.
//
// Applications should define a model type M which implements AppModel[M],
// and some number of types that implement Message[M]. The app can then
// be constructed by passing the initial model to NewApp, and then attached
// to the DOM and executed with App.Run().
package tea

import (
	"context"
	"syscall/js"

	"zenhack.net/go/tea/vdom"
)

// A Message[Model] is a message carrying updates to make to a Model.
type Message[Model any] interface {
	// Return an updated model, and a function to be run asynchronously
	// with the ability to send more messages.
	//
	// Update MUST NOT block, as doing so could hang the user interface.
	// Anything that could block must instead be done in the returned
	// callback, whose second parameter can be used to send messages
	// for further updates.
	Update(Model) (Model, func(context.Context, func(Message[Model])))
}

// An AppModel holds application state, and knows how to render itself as
// a virtual DOM. The Model type parameter will typically be the same as the
// type of the receiver.
type AppModel[Model any] interface {
	// View renders a virtual DOM from the model. The MessageSender
	// can be used to implement event handlers which send update messages
	// to the app.
	View(MessageSender[Model]) vdom.VNode
}

// An App[M] manages an application with model/state type M.
type App[M AppModel[M]] struct {
	model M               // The current application state
	msgs  chan Message[M] // Incoming messages
}

// A MessageSender[M] sends messages to an application with model type M.
type MessageSender[Model any] interface {
	// Send sends a message to the application.
	Send(Message[Model])

	// Event returns a vdom event handler which sends the specified
	// message when triggered.
	Event(Message[Model]) vdom.EventHandler
}

type messageSender[M AppModel[M]] struct {
	app *App[M]
}

func (ms messageSender[Model]) Send(msg Message[Model]) {
	ms.app.SendMessage(msg)
}

func (ms messageSender[Model]) Event(msg Message[Model]) vdom.EventHandler {
	ret := func(vdom.Event) any {
		ms.app.SendMessage(msg)
		return nil
	}
	return &ret
}

// NewApp creates a new application with model/state type M.
func NewApp[M AppModel[M]](model M) *App[M] {
	return &App[M]{
		model: model,
		msgs:  make(chan Message[M]),
	}
}

// SendMessage sends a message to the application
func (app *App[Model]) SendMessage(msg Message[Model]) {
	app.msgs <- msg
}

// Run replaces node in the DOM with a node managed by the application,
// and then runs the application.
func (app *App[Model]) Run(ctx context.Context, node vdom.DomNode) {
	model := app.model
	parent := vdom.DomNode{Value: node.Value.Get("parentNode")}
	var animationFrame struct {
		ch        chan struct{}
		requested bool
	}
	animationFrame.ch = make(chan struct{}, 1)

	onRequestAnimationFrame := js.FuncOf(func(this js.Value, args []js.Value) any {
		animationFrame.ch <- struct{}{}
		return nil
	})
	defer onRequestAnimationFrame.Release()

	ms := messageSender[Model]{app: app}

	vnode := model.View(ms)
	oldVNode := vnode
	node = vdom.ReplacePatch{Replacement: vnode}.Patch(parent, node)

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-app.msgs:
			var cmd func(context.Context, func(Message[Model]))
			model, cmd = msg.Update(model)
			if cmd != nil {
				go cmd(ctx, app.SendMessage)
			}

			if !animationFrame.requested {
				js.Global().Get("window").Call(
					"requestAnimationFrame",
					onRequestAnimationFrame,
				)
				animationFrame.requested = true
			}
		case <-animationFrame.ch:
			animationFrame.requested = false
			vnode = model.View(ms)
			patch := oldVNode.Diff(vnode)
			node = patch.Patch(parent, node)
			oldVNode = vnode
		}
	}
}
