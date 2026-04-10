<template>
  <canvas
    ref="canvasRef"
    class="fixed inset-0 pointer-events-none"
    :style="{ zIndex: 0 }"
  />
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const canvasRef = ref<HTMLCanvasElement>()

const CELL_SIZE = 40 // 网格单元格大小（像素）
const GLOW_RADIUS = 2 // 发光半径（以单元格为单位）
const TRAIL_DURATION = 500 // 轨迹持续时间（毫秒）
const MAX_TRAIL_CELLS = 30 // 最多同时存在的轨迹单元格数量

let ctx: CanvasRenderingContext2D | null = null
let width = 0
let height = 0
let cols = 0
let rows = 0
let mouseX = -9999
let mouseY = -9999
let animFrame: number | null = null
let isDark = true

interface TrailCell {
  row: number
  col: number
  intensity: number
  lastSeen: number
}
const trail = new Map<string, TrailCell>()

function cellKey(row: number, col: number) {
  return `${row},${col}`
}

function resize() {
  const canvas = canvasRef.value
  if (!canvas) return
  const dpr = window.devicePixelRatio || 1
  width = window.innerWidth
  height = window.innerHeight
  canvas.width = width * dpr
  canvas.height = height * dpr
  canvas.style.width = `${width}px`
  canvas.style.height = `${height}px`
  ctx = canvas.getContext('2d')
  ctx?.scale(dpr, dpr)
  cols = Math.ceil(width / CELL_SIZE) + 1
  rows = Math.ceil(height / CELL_SIZE) + 1
}

function onMouseMove(e: MouseEvent) {
  mouseX = e.clientX
  mouseY = e.clientY
}

function onMouseLeave() {
  mouseX = -9999
  mouseY = -9999
}

function update(now: number) {
  // Find grid cell under cursor
  const curCol = Math.floor(mouseX / CELL_SIZE)
  const curRow = Math.floor(mouseY / CELL_SIZE)

  // Activate cells within glow radius
  const radius = Math.ceil(GLOW_RADIUS)
  for (let dr = -radius; dr <= radius; dr++) {
    for (let dc = -radius; dc <= radius; dc++) {
      const r = curRow + dr
      const c = curCol + dc
      const dist = Math.sqrt(dr * dr + dc * dc)
      if (dist > GLOW_RADIUS) continue

      const key = cellKey(r, c)
      const intensity = Math.max(0, 1 - dist / GLOW_RADIUS)
      const existing = trail.get(key)
      if (existing) {
        existing.intensity = Math.max(existing.intensity, intensity)
        existing.lastSeen = now
      } else {
        if (trail.size < MAX_TRAIL_CELLS) {
          trail.set(key, { row: r, col: c, intensity, lastSeen: now })
        }
      }
    }
  }

  // Fade out old trail cells
  for (const [key, cell] of trail) {
    const age = now - cell.lastSeen
    if (age > TRAIL_DURATION) {
      trail.delete(key)
    }
  }
}

function getThemeColors(): {
  baseLine: string
  accentR: number
  accentG: number
  accentB: number
} {
  const style = getComputedStyle(document.body)
  const lineColor = style.getPropertyValue('--grid-line-color').trim()
  const accent = style.getPropertyValue('--neon-cyan').trim()

  // Parse accent color
  let accentR = 6, accentG = 182, accentB = 212
  if (accent.startsWith('#') && accent.length === 7) {
    accentR = parseInt(accent.slice(1, 3), 16)
    accentG = parseInt(accent.slice(3, 5), 16)
    accentB = parseInt(accent.slice(5, 7), 16)
  }

  return {
    baseLine: lineColor || 'rgba(14, 165, 233, 0.06)',
    accentR,
    accentG,
    accentB
  }
}

function draw(now: number) {
  if (!ctx) return
  ctx.clearRect(0, 0, width, height)

  const theme = getThemeColors()

  // Draw base grid lines
  ctx.strokeStyle = theme.baseLine
  ctx.lineWidth = 1

  ctx.beginPath()
  for (let c = 0; c <= cols; c++) {
    const x = c * CELL_SIZE
    ctx.moveTo(x, 0)
    ctx.lineTo(x, height)
  }
  for (let r = 0; r <= rows; r++) {
    const y = r * CELL_SIZE
    ctx.moveTo(0, y)
    ctx.lineTo(width, y)
  }
  ctx.stroke()

  // Draw trail cells
  for (const cell of trail.values()) {
    const age = now - cell.lastSeen
    const fadeFactor = age < TRAIL_DURATION
      ? Math.max(0, 1 - (age / TRAIL_DURATION) * (age / TRAIL_DURATION))
      : 0
    const alpha = cell.intensity * fadeFactor
    if (alpha < 0.01) continue

    const x = cell.col * CELL_SIZE
    const y = cell.row * CELL_SIZE
    const { accentR, accentG, accentB } = theme

    // Cell fill
    ctx.fillStyle = `rgba(${accentR}, ${accentG}, ${accentB}, ${alpha * 0.08})`
    ctx.fillRect(x, y, CELL_SIZE, CELL_SIZE)

    // Cell border glow
    ctx.strokeStyle = `rgba(${accentR}, ${accentG}, ${accentB}, ${alpha * 0.5})`
    ctx.lineWidth = 1
    ctx.strokeRect(x + 0.5, y + 0.5, CELL_SIZE - 1, CELL_SIZE - 1)

    // Intersection dots
    if (alpha > 0.3) {
      const corners = [
        [x, y],
        [x + CELL_SIZE, y],
        [x, y + CELL_SIZE],
        [x + CELL_SIZE, y + CELL_SIZE]
      ]
      for (const [cx, cy] of corners) {
        ctx.beginPath()
        ctx.arc(cx, cy, 1.5 * alpha, 0, Math.PI * 2)
        ctx.fillStyle = `rgba(${accentR}, ${accentG}, ${accentB}, ${alpha * 0.8})`
        ctx.fill()
      }
    }
  }

  // Draw bright glow around cursor position
  if (mouseX > -1000 && mouseY > -1000) {
    const gradient = ctx.createRadialGradient(
      mouseX, mouseY, 0,
      mouseX, mouseY, CELL_SIZE * GLOW_RADIUS
    )
    gradient.addColorStop(0, `rgba(${theme.accentR}, ${theme.accentG}, ${theme.accentB}, 0.12)`)
    gradient.addColorStop(1, `rgba(${theme.accentR}, ${theme.accentG}, ${theme.accentB}, 0)`)
    ctx.fillStyle = gradient
    ctx.fillRect(
      mouseX - CELL_SIZE * GLOW_RADIUS,
      mouseY - CELL_SIZE * GLOW_RADIUS,
      CELL_SIZE * GLOW_RADIUS * 2,
      CELL_SIZE * GLOW_RADIUS * 2
    )
  }
}

function loop() {
  const now = performance.now()
  update(now)
  draw(now)
  animFrame = requestAnimationFrame(loop)
}

onMounted(() => {
  resize()
  isDark = document.body.classList.contains('dark')
  window.addEventListener('resize', resize)
  window.addEventListener('mousemove', onMouseMove)
  document.addEventListener('mouseleave', onMouseLeave)
  loop()
})

onUnmounted(() => {
  if (animFrame) cancelAnimationFrame(animFrame)
  window.removeEventListener('resize', resize)
  window.removeEventListener('mousemove', onMouseMove)
  document.removeEventListener('mouseleave', onMouseLeave)
})
</script>
