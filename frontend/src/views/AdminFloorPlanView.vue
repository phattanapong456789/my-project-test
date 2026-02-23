<template>
  <div class="editor-page">

    <!-- Sidebar -->
    <div class="sidebar">
      <div class="sidebar-header">
        <router-link to="/dashboard" class="btn-back">← กลับหน้าหลัก</router-link>
        <h2>🗺️ ผังร้าน</h2>
      </div>

      <!-- Grid Size Control -->
      <div class="panel">
        <div class="panel-title">⚙️ Grid Settings</div>
        
        <label class="checkbox-row">
          <input type="checkbox" v-model="showGridLines" />
          <span>แสดงเส้น Grid</span>
        </label>
      </div>

      <!-- Add Table -->
      <div class="panel">
        <div class="panel-title">🪑 เพิ่มโต๊ะ</div>
        <div class="form-group">
          <label>ราคา (บาท)</label>
          <input v-model.number="newTable.price" type="number" min="0" step="1" placeholder="0" />
        </div>
        <button @click="addTable" class="btn-add-table" :disabled="addingTable">
          {{ addingTable ? '...' : '+ เพิ่มโต๊ะ' }}
        </button>
        <p class="hint">โต๊ะจะถูกวางกลาง ลากไปวางตามต้องการ</p>
      </div>

      <!-- Add Floor Items -->
      <div class="panel">
        <div class="panel-title">🏗️ เพิ่ม Element</div>
        <div class="element-grid">
          <button v-for="el in elementTypes" :key="el.type" @click="addFloorItem(el)" class="btn-element">
            
            <span class="el-label">{{ el.label }}</span>
          </button>
        </div>
      </div>

      <!-- Selected Item Settings -->
      <div class="panel" v-if="selected">
        <div class="panel-title">⚙️ ตั้งค่า</div>

        <template v-if="selected.itemType === 'table'">
          <div class="form-group">
            <label>เลขโต๊ะ</label>
            <input v-model.number="selected.number" disabled style="background: #f5f5f5;" />
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

        <button @click="deleteSelected" class="btn-delete">🗑️ ลบ</button>
      </div>

      
    </div>

    <!-- Canvas -->
    <div class="canvas-wrapper">
      <div class="canvas-toolbar">
        <div class="zoom-controls">
          <button @click="zoomOut">−</button>
          <span>{{ Math.round(zoom * 100) }}%</span>
          <button @click="zoomIn">+</button>
          <button @click="resetZoom">Reset</button>
        </div>
        <span class="canvas-title">
          ผังร้าน
        </span>
        <span class="canvas-hint">
          {{ tables.length }} โต๊ะ · {{ floorItems.length }} elements
        
        </span>
      </div>
      <div class="canvas-scroll">
        <!-- BUG 1 แก้: ย้าย event handlers มาที่ div ใน canvas-scroll แทน canvas-scroll -->
        <div
          class="canvas"
          ref="canvas"
          :style="canvasStyle"
          @mousedown="handleCanvasMouseDown"
          @mousemove="onMouseMove"
          @mouseup="onMouseUp"
          @touchmove.prevent="onTouchMove"
          @touchend="onTouchEnd"  
          @touchstart="handleCanvasTouchStart"
          @wheel.prevent="onWheelZoom"        
        >
          <!-- Grid lines overlay -->
          <svg v-if="showGridLines" class="grid-svg" :width="CANVAS_W" :height="CANVAS_H">
            <line v-for="x in gridLinesCols" :key="'v' + x"
              :x1="x" y1="0" :x2="x" :y2="CANVAS_H" stroke="#c8ccdd" stroke-width="0.5" />
            <line v-for="y in gridLinesRows" :key="'h' + y"
              x1="0" :y1="y" :x2="CANVAS_W" :y2="y" stroke="#c8ccdd" stroke-width="0.5" />
          </svg>

          <!-- Ghost cell highlight ขณะ drag -->
          <div v-if="ghostCell" class="ghost-cell" :style="{
            left: ghostCell.x + 'px', top: ghostCell.y + 'px',
            width: ghostCell.w + 'px', height: ghostCell.h + 'px',
          }" />

          <!-- Floor Items -->
          <!-- BUG 2 แก้: ย้าย floor-item ออกมาจากข้างใน ghost-cell div -->
          <div
            v-for="item in floorItems" :key="'fi-' + item.id"
            class="floor-item"
            :class="[`fi-${item.type}`, { selected: selected?.id === item.id && selected?.itemType === 'floor' }]"
            :style="{ left: item.pos_x + 'px', top: item.pos_y + 'px', width: item.width + 'px', height: item.height + 'px' }"
            @mousedown.stop="startDragFloor($event, item)"
            @touchstart.stop.prevent="startDragFloorTouch($event, item)"
            @click.stop="selectFloor(item)"
          >
            <span class="fi-label">{{ item.label }}</span>
          </div>

          <!-- Tables -->
          <div
            v-for="table in tables" :key="'t-' + table.id"
            class="table-node"
            :class="{
              selected: selected?.id === table.id && selected?.itemType === 'table',
              inactive: !table.is_active,
              dragging: isDragging && dragging?.item?.id === table.id
            }"
            :style="{ left: table.pos_x + 'px', top: table.pos_y + 'px', width: tableSize + 'px', height: tableSize + 'px' }"
            @mousedown.stop="startDragTable($event, table)"
            @touchstart.stop.prevent="startDragTableTouch($event, table)"
            @click.stop="selectTable(table)"
          >
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
  import { ref, computed, onMounted, onUnmounted } from 'vue'
  import { tableApi } from '../api/auth'

  // ---- Constants ----
    const CANVAS_W = 1200
    const CANVAS_H = 800

  // ---- State ----
    const tables = ref([])
    const floorItems = ref([])
    const selected = ref(null)
    const addingTable = ref(false)
    const newTable = ref({ price: 0 })
    const canvas = ref(null)
    const isDragging = ref(false)   // BUG 3 แก้: ลบ let isDragging ออก เหลือแค่ ref นี้
    const ghostCell = ref(null)
    const gridSize = ref(24)
    const showGridLines = ref(true)

    let dragging = null
    let dragOffsetX = 0
    let dragOffsetY = 0
    let saveTimer = null
  
  //เลื่อน-ซูม บนcanvas
    const isPanning = ref(false)
    let panStartX = 0
    let panStartY = 0
    let panOriginX = 0
    let panOriginY = 0

    let panTouchStartX = 0
    let panTouchStartY = 0
    let panTouchOriginX = 0
    let panTouchOriginY = 0

    const zoom = ref(1)
    const MIN_ZOOM = 0.5
    const MAX_ZOOM = 3

    const panX = ref(0)
    const panY = ref(0)

    let lastTouchDistance = null
      function zoomIn() {
        zoom.value = clamp(zoom.value + 0.1, MIN_ZOOM, MAX_ZOOM)
      }

      function zoomOut() {
        zoom.value = clamp(zoom.value - 0.1, MIN_ZOOM, MAX_ZOOM)
      }

      function resetZoom() {
        zoom.value = 1
      }
      function onWheelZoom(e) {
        const rect = canvas.value.getBoundingClientRect()

        const mouseX = e.clientX - rect.left
        const mouseY = e.clientY - rect.top

        const zoomAmount = -e.deltaY * 0.001
        const newZoom = clamp(zoom.value + zoomAmount, MIN_ZOOM, MAX_ZOOM)

        if (newZoom === zoom.value) return

        // คำนวณตำแหน่งก่อน zoom
        const worldX = (mouseX - panX.value) / zoom.value
        const worldY = (mouseY - panY.value) / zoom.value

        zoom.value = newZoom

        // ปรับ pan ให้ zoom เข้าไปจุดเมาส์
        panX.value = mouseX - worldX * newZoom
        panY.value = mouseY - worldY * newZoom
      }
      function startPan(e) {
        // ถ้ากำลังลาก item อยู่ ไม่ต้อง pan
        if (dragging) return

        isPanning.value = true

        panStartX = e.clientX
        panStartY = e.clientY

        panOriginX = panX.value
        panOriginY = panY.value

        window.addEventListener('mousemove', onPanMove)
        window.addEventListener('mouseup', stopPan)
      }
      function onPanMove(e) {
        if (!isPanning.value) return

        const dx = e.clientX - panStartX
        const dy = e.clientY - panStartY

        panX.value = panOriginX + dx
        panY.value = panOriginY + dy
      }
      function stopPan() {
        isPanning.value = false
        window.removeEventListener('mousemove', onPanMove)
        window.removeEventListener('mouseup', stopPan)
      }
      function handleCanvasMouseDown(e) {
        // ถ้าคลิกโดนโต๊ะ หรือ floor-item → ไม่ pan
        if (
          e.target.closest('.table-node') ||
          e.target.closest('.floor-item')
        ) {
          return
        }

        deselectAll()
        startPan(e)
      }
      function handleCanvasTouchStart(e) {
      // ถ้ามี 2 นิ้ว → pinch zoom (ปล่อยให้ onTouchMove จัดการ)
      if (e.touches.length === 2) return

      const target = e.target

      // ถ้าแตะโดนโต๊ะหรือ floor-item → ไม่ pan
      if (
        target.closest('.table-node') ||
        target.closest('.floor-item')
      ) {
        return
      }

      deselectAll()

      const touch = e.touches[0]

      isPanning.value = true
      panTouchStartX = touch.clientX
      panTouchStartY = touch.clientY
      panTouchOriginX = panX.value
      panTouchOriginY = panY.value
    }
