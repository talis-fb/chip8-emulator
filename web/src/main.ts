import { CpuAdapter } from './cpu.ts'
import { loadWasm } from './wasm.ts'

(async () => {
    await loadWasm();
    console.log('WASM loaded')

    const canvas = document.getElementById('canvas') as HTMLCanvasElement
    const ctx = canvas.getContext('2d')!

    const cpu = new CpuAdapter()
    cpu.setRom('IBM')

    ctx.fillStyle = 'black'
    ctx.fillRect(1, 1, 1, 1)

    cpu.onDraw((arr) => {
        console.log('arr onDraw')
        console.log(arr)
        const pixels = Array.from(arr)
        console.log('pixels to paint')
        console.log(pixels.map((p, i) => [i, p] as const).filter(([_, p]) => p !== 0).map(([i]) => i))
        
        pixels.map((p) => Boolean(p)).forEach((p, i) => {
            const x = i % 64
            const y = Math.floor(i / 64)

            ctx.fillRect(x, y, 1, 1)
            if (p) {
                ctx.fillStyle = 'black'
            } else {
                ctx.fillStyle = 'white'
            }

        })
        
        // ctx.clearRect(0, 0, canvas.width, canvas.height)
        // ctx.fillStyle = 'black'
        // ctx.fillRect(0, 0, canvas.width, canvas.height)
    })

    cpu.onSound((isOn) => {
        console.log(isOn)
    })


    let isRunning = false

    let animationID: number
    const playRunner = () => {
        animationID = setInterval(() => {
            if (!isRunning) {
                window.clearInterval(animationID)
                return;
            }
            
            cpu.cycle()
        }, 400)
    }

    // const playRunner = () => {
    //     if (!isRunning) {
    //         window.cancelAnimationFrame(animationID)
    //         return;
    //     }
        
    //     cpu.cycle()

    //     animationID = requestAnimationFrame(playRunner)
    // }

    document.addEventListener('keydown', (event) => {
        // event.preventDefault();
        cpu.setKey(event.key, true)
    })

    document.addEventListener('keyup', (event) => {
        // event.preventDefault();
        cpu.setKey(event.key, false)
    })


    document.getElementById('btn-play')?.addEventListener('click', () => {
        console.log('play')
        if (!isRunning) {
            isRunning = true
            requestAnimationFrame(playRunner)
        }

    })

    document.getElementById('btn-pause')?.addEventListener('click', () => {
        console.log('pause')
        isRunning = false
    })

    document.getElementById('btn-reset')?.addEventListener('click', () => {
        console.log('reset')
        isRunning = false
    })  

})()
