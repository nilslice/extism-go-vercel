package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	extism "github.com/extism/go-sdk"
)

var text2img *extism.Plugin

func init() {
	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmUrl{
				Url: "https://cdn.modsurfer.dylibso.com/api/v1/module/2c9eb901052b1e6397d2414bdb796975407cc87085e6b5fe9564932538d8af51.wasm",
			},
		},
	}
	var err error
	ctx := context.Background()
	text2img, err = extism.NewPlugin(ctx, manifest, extism.PluginConfig{}, nil)
	if err != nil {
		panic(err)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	input := make(map[string]any)
	input["value"] = params.Get("text")
	input["color"] = params.Get("color")
	fontSize, err := strconv.Atoi(params.Get("font_size"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %v", err)
	}
	input["font_size"] = fontSize

	data, err := json.Marshal(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %v", err)
	}

	code, output, err := text2img.Call("handle", data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %v (code = %d)", err, code)
		return
	}

	out := make(map[string]string)
	err = json.Unmarshal(output, &out)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %v", err)
		return
	}

	img, err := base64.StdEncoding.DecodeString(out["value"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %v", err)
		return
	}

	w.WriteHeader(200)
	w.Header().Add("content-type", "image/png")
	w.Write(img)
}
