<template>
  <div class="editor-page">

    <!-- Sidebar -->
    <div class="sidebar" :class="{ open: isSidebarOpen }">
      <div class="sidebar-header">
        <div class="sidebar-title-row">

          <button class="sidebar-close-btn" @click="isSidebarOpen = false">
            <h2>X</h2>
            <X :size="14" />
          </button>

        </div>
      </div>



      <!-- Add Table -->
      <div class="panel">
        <div class="panel-title">เพิ่มโต๊ะ</div>
        <div class="form-group">
          <label>ที่นั่ง (ค่าเริ่มต้น)</label>
          <input v-model.number="defaultTableSeats" type="number" min="1" step="1" />
        </div>
        <div class="form-group">
          <label>ราคา (บาท)</label>
          <input v-model.number="newTable.price" type="number" min="0" step="1" placeholder="0" />
        </div>
        <div class="form-group">
          <label>ที่นั่ง</label>
          <input v-model.number="newTable.seats" type="number" min="1" step="1" />
        </div>
        <button @click="addTable" class="btn-add-table" :disabled="addingTable">
          {{ addingTable ? '...' : '+ เพิ่มโต๊ะ' }}
        </button>

      </div>

      <!-- Add Floor Items -->
      <div class="panel">
        <div class="panel-title">
          <Plus :size="12" /> เพิ่ม Element
        </div>
        <div class="element-grid">
          <button v-for="el in elementTypes" :key="el.type" @click="addFloorItem(el)" class="btn-element">

            <span class="el-label">{{ el.label }}</span>

          </button>
        </div>
      </div>

      <!-- Selected Item Settings -->
      <div class="panel" v-if="selected">
        <div class="panel-title">
          <Settings :size="12" /> ตั้งค่า
        </div>

        <template v-if="selected.itemType === 'table'">
          <div class="form-group">
            <label>เลขโต๊ะ</label>
            <input v-model.number="selected.number" disabled style="background: #f5f5f5;" />
          </div>
          <div class="form-group">
            <label>ที่นั่ง</label>
            <input v-model.number="selected.seats" type="number" min="1" step="1" @input="saveTableEdit" />
          </div>
          <div class="form-group">
            <label>ราคา (บาท)</label>
            <input v-model.number="selected.price" type="number" min="0" step="1" @input="saveTableEdit" />
          </div>
          <div class="form-group toggle-row">
            <label>เปิดให้บริการ</label>
            <input type="checkbox" v-model="selected.is_active" @change="saveTableEdit" class="toggle" />
          </div>
        </template>

        <template v-if="selected.itemType === 'floor'">
          <div class="form-group">
            <label>ชื่อที่แสดง</label>
            <input v-model="selected.label" @input="saveFloorEdit" />
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>กว้าง (px)</label>
              <input v-model.number="selected.width" type="number" min="40" @input="saveFloorEdit" />
            </div>
            <div class="form-group">
              <label>สูง (px)</label>
              <input v-model.number="selected.height" type="number" min="30" @input="saveFloorEdit" />
            </div>
          </div>
        </template>

        <button @click="deleteSelected" class="btn-delete">
          <Trash2 :size="13" /> ลบ
        </button>
      </div>


    </div>

    <!-- Canvas -->
    <div class="canvas-wrapper">
      <div class="canvas-toolbar">
        <button v-if="!isSidebarOpen" class="open-sidebar-btn" @click="isSidebarOpen = true">
          <Menu :size="15" /> เมนู
        </button>
        <h3>ผังร้าน</h3>
        <span class="canvas-hint">
          {{ tables.length }} โต๊ะ
        </span>
        <router-link to="/dashboard" class="btn-back">
          <ArrowLeft :size="14" /> กลับหน้าหลัก
        </router-link>
      </div>
      <div class="canvas-scroll">
        <!-- BUG 1 แก้: ย้าย event handlers มาที่ div ใน canvas-scroll แทน canvas-scroll -->
        <div class="canvas" ref="canvas" :style="canvasStyle" @mousedown="handleCanvasMouseDown"
          @mousemove="onMouseMove" @mouseup="onMouseUp" @touchmove.prevent="onTouchMove" @touchend="onTouchEnd"
          @touchstart="handleCanvasTouchStart" @wheel.prevent="onWheelZoom">
          <!-- Grid lines overlay -->
          <svg v-if="showGridLines" class="grid-svg" :width="canvasWidth" :height="canvasHeight">
            <line v-for="x in gridLinesCols" :key="'v' + x" :x1="x" y1="0" :x2="x" :y2="canvasHeight.value"
              stroke="#c8ccdd" stroke-width="0.5" />
            <line v-for="y in gridLinesRows" :key="'h' + y" x1="0" :y1="y" :x2="canvasWidth.value" :y2="y"
              stroke="#c8ccdd" stroke-width="0.5" />
          </svg>

          <!-- Ghost cell highlight ขณะ drag -->
          <div v-if="ghostCell" class="ghost-cell" :style="{
            left: ghostCell.x + 'px', top: ghostCell.y + 'px',
            width: ghostCell.w + 'px', height: ghostCell.h + 'px',
          }" />

          <!-- Floor Items -->
          <!-- BUG 2 แก้: ย้าย floor-item ออกมาจากข้างใน ghost-cell div -->
          <div v-for="item in floorItems" :key="'fi-' + item.id" class="floor-item"
            :class="[`fi-${item.type}`, { selected: selected?.id === item.id && selected?.itemType === 'floor' }]"
            :style="{ left: item.pos_x + 'px', top: item.pos_y + 'px', width: item.width + 'px', height: item.height + 'px' }"
            @mousedown.stop="startDragFloor($event, item)" @touchstart.stop.prevent="startDragFloorTouch($event, item)"
            @click.stop="selectFloor(item)">
            <span class="fi-label">{{ item.label }}</span>
          </div>

          <!-- Tables -->
          <div v-for="table in tables" :key="'t-' + table.id" class="table-node" :class="{
            selected: selected?.id === table.id && selected?.itemType === 'table',
            inactive: !table.is_active,
            dragging: isDragging && dragging?.item?.id === table.id
          }"
            :style="{ left: table.pos_x + 'px', top: table.pos_y + 'px', width: tableSize + 'px', height: tableSize + 'px' }"
            @mousedown.stop="startDragTable($event, table)"
            @touchstart.stop.prevent="startDragTableTouch($event, table)" @click.stop="selectTable(table)">
            <div class="t-name">{{ table.number }}</div>
            <div v-if="table.price > 0" class="t-price">
              ฿{{ table.price.toLocaleString('th-TH', { minimumFractionDigits: 0, maximumFractionDigits: 0 }) }}
            </div>
            <div v-if="!table.is_active" class="t-closed">ปิด</div>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { tableApi } from '../api/auth'
