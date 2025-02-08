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
        // ctx.clearRect(0, 0, canvas.width, canvas.height)
        // ctx.fillStyle = 'black'
        // ctx.fillRect(0, 0, canvas.width, canvas.height)
    })

    cpu.onSound((isOn) => {
        console.log(isOn)
    })

    cpu.cycle()

    document.addEventListener('keydown', (event) => {
        cpu.setKey(event.keyCode, true)
    })

    document.addEventListener('keyup', (event) => {
        cpu.setKey(event.keyCode, false)
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