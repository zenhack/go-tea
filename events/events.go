package events

import "zenhack.net/go/vdom"

func handler(f func(vdom.Event) any) vdom.EventHandler {
	return &f
}

func OnInput(f func(string)) vdom.EventHandler {
	return handler(func(e vdom.Event) any {
		f(e.Value.Get("target").Get("value").String())
		return nil
	})
}
