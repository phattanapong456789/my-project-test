<template>
  <div class="editor-page">

    <!-- Sidebar -->
    <div class="sidebar">
      <div class="sidebar-header">
        <router-link to="/admin/reservations" class="btn-back">← กลับ</router-link>
        <h2>🗺️ ผังร้าน</h2>
      </div>

      <!-- Add Table -->
      <div class="panel">
        <div class="panel-title">🪑 เพิ่มโต๊ะ</div>
        <div class="form-row">
          <div class="form-group">
            <label>ราคา (บาท)</label>
            <input v-model.number="newTable.price" type="number" min="0" step="0.01" placeholder="0" />
          </div>
        </div>
        <button @click="addTable" class="btn-add-table" :disabled="addingTable">
          {{ addingTable ? '...' : '+ เพิ่มโต๊ะ' }}
        </button>
        <p class="hint">โต๊ะจะถูกวางกลาง canvas — ลากไปวางตามต้องการ</p>
      </div>

      <!-- Add Floor Items -->
      <div class="panel">
        <div class="panel-title">🏗️ เพิ่ม Element</div>
        <div class="element-grid">
          <button
            v-for="el in elementTypes"
            :key="el.type"
            @click="addFloorItem(el)"
            class="btn-element"
          >
            <span class="el-icon">{{ el.icon }}</span>
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
            <input v-model.number="selected.price" type="number" min="0" step="0.01" @input="saveTableEdit" />
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

      <div class="panel tips">
        <div class="panel-title">💡 วิธีใช้</div>
        <ul>
          <li>คลิกเพื่อเลือก</li>
          <li>ลากเพื่อย้ายตำแหน่ง</li>
          <li>แก้ไขได้ที่ sidebar</li>
          <li>ตำแหน่งบันทึกอัตโนมัติ</li>
        </ul>
      </div>
    </div>

    <!-- Canvas -->
    <div class="canvas-wrapper" ref="canvasWrapper">
      <div class="canvas-toolbar">
        <span class="canvas-title"> ผังร้าน</span>
        <span class="canvas-hint">{{ tables.length }} โต๊ะ · {{ floorItems.length }} elements</span>
      </div>

        <div
          class="canvas"
          ref="canvas"
          @mousedown.self="deselectAll"
          @mousemove="onMouseMove"
          @mouseup="onMouseUp"
          @touchmove.prevent="onTouchMove"
          @touchend="onTouchEnd"
        >
        <!-- Floor Items (อยู่ข้างล่างโต๊ะ) -->
        <div
          v-for="item in floorItems"
          :key="'fi-' + item.id"
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
          v-for="table in tables"
          :key="'t-' + table.id"
          class="table-node"
          :class="{
            selected: selected?.id === table.id && selected?.itemType === 'table',
            inactive: !table.is_active
          }"
          :style="{ left: table.pos_x + 'px', top: table.pos_y + 'px' }"
          @mousedown.stop="startDragTable($event, table)"
          @touchstart.stop.prevent="startDragTableTouch($event, table)"
          @click.stop="selectTable(table)"
        >
          
          <div class="t-name">{{ table.number }}</div>
          <div v-if="table.price > 0" class="t-price">฿{{ table.price.toLocaleString('th-TH', { minimumFractionDigits: 0, maximumFractionDigits: 0 }) }}</div>
          
          <div v-if="!table.is_active" class="t-closed">ปิด</div>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { tableApi } from '../api/auth'

const tables = ref([])
const floorItems = ref([])
const selected = ref(null)
const addingTable = ref(false)
const newTable = ref({ seats: 4, price: 0 })
const canvasWrapper = ref(null)
const canvas = ref(null)

// drag state
let dragging = null
let dragOffsetX = 0
let dragOffsetY = 0
let saveTimer = null

const elementTypes = [
  { type: 'stage',    icon: '🎤', label: 'เวที',     width: 200, height: 100 },

]

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

// ---- Add ----
async function addTable() {
  addingTable.value = true
  try {
    const canvasRect = canvas.value?.getBoundingClientRect()
    const cx = canvasRect ? Math.floor((canvasRect.width / 2) - 40) : 200
    const cy = canvasRect ? Math.floor((canvasRect.height / 2) - 40) : 200
    const res = await tableApi.adminCreateTable({
      seats: newTable.value.seats || 4,
      price: newTable.value.price || 0,
      pos_x: cx, pos_y: cy,
    })
    tables.value.push({ ...res.data, itemType: 'table' })
    newTable.value = { seats: 4, price: 0 }
  } finally {
    addingTable.value = false
  }
}