// ---- Computed ----
  const tableSize = computed(() => gridSize.value * 2)

  const canvasStyle = computed(() => ({
    backgroundSize: `${gridSize.value}px ${gridSize.value}px`,
    width: CANVAS_W + 'px',
    height: CANVAS_H + 'px',
    transform: `translate(${panX.value}px, ${panY.value}px) scale(${zoom.value})`,
    transformOrigin: '0 0',
  }))

  const gridLinesCols = computed(() => {
    const cols = []
    for (let x = gridSize.value; x < CANVAS_W; x += gridSize.value) cols.push(x)
    return cols
  })
  const gridLinesRows = computed(() => {
    const rows = []
    for (let y = gridSize.value; y < CANVAS_H; y += gridSize.value) rows.push(y)
    return rows
  })

  const elementTypes = [
    { type: 'stage',     label: 'เวที',    width: 192, height: 96 },
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
      const cx = snap(CANVAS_W / 2 - tableSize.value / 2, gridSize.value)
      const cy = snap(CANVAS_H / 2 - tableSize.value / 2, gridSize.value)
      const res = await tableApi.adminCreateTable({
        price: newTable.value.price || 0,
        pos_x: cx, pos_y: cy,
      })
      tables.value.push({ ...res.data, itemType: 'table' })
      newTable.value = { price: 0 }
    } finally {
      addingTable.value = false
    }
  }

  // ---- Add Floor Item ----
  async function addFloorItem(el) {
    const w = snap(el.width, gridSize.value) || gridSize.value * 2
    const h = snap(el.height, gridSize.value) || gridSize.value * 2
    const cx = snap(CANVAS_W / 2 - w / 2, gridSize.value)
    const cy = snap(CANVAS_H / 2 - h / 2, gridSize.value)
    const res = await tableApi.adminCreateFloorItem({
      type: el.type, label:  el.label,
      pos_x: cx, pos_y: cy, width: w, height: h,
    })
    floorItems.value.push({ ...res.data, itemType: 'floor' })
  }

  // ---- Select ----
  function selectTable(t) { selected.value = t }
  function selectFloor(f) { selected.value = f }
  function deselectAll()  { selected.value = null }

  // ---- Edit ----
  function saveTableEdit() {
    clearTimeout(saveTimer)
    saveTimer = setTimeout(async () => {
      if (!selected.value || selected.value.itemType !== 'table') return
      await tableApi.adminUpdateTable(selected.value.id, {
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
    const hgt  = dragging.type === 'table' ? tableSize.value : dragging.item.height
    dragging.item.pos_x = clamp(snap(rawX, gridSize.value), 0, CANVAS_W - size)
    dragging.item.pos_y = clamp(snap(rawY, gridSize.value), 0, CANVAS_H - hgt)
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

        dragging.item.pos_x = clamp(snap(rawX, gridSize.value), 0, CANVAS_W - size)
        dragging.item.pos_y = clamp(snap(rawY, gridSize.value), 0, CANVAS_H - hgt)

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
  .editor-page { display: flex; height: 100vh; overflow: hidden; background: #f0f2f8; }

   .sidebar {
    width: 20%; flex-shrink: 0; background: white;
    box-shadow: 2px 0 12px rgba(0,0,0,0.1); overflow-y: auto;
    display: flex; flex-direction: column;
  } 
  .canvas {
  cursor: grab;
}

.canvas:active {
  cursor: grabbing;
}
  .editor-page { 
    display: flex; 
    height: 100vh; 
    overflow: hidden; 
    background: #f0f2f8; 
  }

  @media (max-width: 768px) {
    .editor-page { flex-direction: column; }
    .sidebar {
      width: 100%;
      height: 40vh; /* Takes top 40% of screen */
      border-right: none;
      border-bottom: 2px solid #ddd;
    }
    .canvas-wrapper { height: 60vh; }
  }

  .sidebar-header { padding: 15px 20px; border-bottom: 2px solid #f0f0f0; }

  .btn-back { color: #888; text-decoration: none; font-size: 0.85rem; display: block; margin-bottom: 6px; }
  .btn-back:hover { color: #333; }
  .sidebar-header h2 { font-size: 1.2rem; color: #1a1a2e; margin: 0; }

  .panel { padding: 12px 20px; border-bottom: 1px solid #f0f0f0; }
  /* .panel-title { font-weight: 700; font-size: 0.85rem; color: #667eea; margin-bottom: 12px; text-transform: uppercase; letter-spacing: 0.05em; } */
  .panel-title { font-weight: 700; font-size: 0.8rem; color: #667eea; margin-bottom: 10px; }

  .form-group { display: flex; flex-direction: column; gap: 4px; margin-bottom: 10px; }
  .form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; }
  label { font-size: 0.8rem; font-weight: 600; color: #666; }
  input[type="text"], input[type="number"] {
    padding: 8px 10px; border: 1.5px solid #e8e8e8;
    border-radius: 8px; font-size: 0.9rem; outline: none; width: 100%;
  }
  input:focus { border-color: #667eea; }

  .slider { width: 100%; accent-color: #667eea; margin: 4px 0; }
  .grid-labels { display: flex; justify-content: space-between; font-size: 0.7rem; color: #aaa; }
  .checkbox-row { display: flex; align-items: center; gap: 8px; font-size: 0.85rem; font-weight: 600; color: #555; cursor: pointer; margin-top: 6px; }
  .toggle-row { flex-direction: row; align-items: center; justify-content: space-between; }
  .toggle { width: 20px; height: 20px; cursor: pointer; }

  .btn-add-table {
    width: 100%; padding: 9px;
    background: linear-gradient(135deg, #667eea, #764ba2);
    color: white; border: none; border-radius: 8px;
    font-weight: 600; cursor: pointer; margin-bottom: 8px;
  }
  .btn-add-table:disabled { opacity: 0.6; cursor: not-allowed; }
  .hint { font-size: 0.75rem; color: #aaa; margin: 0; }

  .pos-info { font-size: 0.72rem; color: #888; background: #f8f8ff; padding: 6px 8px; border-radius: 6px; margin-top: 8px; font-family: monospace; }

  /* .element-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; } */

  /* Make element grid adapt to sidebar width */
  .element-grid { 
    display: grid; 
    grid-template-columns: repeat(auto-fill, minmax(80px, 1fr)); 
    gap: 8px; 
  }

  .btn-element {
    padding: 10px 6px; background: #f8f8ff; border: 1.5px solid #e8e8e8;
    border-radius: 10px; cursor: pointer; display: flex;
    flex-direction: column; align-items: center; gap: 4px; transition: all 0.15s;
  }
  .btn-element:hover { background: #f0f4ff; border-color: #667eea; }
  .el-icon { font-size: 1.4rem; }
  .el-label { font-size: 0.75rem; font-weight: 600; color: #555; }

  .btn-delete {
    width: 100%; padding: 9px; background: #fff5f5; color: #c53030;
    border: 1.5px solid #fed7d7; border-radius: 8px; cursor: pointer;
    font-weight: 600; margin-top: 4px;
  }
  .btn-delete:hover { background: #fed7d7; }
  .tips ul { margin: 0; padding-left: 18px; }
  .tips li { font-size: 0.8rem; color: #888; margin-bottom: 4px; }



  .canvas-wrapper { flex: 1; display: flex; flex-direction: column; overflow: hidden; min-width: 0; }
  .canvas-toolbar {
    background: white; padding: 12px 20px; border-bottom: 1px solid #e8e8e8;
    display: flex; justify-content: space-between; align-items: center; flex-shrink: 0;
  }
  .canvas-title { font-weight: 700; color: #1a1a2e; display: flex; align-items: center; gap: 8px; }
  .grid-badge { font-size: 0.72rem; background: #f0f4ff; color: #667eea; padding: 2px 8px; border-radius: 10px; font-weight: 600; }
  .canvas-hint { font-size: 0.85rem; color: #888; }
  .drag-indicator { color: #667eea; font-weight: 600; margin-left: 8px; }

  .canvas-scroll { 
    flex: 1; 
    overflow: auto;
    -webkit-overflow-scrolling: touch;
    display: flex;             
    justify-content: center;    
    align-items: center;       
    background: #e5e7eb;
    padding: 20px;
  }
  .canvas {
    position: relative;
    background-image: radial-gradient(circle, #c8ccdd 1px, transparent 1px);
    background-color: #ffffff;
    cursor: default; 
    user-select: none;
    box-shadow: 0 0 20px rgba(0,0,0,0.1);
    transition: transform 0.05s linear;
    touch-action: none;
    flex-shrink: 0; /* ⭐ ป้องกันถูกบีบ */
  }
  

  .grid-svg { position: absolute; top: 0; left: 0; pointer-events: none; z-index: 1; }

  .ghost-cell {
    position: absolute;
    background: rgba(102,126,234,0.12);
    border: 2px dashed #667eea;
    border-radius: 10px;
    pointer-events: none;
    z-index: 8;
  }

  .table-node {
    position: absolute; background: white; border: 2.5px solid #667eea;
    border-radius: 14px; display: flex; flex-direction: column;
    align-items: center; justify-content: center; cursor: grab;
    box-shadow: 0 2px 10px rgba(102,126,234,0.15);
    transition: box-shadow 0.15s, border-color 0.15s, opacity 0.15s;
    z-index: 10;
  }
  .table-node:active { cursor: grabbing; }
  .table-node.selected { border-color: #f6ad55; box-shadow: 0 0 0 3px rgba(246,173,85,0.4); }
  .table-node.inactive { border-color: #cbd5e0; opacity: 0.5; }
  .table-node.dragging { opacity: 0.85; box-shadow: 0 8px 28px rgba(102,126,234,0.35); z-index: 20; }

  .t-name  { font-size: 1.3rem; font-weight: 800; color: #1a1a2e; line-height: 1; }
  .t-price { font-size: 0.7rem; color: #48bb78; font-weight: 700; margin-top: 2px; }
  .t-closed { font-size: 0.58rem; background: #fed7d7; color: #c53030; padding: 1px 5px; border-radius: 4px; margin-top: 2px; }

  .floor-item {
    position: absolute; border-radius: 10px; display: flex;
    align-items: center; justify-content: center; cursor: grab;
    z-index: 5; border: 2px solid transparent; transition: border-color 0.15s;
  }
  .floor-item:active { cursor: grabbing; }
  .floor-item.selected { border-color: #f6ad55 !important; box-shadow: 0 0 0 3px rgba(246,173,85,0.3); }

  .fi-stage    { background: rgba(237,137,54,0.15); border-color: rgba(237,137,54,0.4); }
  .fi-bar      { background: rgba(159,122,234,0.15); border-color: rgba(159,122,234,0.4); }
  .fi-restroom { background: rgba(72,187,120,0.15);  border-color: rgba(72,187,120,0.4); }
  .fi-entrance { background: rgba(66,153,225,0.15);  border-color: rgba(66,153,225,0.4); }

  .fi-label { font-size: 0.85rem; font-weight: 700; color: #1a1a2e; text-align: center; pointer-events: none; }

  /* --- Items --- */
  .table-node, .floor-item {
    touch-action: none; /* Critical: Stops browser from scrolling while dragging items */
  }
  @media (max-width: 768px) {
    .canvas-scroll {
      align-items: flex-start; /* ⭐ ไม่ต้องจัดกลางแนวตั้งบนมือถือ */
      justify-content: flex-start;
    }
  }
    .floor-background {
      position: absolute;
      background: #ffffff;
      border-radius: 30px;
      box-shadow: 0 15px 60px rgba(0,0,0,0.12);
      z-index: 0;
    }
    .zoom-controls {
  display: flex;
  align-items: center;
  gap: 6px;
}

.zoom-controls button {
  padding: 4px 8px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 600;
}

.zoom-controls button:hover {
  background: #f0f4ff;
  border-color: #667eea;
}

</style> 