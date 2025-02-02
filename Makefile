setup:
	docker run -it golang:1.23-alpine cat /usr/local/go/misc/wasm/wasm_exec.js > frontend/wasm/wasm_exec.js
