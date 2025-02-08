export type CpuExports = {
    setRom: (rom: Uint8Array) => void;
    onDraw: (cb: (arr: Uint8Array) => void) => void;
    onSound: (cb: (isOn: boolean) => void) => void;
    setKey: (key: number, isPressed: boolean) => void;
    reset: () => void;
    cycle: () => void;
}

// It assumes that loadWasm() has been called
//  and all WASM globals exported functions are available
const wasmExports = (window as unknown as Window & CpuExports)

export class CpuAdapter {
    public async setRom(romName: string) {
        try {
            const response = await fetch(`/roms/${romName}`);
            if (!response.ok) {
                throw new Error(`Failed to load ROM: ${romName}`);
            }

            const arrayBuffer = await response.arrayBuffer();
            const romData = new Uint8Array(arrayBuffer);

            //wasmExports.setRom(romData);

            console.log(`ROM "${romName}" loaded successfully.`);
        } catch (error) {
            console.error(`Error loading ROM "${romName}":`, error);
        }
    }

    public onDraw(cb: (arr: Uint8Array) => void) {
        wasmExports.onDraw(cb);
    }

    public onSound(cb: (isOn: boolean) => void) {
        // wasmExports.onSound(cb);
    }

    public setKey(key: number, isPressed: boolean) {
        //
    }

    public reset() {
        //
    }

    public cycle() {
        //
    }
}