async function addFloorItem(el) {
  const canvasRect = canvas.value?.getBoundingClientRect()
  const cx = canvasRect ? Math.floor((canvasRect.width / 2) - el.width / 2) : 200
  const cy = canvasRect ? Math.floor((canvasRect.height / 2) - el.height / 2) : 200
  const res = await tableApi.adminCreateFloorItem({
    type: el.type, label: el.icon + ' ' + el.label,
    pos_x: cx, pos_y: cy, width: el.width, height: el.height,
  })
  floorItems.value.push({ ...res.data, itemType: 'floor' })
}

// ---- Select ----
function selectTable(t) { selected.value = t }
function selectFloor(f) { selected.value = f }
function deselectAll() { selected.value = null }

// ---- Edit with debounce ----
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
  if (!confirm(`ลบ "${selected.value.itemType === 'table' ? 'โต๊ะ ' + selected.value.number : selected.value.label}" ?`)) return
  if (selected.value.itemType === 'table') {
    await tableApi.adminDeleteTable(selected.value.id)
    tables.value = tables.value.filter(t => t.id !== selected.value.id)
  } else {
    await tableApi.adminDeleteFloorItem(selected.value.id)
    floorItems.value = floorItems.value.filter(f => f.id !== selected.value.id)
  }
  selected.value = null
}

// ---- Drag (Mouse) ----
function startDragTable(e, table) {
  dragging = { item: table, type: 'table' }
  const rect = e.currentTarget.getBoundingClientRect()
  dragOffsetX = e.clientX - rect.left
  dragOffsetY = e.clientY - rect.top
  window.addEventListener('mousemove', onMouseMove)
  window.addEventListener('mouseup', onMouseUp)
}

function startDragFloor(e, item) {
  dragging = { item, type: 'floor' }
  const rect = e.currentTarget.getBoundingClientRect()
  dragOffsetX = e.clientX - rect.left
  dragOffsetY = e.clientY - rect.top
  window.addEventListener('mousemove', onMouseMove)
  window.addEventListener('mouseup', onMouseUp)
}

function onMouseMove(e) {
  if (!dragging || !canvas.value) return
  const canvasRect = canvas.value.getBoundingClientRect()
  const x = Math.max(0, Math.min(e.clientX - canvasRect.left - dragOffsetX, canvasRect.width - 80))
  const y = Math.max(0, Math.min(e.clientY - canvasRect.top - dragOffsetY, canvasRect.height - 80))
  dragging.item.pos_x = Math.round(x)
  dragging.item.pos_y = Math.round(y)
}

function onMouseUp() {
  if (!dragging) return
  savePosition(dragging)
  dragging = null
  window.removeEventListener('mousemove', onMouseMove)
  window.removeEventListener('mouseup', onMouseUp)
}

// ---- Drag (Touch) ----
function startDragTableTouch(e, table) {
  selectTable(table)
  dragging = { item: table, type: 'table' }
  const touch = e.touches[0]
  const rect = e.currentTarget.getBoundingClientRect()
  dragOffsetX = touch.clientX - rect.left
  dragOffsetY = touch.clientY - rect.top
}

function startDragFloorTouch(e, item) {
  selectFloor(item)
  dragging = { item, type: 'floor' }
  const touch = e.touches[0]
  const rect = e.currentTarget.getBoundingClientRect()
  dragOffsetX = touch.clientX - rect.left
  dragOffsetY = touch.clientY - rect.top
}

function onTouchMove(e) {
  if (!dragging || !canvas.value) return
  const touch = e.touches[0]
  const canvasRect = canvas.value.getBoundingClientRect()
  const x = Math.max(0, Math.min(touch.clientX - canvasRect.left - dragOffsetX, canvasRect.width - 80))
  const y = Math.max(0, Math.min(touch.clientY - canvasRect.top - dragOffsetY, canvasRect.height - 80))
  dragging.item.pos_x = Math.round(x)
  dragging.item.pos_y = Math.round(y)
}