import { X, Plus, Settings, Trash2, Menu, ArrowLeft } from 'lucide-vue-next'

// ---- Constants ----
const CANVAS_PADDING = 50
//Sidebar ตรวจขนาดจอ
const isSidebarOpen = ref(true)
const isMobile = ref(false)

function checkScreen() {
  isMobile.value = window.innerWidth <= 768
  if (isMobile.value) {
    isSidebarOpen.value = false
  } else {
    isSidebarOpen.value = true
  }
}

onMounted(() => {
  checkScreen()
  window.addEventListener('resize', checkScreen)

  // reset default seats to 4 (once)
  if (!localStorage.getItem('default_table_seats_reset_4_v1')) {
    localStorage.setItem('default_table_seats', '4')
    localStorage.setItem('default_table_seats_reset_4_v1', '1')
  }

  defaultTableSeats.value = Number(localStorage.getItem('default_table_seats') || 4)
  if (!newTable.value.seats || newTable.value.seats <= 0) {
    newTable.value.seats = defaultTableSeats.value
  }
})
// ---- State ----
const tables = ref([])
const floorItems = ref([])
const selected = ref(null)
const addingTable = ref(false)
const defaultTableSeats = ref(Number(localStorage.getItem('default_table_seats') || 4))
const newTable = ref({ price: 0, seats: defaultTableSeats.value })

const canvas = ref(null)
const isDragging = ref(false)   // BUG 3 แก้: ลบ let isDragging ออก เหลือแค่ ref นี้
const ghostCell = ref(null)
const gridSize = ref(24)
const showGridLines = ref(true)

let saveTimer = null
let dragging = null
let dragOffsetX = 0
let dragOffsetY = 0

const panX = ref(0)
const panY = ref(0)
const zoom = ref(1)
const MIN_ZOOM = 0.4
const MAX_ZOOM = 2.5

