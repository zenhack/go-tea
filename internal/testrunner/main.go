package main

import (
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var index = []byte(`<!doctype html>
<html>
  <head>
    <title>vdom tests</title>
    <script src="wasm_exec.js"></script>
    <script>
      const go = new Go();
      WebAssembly.
        instantiateStreaming(fetch("vdom.test"), go.importObject).
        then((result) => {
          go.run(result.instance);
        });
    </script>
  </head>
  <body>
  </body>
</html>`)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write(index)
	})
	http.HandleFunc("/vdom.test", func(w http.ResponseWriter, req *http.Request) {
		data, err := os.ReadFile("vdom.test")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(data)
	})
	http.HandleFunc("/wasm_exec.js", func(w http.ResponseWriter, req *http.Request) {
		out, err := exec.Command("go", "env", "GOROOT").Output()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		data, err := os.ReadFile(strings.TrimSpace(string(out)) + "/misc/wasm/wasm_exec.js")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(data)
	})

	const addr = ":4444"
	println("Open http://localhost" + addr + " in the browser to run the tests.")
	panic(http.ListenAndServe(addr, nil))
}