function onTouchEnd() {
  if (!dragging) return
  savePosition(dragging)
  dragging = null
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
.editor-page {
  display: flex;
  height: 100vh;
  overflow: hidden;
  background: #f0f2f8;
}

/* Sidebar */
.sidebar {
  width: 280px;
  flex-shrink: 0;
  background: white;
  box-shadow: 2px 0 12px rgba(0,0,0,0.1);
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 0;
}

.sidebar-header {
  padding: 20px 20px 12px;
  border-bottom: 1px solid #f0f0f0;
}
.btn-back { color: #888; text-decoration: none; font-size: 0.85rem; display: block; margin-bottom: 6px; }
.btn-back:hover { color: #333; }
.sidebar-header h2 { font-size: 1.2rem; color: #1a1a2e; margin: 0; }

.panel {
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
}
.panel-title { font-weight: 700; font-size: 0.85rem; color: #667eea; margin-bottom: 12px; text-transform: uppercase; letter-spacing: 0.05em; }

.form-group { display: flex; flex-direction: column; gap: 4px; margin-bottom: 10px; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; }
label { font-size: 0.8rem; font-weight: 600; color: #666; }
input[type="text"], input[type="number"] {
  padding: 8px 10px;
  border: 1.5px solid #e8e8e8;
  border-radius: 8px;
  font-size: 0.9rem;
  outline: none;
  width: 100%;
}
input:focus { border-color: #667eea; }

.toggle-row { flex-direction: row; align-items: center; justify-content: space-between; }
.toggle { width: 20px; height: 20px; cursor: pointer; }

.btn-add-table {
  width: 100%;
  padding: 9px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  margin-bottom: 8px;
}
.btn-add-table:disabled { opacity: 0.6; cursor: not-allowed; }
.hint { font-size: 0.75rem; color: #aaa; margin: 0; }

.element-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; }
.btn-element {
  padding: 10px 6px;
  background: #f8f8ff;
  border: 1.5px solid #e8e8e8;
  border-radius: 10px;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  transition: all 0.15s;
}
.btn-element:hover { background: #f0f4ff; border-color: #667eea; }
.el-icon { font-size: 1.4rem; }
.el-label { font-size: 0.75rem; font-weight: 600; color: #555; }

.btn-delete {
  width: 100%;
  padding: 9px;
  background: #fff5f5;
  color: #c53030;
  border: 1.5px solid #fed7d7;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
  margin-top: 4px;
}
.btn-delete:hover { background: #fed7d7; }

.tips ul { margin: 0; padding-left: 18px; }
.tips li { font-size: 0.8rem; color: #888; margin-bottom: 4px; }

/* Canvas */
.canvas-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.canvas-toolbar {
  background: white;
  padding: 12px 20px;
  border-bottom: 1px solid #e8e8e8;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.canvas-title { font-weight: 700; color: #1a1a2e; }
.canvas-hint { font-size: 0.85rem; color: #888; }

.canvas {
  flex: 1;
  position: relative;
  overflow: auto;
  background:
    radial-gradient(circle, #d0d4e8 1px, transparent 1px);
  background-size: 24px 24px;
  background-color: #f4f5fb;
  min-width: 900px;
  min-height: 600px;
  cursor: default;
  user-select: none;
}

/* Table Nodes */
.table-node {
  position: absolute;
  width: 50px;
  height: 50px;
  background: white;
  border: 2.5px solid #667eea;
  border-radius: 14px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: grab;
  box-shadow: 0 2px 10px rgba(102,126,234,0.15);
  transition: box-shadow 0.15s, border-color 0.15s;
  z-index: 10;
  padding: 2px;
}
.table-node:active { cursor: grabbing; }
.table-node.selected {
  border-color: #f6ad55;
  box-shadow: 0 0 0 3px rgba(246,173,85,0.4);
}
.table-node.inactive {
  border-color: #cbd5e0;
  opacity: 0.55;
}
.t-number { font-size: 1.2rem; font-weight: 800; color: #1a1a2e; line-height: 1; }
.t-name { font-size: 0.9rem; font-weight: 800; color: #1a1a2e; line-height: 1; }
.t-price { font-size: 0.55rem; color: #667eea; font-weight: 600; margin-top: 1px; line-height: 1; }
.t-seats { font-size: 0.65rem; color: #aaa; }
.t-closed { font-size: 0.5rem; background: #fed7d7; color: #c53030; padding: 1px 4px; border-radius: 4px; margin-top: 2px; }

/* Floor Items */
.floor-item {
  position: absolute;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: grab;
  z-index: 5;
  border: 2px dashed transparent;
  transition: border-color 0.15s;
}
.floor-item:active { cursor: grabbing; }
.floor-item.selected { border-color: #f6ad55 !important; }

.fi-stage    { background: rgba(237,137,54,0.15); border-color: rgba(237,137,54,0.3); }
.fi-bar      { background: rgba(159,122,234,0.15); border-color: rgba(159,122,234,0.3); }
.fi-restroom { background: rgba(72,187,120,0.15); border-color: rgba(72,187,120,0.3); }
.fi-entrance { background: rgba(66,153,225,0.15); border-color: rgba(66,153,225,0.3); }

.fi-label { font-size: 0.85rem; font-weight: 700; color: #1a1a2e; text-align: center; pointer-events: none; }
</style>