const isPanning = ref(false)
let panStartClientX = 0
let panStartClientY = 0
let panOriginX = 0
let panOriginY = 0

let panTouchStartX = 0
let panTouchStartY = 0
let panTouchOriginX = 0
let panTouchOriginY = 0
let lastTouchDistance = null

function onPanMouseMove(e) {
  if (!isPanning.value) return
  const dx = e.clientX - panStartClientX
  const dy = e.clientY - panStartClientY
  panX.value = panOriginX + dx
  panY.value = panOriginY + dy
}

function onPanMouseUp() {
  if (!isPanning.value) return
  isPanning.value = false
  window.removeEventListener('mousemove', onPanMouseMove)
  window.removeEventListener('mouseup', onPanMouseUp)
}

function handleCanvasMouseDown(e) {
  if (dragging) return
  if (e.button !== 0) return
  isPanning.value = true
  panStartClientX = e.clientX
  panStartClientY = e.clientY
  panOriginX = panX.value
  panOriginY = panY.value
  window.addEventListener('mousemove', onPanMouseMove)
  window.addEventListener('mouseup', onPanMouseUp)
}

function handleCanvasTouchStart(e) {
  if (dragging) return
  if (e.touches.length !== 1) return
  const touch = e.touches[0]
  isPanning.value = true
  panTouchStartX = touch.clientX
  panTouchStartY = touch.clientY
  panTouchOriginX = panX.value
  panTouchOriginY = panY.value
}

function onWheelZoom(e) {
  if (!canvas.value) return
  const rect = canvas.value.getBoundingClientRect()
  const cx = e.clientX - rect.left
  const cy = e.clientY - rect.top

  const worldX = (cx - panX.value) / zoom.value
  const worldY = (cy - panY.value) / zoom.value

  const direction = e.deltaY > 0 ? -1 : 1
  const factor = direction > 0 ? 1.08 : 0.92
  const newZoom = clamp(zoom.value * factor, MIN_ZOOM, MAX_ZOOM)

  zoom.value = newZoom
  panX.value = cx - worldX * newZoom
  panY.value = cy - worldY * newZoom
}

watch(defaultTableSeats, (v) => {
  const n = Number(v)
  if (!Number.isFinite(n) || n <= 0) return
  localStorage.setItem('default_table_seats', String(n))
  if (!newTable.value.seats || newTable.value.seats <= 0) {
    newTable.value.seats = n
  }
})

// ---- Computed ----
//สร้าง canvasWidth / canvasHeight แบบ dynamic
//สร้าง canvasWidth / canvasHeight แบบ dynamic
const MIN_CANVAS_W = 600
const MIN_CANVAS_H = 600


const canvasWidth = computed(() => {
  let maxRight = 0

  tables.value.forEach(t => {
    maxRight = Math.max(maxRight, t.pos_x + tableSize.value)
  })

  floorItems.value.forEach(f => {
    maxRight = Math.max(maxRight, f.pos_x + f.width)
  })

  return Math.max(MIN_CANVAS_W, maxRight + CANVAS_PADDING)
})

const canvasHeight = computed(() => {
  let maxBottom = 0

  tables.value.forEach(t => {
    maxBottom = Math.max(maxBottom, t.pos_y + tableSize.value)
  })

  floorItems.value.forEach(f => {
    maxBottom = Math.max(maxBottom, f.pos_y + f.height)
  })

  return Math.max(MIN_CANVAS_H, maxBottom + CANVAS_PADDING)
})
const tableSize = computed(() => gridSize.value * 2)

const canvasStyle = computed(() => ({
  backgroundSize: `${gridSize.value}px ${gridSize.value}px`,
  width: canvasWidth.value + 'px',
  height: canvasHeight.value + 'px',
  transform: `translate(${panX.value}px, ${panY.value}px) scale(${zoom.value})`,
  transformOrigin: '0 0',
}))

const gridLinesCols = computed(() => {
  const cols = []
  for (let x = gridSize.value; x < canvasWidth.value; x += gridSize.value) cols.push(x)
  return cols
})
const gridLinesRows = computed(() => {
  const rows = []
  for (let y = gridSize.value; y < canvasHeight.value; y += gridSize.value) rows.push(y)
  return rows
})

