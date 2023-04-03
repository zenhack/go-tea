package tea

import (
	"context"
	"syscall/js"

	"zenhack.net/go/vdom"
)

type Message[Model any] interface {
	Update(Model) (Model, func(context.Context, func(Message[Model])))
}

type AppModel[Model any] interface {
	View(MessageSender[Model]) vdom.VNode
}

type App[M AppModel[M]] struct {
	model M
	msgs  chan Message[M]
}

type MessageSender[Model any] interface {
	Send(Message[Model])
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

func NewApp[M AppModel[M]](model M) *App[M] {
	return &App[M]{
		model: model,
		msgs:  make(chan Message[M]),
	}
}

func (app *App[Model]) SendMessage(msg Message[Model]) {
	app.msgs <- msg
}

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
