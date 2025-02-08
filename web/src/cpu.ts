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

            wasmExports.setRom(romData);

            console.log(`ROM "${romName}" loaded successfully.`);
        } catch (error) {
            console.error(`Error loading ROM "${romName}":`, error);
        }
    }

    public onDraw(cb: (arr: Uint8Array) => void) {
        wasmExports.onDraw(cb);
    }

    public onSound(cb: (isOn: boolean) => void) {
        wasmExports.onSound(cb);
    }

    public setKey(key: string, isPressed: boolean) {
        const keyCodes = {
            '1': 0x1,
            '2': 0x2,
            '3': 0x3,
            '4': 0xC,
            'q': 0x4,
            'w': 0x5,
            'e': 0x6,
            'r': 0xD,
            'a': 0x7,
            's': 0x8,
            'd': 0x9,
            'f': 0xE,
            'z': 0xA,
            'x': 0x0,
            'c': 0xB,
            'v': 0xF
        } as Record<string, number | undefined>;
 
        const cpuKeyCode = keyCodes[key];
        if (cpuKeyCode) {
            wasmExports.setKey(cpuKeyCode, isPressed)
        }
    }

    public reset() {
        wasmExports.reset()
    }

    public cycle() {
        // let currentSecond = new Date().getSeconds()
        // console.log(`cycle ${currentSecond}`)
        wasmExports.cycle()
    }
}
