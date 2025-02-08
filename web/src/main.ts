import './ui.ts'

import { CpuAdapter } from './cpu.ts'
import { loadWasm } from './wasm.ts'

(async () => {
    await loadWasm();
    console.log('WASM loaded')

    const canvas = document.getElementById('canvas') as HTMLCanvasElement
    const ctx = canvas.getContext('2d')!

    const cpu = new CpuAdapter()
    cpu.setRom('IBM')

    cpu.onDraw((arr) => {
        console.log('arr onDraw')
        console.log(arr)
        const pixels = Array.from(arr)
        console.log('pixels')
        console.log(pixels)

        
        // ctx.clearRect(0, 0, canvas.width, canvas.height)
        // ctx.fillStyle = 'black'
        // ctx.fillRect(0, 0, canvas.width, canvas.height)
    })

    cpu.onSound((isOn) => {
        console.log(isOn)
    })

    // cpu.cycle()a

    document.addEventListener('keydown', (event) => {
        event.preventDefault();
        cpu.setKey(event.key, true)
    })

    document.addEventListener('keyup', (event) => {
        event.preventDefault();
        cpu.setKey(event.key, false)
    })


    document.getElementById('btn-play')?.addEventListener('click', () => {
        console.log('play')
    })

    document.getElementById('btn-pause')?.addEventListener('click', () => {
        console.log('pause')
    })

    document.getElementById('btn-reset')?.addEventListener('click', () => {
        console.log('reset')
    })  

})()

// const canvas = document.getElementById('canvas') as HTMLCanvasElement
// const ctx = canvas.getContext('2d')!

// const cpu = new CpuAdapter()
// cpu.setRom(ROMS[0])

// cpu.onDraw(() => {
//     ctx.clearRect(0, 0, canvas.width, canvas.height)
//     ctx.fillStyle = 'black'
//     ctx.fillRect(0, 0, canvas.width, canvas.height)
// })

// cpu.onSound((isOn) => {
//     console.log(isOn)
// })

// cpu.cycle()