const elementTypes = [
  { type: 'stage', label: 'เวที', width: 192, height: 96 },
]

// ---- Snap / Clamp ----
function snap(val, size) { return Math.round(val / size) * size }
function clamp(val, min, max) { return Math.max(min, Math.min(max, val)) }

function onGridSizeChange() {
  tables.value.forEach(t => {
    t.pos_x = snap(t.pos_x, gridSize.value)
    t.pos_y = snap(t.pos_y, gridSize.value)
  })
  floorItems.value.forEach(fi => {
    fi.pos_x = snap(fi.pos_x, gridSize.value)
    fi.pos_y = snap(fi.pos_y, gridSize.value)
  })
}

onMounted(loadAll)
onUnmounted(() => {
  window.removeEventListener('mousemove', onMouseMove)
  window.removeEventListener('mouseup', onMouseUp)

  window.removeEventListener('mousemove', onPanMouseMove)
  window.removeEventListener('mouseup', onPanMouseUp)
})

async function loadAll() {
  const [tablesRes, floorRes] = await Promise.all([
    tableApi.adminGetTables(),
    tableApi.adminGetFloorItems(),
  ])
  tables.value = tablesRes.data.map(t => ({ ...t, itemType: 'table' }))
  floorItems.value = floorRes.data.map(f => ({ ...f, itemType: 'floor' }))
}

// ---- Add Table ----
async function addTable() {
  addingTable.value = true
  try {
    const cx = snap(canvasWidth.value / 2 - tableSize.value / 2, gridSize.value)
    const cy = snap(canvasHeight.value / 2 - tableSize.value / 2, gridSize.value)
    const res = await tableApi.adminCreateTable({
      seats: newTable.value.seats || defaultTableSeats.value || 4,
      price: newTable.value.price || 0,
      pos_x: cx, pos_y: cy,
    })
    tables.value.push({ ...res.data, itemType: 'table' })
    newTable.value = { price: 0, seats: defaultTableSeats.value }
  } finally {
    addingTable.value = false
  }
}

// ---- Add Floor Item ----
async function addFloorItem(el) {
  const w = snap(el.width, gridSize.value) || gridSize.value * 2
  const h = snap(el.height, gridSize.value) || gridSize.value * 2
  const cx = snap(canvasWidth.value / 2 - w / 2, gridSize.value)
  const cy = snap(canvasHeight.value / 2 - h / 2, gridSize.value)
  const res = await tableApi.adminCreateFloorItem({
    type: el.type, label: el.label,
    pos_x: cx, pos_y: cy, width: w, height: h,
  })
  floorItems.value.push({ ...res.data, itemType: 'floor' })
}

// ---- Select ----
function selectTable(t) { selected.value = t }
function selectFloor(f) { selected.value = f }
function deselectAll() { selected.value = null }

// ---- Edit ----
function saveTableEdit() {
  clearTimeout(saveTimer)
  saveTimer = setTimeout(async () => {
    if (!selected.value || selected.value.itemType !== 'table') return
    await tableApi.adminUpdateTable(selected.value.id, {
      seats: selected.value.seats,
      price: selected.value.price,
      is_active: selected.value.is_active,
    })
  }, 600)
}

function saveFloorEdit() {
  clearTimeout(saveTimer)
  saveTimer = setTimeout(async () => {
    if (!selected.value || selected.value.itemType !== 'floor') return
    await tableApi.adminUpdateFloorItem(selected.value.id, {
      label: selected.value.label,
      width: selected.value.width,
      height: selected.value.height,
    })
  }, 600)
}

// ---- Delete ----
async function deleteSelected() {
  if (!selected.value) return
  const name = selected.value.itemType === 'table'
    ? 'โต๊ะ ' + selected.value.number
    : selected.value.label
  if (!confirm(`ลบ "${name}" ?`)) return
  if (selected.value.itemType === 'table') {
    await tableApi.adminDeleteTable(selected.value.id)
    tables.value = tables.value.filter(t => t.id !== selected.value.id)
  } else {
    await tableApi.adminDeleteFloorItem(selected.value.id)
    floorItems.value = floorItems.value.filter(f => f.id !== selected.value.id)
  }
  selected.value = null
}

