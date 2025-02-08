export let isWasmReady = false;

export async function loadWasm() {
    // Imported at index.html
    const go = new (window as any).Go();
    const wasm = await WebAssembly.instantiateStreaming(fetch("/wasm/cpu.wasm"), go.importObject);
    go.run(wasm.instance);

    (window as any).helloWasm();

    isWasmReady = true;
}