// ---- Drag Mouse ----
function startDragTable(e, table) {
  selectTable(table)
  dragging = { item: table, type: 'table' }
  const canvasRect = canvas.value.getBoundingClientRect()

  const pointerX = (e.clientX - canvasRect.left - panX.value) / zoom.value
  const pointerY = (e.clientY - canvasRect.top - panY.value) / zoom.value

  dragOffsetX = pointerX - table.pos_x
  dragOffsetY = pointerY - table.pos_y
  isDragging.value = true
  window.addEventListener('mousemove', onMouseMove)
  window.addEventListener('mouseup', onMouseUp)
}

function startDragFloor(e, item) {
  selectFloor(item)
  dragging = { item, type: 'floor' }
  const canvasRect = canvas.value.getBoundingClientRect()

  const pointerX = (e.clientX - canvasRect.left - panX.value) / zoom.value
  const pointerY = (e.clientY - canvasRect.top - panY.value) / zoom.value

  dragOffsetX = pointerX - item.pos_x
  dragOffsetY = pointerY - item.pos_y
  isDragging.value = true
  window.addEventListener('mousemove', onMouseMove)
  window.addEventListener('mouseup', onMouseUp)
}

function onMouseMove(e) {
  if (!dragging || !canvas.value) return
  const canvasRect = canvas.value.getBoundingClientRect()
  const rawX =
    (e.clientX - canvasRect.left - panX.value) / zoom.value - dragOffsetX
  const rawY =
    (e.clientY - canvasRect.top - panY.value) / zoom.value - dragOffsetY
  const size = dragging.type === 'table' ? tableSize.value : dragging.item.width
  const hgt = dragging.type === 'table' ? tableSize.value : dragging.item.height
  dragging.item.pos_x = clamp(snap(rawX, gridSize.value), 0, canvasWidth.value - size)
  dragging.item.pos_y = clamp(snap(rawY, gridSize.value), 0, canvasHeight.value - hgt)
  ghostCell.value = { x: dragging.item.pos_x, y: dragging.item.pos_y, w: size, h: hgt }
}

function onMouseUp() {
  if (!dragging) return
  ghostCell.value = null
  isDragging.value = false
  savePosition(dragging)
  dragging = null
  window.removeEventListener('mousemove', onMouseMove)
  window.removeEventListener('mouseup', onMouseUp)
}

// ---- Drag Touch ----
function startDragTableTouch(e, table) {
  selectTable(table)
  dragging = { item: table, type: 'table' }

  const touch = e.touches[0]
  const canvasRect = canvas.value.getBoundingClientRect()

  const pointerX =
    (touch.clientX - canvasRect.left - panX.value) / zoom.value

  const pointerY =
    (touch.clientY - canvasRect.top - panY.value) / zoom.value

  dragOffsetX = pointerX - table.pos_x
  dragOffsetY = pointerY - table.pos_y

  isDragging.value = true
}

function startDragFloorTouch(e, item) {
  selectFloor(item)
  dragging = { item, type: 'floor' }

  const touch = e.touches[0]
  const canvasRect = canvas.value.getBoundingClientRect()

  const pointerX =
    (touch.clientX - canvasRect.left - panX.value) / zoom.value

  const pointerY =
    (touch.clientY - canvasRect.top - panY.value) / zoom.value

  dragOffsetX = pointerX - item.pos_x
  dragOffsetY = pointerY - item.pos_y

  isDragging.value = true
}

function onTouchMove(e) {
  if (!canvas.value) return
  // 🟢 Pinch Zoom
  if (e.touches.length === 2) {
    const t1 = e.touches[0]
    const t2 = e.touches[1]

    const dx = t1.clientX - t2.clientX
    const dy = t1.clientY - t2.clientY
    const distance = Math.sqrt(dx * dx + dy * dy)

    if (!lastTouchDistance) {
      lastTouchDistance = distance
      return
    }

    const delta = distance - lastTouchDistance
    const zoomAmount = delta * 0.002
    const newZoom = clamp(zoom.value + zoomAmount, MIN_ZOOM, MAX_ZOOM)

    const rect = canvas.value.getBoundingClientRect()
    const centerX = (t1.clientX + t2.clientX) / 2 - rect.left
    const centerY = (t1.clientY + t2.clientY) / 2 - rect.top

    const worldX = (centerX - panX.value) / zoom.value
    const worldY = (centerY - panY.value) / zoom.value

    zoom.value = newZoom

    panX.value = centerX - worldX * newZoom
    panY.value = centerY - worldY * newZoom

    lastTouchDistance = distance
    return
  }

  // 🔵 Drag item
  if (dragging) {
    const touch = e.touches[0]
    const canvasRect = canvas.value.getBoundingClientRect()

    const rawX =
      (touch.clientX - canvasRect.left - panX.value) / zoom.value - dragOffsetX

    const rawY =
      (touch.clientY - canvasRect.top - panY.value) / zoom.value - dragOffsetY

    const size =
      dragging.type === 'table' ? tableSize.value : dragging.item.width

    const hgt =
      dragging.type === 'table' ? tableSize.value : dragging.item.height

    dragging.item.pos_x = clamp(snap(rawX, gridSize.value), 0, canvasWidth.value - size)
    dragging.item.pos_y = clamp(snap(rawY, gridSize.value), 0, canvasHeight.value - hgt)

    ghostCell.value = {
      x: dragging.item.pos_x,
      y: dragging.item.pos_y,
      w: size,
      h: hgt,
    }

    return
  }

  // 🔵 Pan canvas (ถ้าไม่ได้ลาก item)
  if (isPanning.value && e.touches.length === 1) {
    const touch = e.touches[0]

    const dx = touch.clientX - panTouchStartX
    const dy = touch.clientY - panTouchStartY

    panX.value = panTouchOriginX + dx
    panY.value = panTouchOriginY + dy
  }
}

function onTouchEnd() {
  if (dragging) {
    ghostCell.value = null
    isDragging.value = false
    savePosition(dragging)
    dragging = null
  }

  isPanning.value = false
  lastTouchDistance = null
}

// ---- Save Position ----
async function savePosition({ item, type }) {
  if (type === 'table') {
    await tableApi.adminUpdateTable(item.id, { pos_x: item.pos_x, pos_y: item.pos_y })
  } else {
    await tableApi.adminUpdateFloorItem(item.id, { pos_x: item.pos_x, pos_y: item.pos_y })
  }
}
</script>

<style scoped>
/* ===== Layout ===== */
.editor-page {
  display: flex;
  height: 100vh;
  overflow: hidden;
  background: var(--bg-base);
}

/* ===== Sidebar ===== */
.sidebar {
  width: 280px;
  background: var(--bg-surface);
  border-right: 1px solid var(--border-subtle);
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

@media (max-width: 768px) {
  .editor-page {
    position: relative;
  }

  .canvas-scroll {
    align-items: flex-start;
    justify-content: flex-start;
  }

  .sidebar {
    position: fixed;
    top: 0;
    left: 0;
    height: 100vh;
    width: 80%;
    max-width: 300px;
    z-index: 1000;
    transform: translateX(-100%);
    transition: transform 0.3s ease;
  }

  .sidebar.open {
    transform: translateX(0);
  }
}

.sidebar-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  z-index: 999;
}

.sidebar-header {
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-subtle);
  background: var(--bg-elevated);
}

.sidebar-title-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.sidebar-header h2 {
  font-size: 1rem;
  font-weight: 800;
  color: var(--text-primary);
  margin: 0;
  letter-spacing: 0.05em;
}

.sidebar-close-btn {
  background: none;
  border: 1px solid var(--border-subtle);
  color: var(--text-muted);
  border-radius: var(--radius-sm);
  width: 28px;
  height: 28px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.8rem;
  transition: all 0.2s;
  flex-shrink: 0;
}

.sidebar-close-btn:hover {
  border-color: var(--border);
  color: var(--text-primary);
}

.panel {
  padding: 14px 18px;
  border-bottom: 1px solid var(--border-subtle);
}

.panel-title {
  font-weight: 700;
  font-size: 0.72rem;
  color: var(--gold);
  margin-bottom: 12px;
  text-transform: uppercase;
  letter-spacing: 0.1em;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-bottom: 10px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

label {
  font-size: 0.78rem;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

input[type="text"],
input[type="number"] {
  padding: 8px 10px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-sm);
  font-size: 0.88rem;
  font-family: inherit;
  color: var(--text-primary);
  outline: none;
  width: 100%;
  transition: border-color 0.2s;
}

input[type="text"]::placeholder,
input[type="number"]::placeholder {
  color: var(--text-muted);
}

input[type="text"]:focus,
input[type="number"]:focus {
  border-color: var(--gold);
  box-shadow: 0 0 0 3px rgba(201, 168, 76, 0.1);
}

input:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  background: rgba(255, 255, 255, 0.02);
}

.toggle-row {
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
}

.toggle {
  width: 18px;
  height: 18px;
  cursor: pointer;
  accent-color: var(--gold);
}

.btn-add-table {
  width: 100%;
  padding: 9px;
  background: linear-gradient(135deg, var(--gold), var(--gold-light));
  color: #1a0e00;
  border: none;
  border-radius: var(--radius-sm);
  font-weight: 700;
  font-size: 0.88rem;
  font-family: inherit;
  cursor: pointer;
  margin-bottom: 8px;
  transition: opacity 0.2s;
}

.btn-add-table:hover {
  opacity: 0.88;
}

.btn-add-table:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.hint {
  font-size: 0.72rem;
  color: var(--text-muted);
  margin: 0;
}

.element-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
  gap: 8px;
}

.btn-element {
  padding: 10px 6px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-sm);
  cursor: pointer;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  transition: all 0.15s;
  color: var(--text-secondary);
}

.btn-element:hover {
  background: var(--gold-dim);
  border-color: var(--gold-border);
}

.el-icon {
  font-size: 1.3rem;
}

.el-label {
  font-size: 0.72rem;
  font-weight: 600;
  color: var(--text-muted);
}

.btn-element:hover .el-label {
  color: var(--gold-light);
}

.btn-delete {
  width: 100%;
  padding: 9px;
  background: rgba(239, 68, 68, 0.08);
  color: #fca5a5;
  border: 1px solid rgba(239, 68, 68, 0.2);
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-weight: 700;
  font-family: inherit;
  margin-top: 4px;
  transition: all 0.2s;
}

.btn-delete:hover {
  background: rgba(239, 68, 68, 0.15);
}

/* ===== Canvas Wrapper ===== */
.canvas-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-width: 0;
}

.canvas-toolbar {
  background: var(--bg-surface);
  border-bottom: 1px solid var(--border-subtle);
  padding: 12px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-shrink: 0;
  gap: 12px;
  flex-wrap: wrap;
}

.canvas-toolbar h3 {
  font-size: 0.95rem;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
}

.canvas-hint {
  font-size: 0.82rem;
  color: var(--text-muted);
}

.open-sidebar-btn {
  padding: 7px 14px;
  background: var(--gold-dim);
  color: var(--gold-light);
  border: 1px solid var(--gold-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-weight: 600;
  font-size: 0.85rem;
  font-family: inherit;
  transition: opacity 0.2s;
}

.open-sidebar-btn:hover {
  opacity: 0.85;
}

.btn-back {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 7px 16px;
  background: rgba(201, 168, 76, 0.06);
  color: var(--gold);
  border: 1px solid var(--gold-border);
  border-radius: 30px;
  text-decoration: none;
  font-weight: 600;
  font-size: 0.85rem;
  transition: all 0.2s;
}

.btn-back:hover {
  background: var(--gold-dim);
  border-color: var(--gold);
  box-shadow: 0 0 12px rgba(201, 168, 76, 0.2);
}

.zoom-controls {
  display: flex;
  align-items: center;
  gap: 6px;
}

.zoom-controls button {
  padding: 5px 10px;
  border: 1px solid var(--border-subtle);
  background: rgba(255, 255, 255, 0.04);
  color: var(--text-secondary);
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-weight: 600;
  font-family: inherit;
  font-size: 0.85rem;
  transition: all 0.2s;
}

.zoom-controls button:hover {
  background: var(--gold-dim);
  border-color: var(--gold-border);
  color: var(--gold-light);
}

/* ===== Canvas Scroll ===== */
.canvas-scroll {
  flex: 1;
  overflow: auto;
  -webkit-overflow-scrolling: touch;
  display: flex;
  justify-content: center;
  align-items: center;
  background: var(--bg-base);
  padding: 20px;
}

/* ===== Canvas ===== */
.canvas {
  position: relative;
  background-image: radial-gradient(circle, rgba(201, 168, 76, 0.08) 1px, transparent 1px);
  background-color: #08080f;
  cursor: default;
  user-select: none;
  box-shadow: 0 0 60px rgba(0, 0, 0, 0.8), 0 0 0 1px var(--border-subtle);
  border-radius: 4px;
  transition: transform 0.05s linear;
  touch-action: none;
  flex-shrink: 0;
}

.canvas:active {
  cursor: grabbing;
}

.grid-svg {
  position: absolute;
  top: 0;
  left: 0;
  pointer-events: none;
  z-index: 1;
}

.ghost-cell {
  position: absolute;
  background: rgba(201, 168, 76, 0.1);
  border: 2px dashed var(--gold);
  border-radius: var(--radius-md);
  pointer-events: none;
  z-index: 8;
}

/* ===== Table Nodes ===== */
.table-node {
  position: absolute;
  background: rgba(255, 255, 255, 0.04);
  border: 2px solid var(--purple);
  border-radius: var(--radius-md);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: grab;
  box-shadow: 0 0 12px rgba(139, 92, 246, 0.2);
  transition: box-shadow 0.15s, border-color 0.15s, opacity 0.15s;
  z-index: 10;
  touch-action: none;
}

.table-node:active {
  cursor: grabbing;
}

.table-node.selected {
  border-color: var(--gold);
  box-shadow: 0 0 0 3px rgba(201, 168, 76, 0.25), 0 0 20px rgba(201, 168, 76, 0.3);
}

.table-node.inactive {
  border-color: rgba(255, 255, 255, 0.12);
  opacity: 0.4;
}

.table-node.dragging {
  opacity: 0.85;
  box-shadow: 0 8px 28px rgba(201, 168, 76, 0.3);
  z-index: 20;
}

.t-name {
  font-size: 1.1rem;
  font-weight: 800;
  color: var(--text-primary);
  line-height: 1;
}

.t-price {
  font-size: 0.68rem;
  color: var(--green);
  font-weight: 700;
  margin-top: 2px;
}

.t-closed {
  font-size: 0.55rem;
  background: rgba(239, 68, 68, 0.7);
  color: white;
  padding: 1px 5px;
  border-radius: 3px;
  margin-top: 2px;
}

/* ===== Floor Items ===== */
.floor-item {
  position: absolute;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: grab;
  z-index: 5;
  border: 1px solid transparent;
  transition: border-color 0.15s;
  touch-action: none;
}

.floor-item:active {
  cursor: grabbing;
}

.floor-item.selected {
  border-color: var(--gold) !important;
  box-shadow: 0 0 0 2px rgba(201, 168, 76, 0.3);
}

.fi-stage {
  background: rgba(237, 137, 54, 0.12);
  border-color: rgba(237, 137, 54, 0.35);
}

.fi-bar {
  background: rgba(139, 92, 246, 0.12);
  border-color: rgba(139, 92, 246, 0.35);
}

.fi-restroom {
  background: rgba(16, 185, 129, 0.12);
  border-color: rgba(16, 185, 129, 0.35);
}

.fi-entrance {
  background: rgba(6, 182, 212, 0.12);
  border-color: rgba(6, 182, 212, 0.35);
}

.fi-label {
  font-size: 0.82rem;
  font-weight: 700;
  color: var(--text-primary);
  text-align: center;
  pointer-events: none;
}

.floor-background {
  position: absolute;
  border-radius: 20px;
  z-index: 0;
}

.pos-info {
  font-size: 0.7rem;
  color: var(--text-muted);
  background: rgba(255, 255, 255, 0.04);
  padding: 5px 8px;
  border-radius: var(--radius-sm);
  margin-top: 6px;
  font-family: monospace;
}

.slider {
  width: 100%;
  accent-color: var(--gold);
  margin: 4px 0;
}

.grid-labels {
  display: flex;
  justify-content: space-between;
  font-size: 0.7rem;
  color: var(--text-muted);
}

.checkbox-row {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.82rem;
  font-weight: 600;
  color: var(--text-secondary);
  cursor: pointer;
  margin-top: 6px;
}

.drag-indicator {
  color: var(--gold);
  font-weight: 600;
  margin-left: 8px;
  font-size: 0.82rem;
}

.grid-badge {
  font-size: 0.72rem;
  background: var(--gold-dim);
  color: var(--gold-light);
  padding: 2px 8px;
  border-radius: 10px;
  font-weight: 600;
  border: 1px solid var(--gold-border);
}

.canvas-title {
  font-weight: 700;
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>