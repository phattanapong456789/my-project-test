<template>
  <div class="page">
    <div class="container">

      <!-- Header -->
      <div class="page-header">
        
          <router-link to="/dashboard" class="btn-my">← กลับ</router-link>
          <h1>🍺 จองโต๊ะ</h1>
        
        <router-link to="/my-reservations" class="btn-my">📋 การจองของฉัน</router-link>
      </div>

      <!-- Step 1: เลือกวัน -->
      <div class="card">
        <h2 class="step-title">1. เลือกวันที่</h2>
        <div class="date-row">
          <div class="form-group">
            <label>วันที่ต้องการจอง</label>
            <input v-model="selectedDate" type="date" :min="today" @change="loadFloorPlan" />
          </div>
          <div class="time-notice">
            <span class="notice-icon">🕘</span>
            <div>
              <div class="notice-title">เวลารับโต๊ะ</div>
              <div class="notice-sub">กรุณามาก่อน <strong>21:00 น.</strong></div>
            </div>
          </div>
        </div>
      </div>

      <!-- Step 2: ผังโต๊ะ -->
      <div v-if="selectedDate" class="card floor-card">
        <div class="floor-header">
          <h2 class="step-title">2. เลือกโต๊ะ</h2>
          <div class="legend">
            <span class="leg"><span class="dot dot-free"></span>ว่าง</span>
            <span class="leg"><span class="dot dot-booked"></span>ไม่ว่าง</span>
            <span class="leg"><span class="dot dot-sel"></span>เลือก</span>
          </div>
        </div>

        <div v-if="loadingFloor" class="loading-floor">กำลังโหลดผัง...</div>

        <div
          ref="viewportRef"
          class="floor-viewport"
          @wheel.prevent="onWheel"
          @mousedown="startPan"
          @touchstart="onTouchStart"
          @touchmove.prevent="onTouchMove"
          @touchend="onTouchEnd"
        >
          <div
            class="floor-canvas"
            ref="floorCanvas"
            :style="canvasStyle"
          >

            <!-- Floor Items -->
            <div
              v-for="fi in floorItems"
              :key="'fi-' + fi.id"
              class="floor-item"
              :class="`fi-${fi.type}`"
              :style="{ left: fi.pos_x + 'px', top: fi.pos_y + 'px', width: fi.width + 'px', height: fi.height + 'px' }"
            >
              {{ fi.label }}
            </div>

            <!-- Tables -->
            <div
              v-for="t in tables"
              :key="'t-' + t.id"
              class="table-node"
              :class="{
                'is-booked': t.is_booked,
                'is-selected': selectedTableId === t.id
              }"
              :style="{ left: t.pos_x + 'px', top: t.pos_y + 'px' }"
              @click="selectTable(t)"
            >
              
              <div class="t-name">{{ t.number }}</div>
              <div v-if="t.price > 0" class="t-price">{{ t.price.toLocaleString('th-TH', { minimumFractionDigits: 0, maximumFractionDigits: 0 }) }}฿</div>
              
              <div v-if="t.is_booked" class="t-booked-tag">ไม่ว่าง</div>
            </div>

            <!-- Empty state -->
            <div v-if="tables.length === 0" class="empty-canvas">
              ยังไม่มีโต๊ะในผัง<br>
              <small>ติดต่อแอดมินเพื่อตั้งค่าผังร้าน</small>
            </div>
          </div>
        </div>

        <!-- Info โต๊ะที่เลือก -->
        <div v-if="selectedTable" class="selected-info">
          <div class="sel-detail">
            <span class="sel-icon">🪑</span>
            <div>
              <strong>โต๊ะ {{ selectedTable.number }}</strong>
              <div class="sel-sub-row">
                <span class="sel-sub">{{ selectedTable.seats }} ที่นั่ง</span>
                <span v-if="selectedTable.price > 0" class="sel-price">{{ selectedTable.price.toLocaleString('th-TH', { minimumFractionDigits: 0, maximumFractionDigits: 0 }) }}฿</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Step 3: ยืนยัน -->
      <div v-if="selectedTableId" class="card">
        <h2 class="step-title">3. ยืนยันการจอง</h2>
        <div class="summary-box">
          <div class="sum-row"><span>📅 วันที่</span><strong>{{ formatDateThai(selectedDate) }}</strong></div>
          <div class="sum-row"><span>🕘 เวลารับโต๊ะ</span><strong>ก่อน 21:00 น.</strong></div>
          <div class="sum-row"><span>🪑 โต๊ะ</span><strong>{{ selectedTable?.number }}</strong></div>
          <div class="sum-row"><span>👥 ที่นั่ง</span><strong>{{ selectedTable?.seats }} คน</strong></div>
          <div v-if="selectedTable?.price > 0" class="sum-row"><span>💰 ราคา</span><strong class="price-highlight">฿{{ selectedTable.price.toLocaleString('th-TH', { minimumFractionDigits: 0, maximumFractionDigits: 0 }) }}</strong></div>
        </div>
        <div class="form-group" style="margin-top:16px">
          <label>หมายเหตุ (ไม่บังคับ)</label>
          <input v-model="note" type="text" placeholder="เช่น วันเกิด, ขอเค้ก..." />
        </div>
        <div v-if="bookError" class="alert-error">⚠️ {{ bookError }}</div>
        <!-- เพิ่ม modal -->
        <div v-if="showConfirm" class="modal-overlay">
          <div class="modal">
            <h4>ยืนยันการจองโต๊ะใช่ไหม?</h4>
            <button class="btn-book" @click="handleBook(); showConfirm = false">✅ ยืนยัน</button>
            <button class="btn-book" @click="showConfirm = false">❌ ยกเลิก</button>
          </div>
        </div>

        <button class="btn-book" @click="showConfirm = true" :disabled="booking">
          <span v-if="booking" class="spinner"></span>
          {{ booking ? 'กำลังจอง...' : '🍺 ยืนยันจองโต๊ะ' }}
        </button>
        
      </div>

    </div>

    <!-- Success Modal -->
    <div v-if="successRes" class="modal-overlay">
      <div class="modal">
        <div class="success-icon">🎉</div>
        <h2>จองสำเร็จ!</h2>
        <p class="success-sub">รอการยืนยันจากทางร้าน<br>เมื่อยืนยันแล้ว กรุณามาก่อน <strong>21:00 น.</strong></p>
        <div class="success-info">
          <div class="sum-row"><span>โต๊ะ</span><strong>{{ successRes.table.number }}</strong></div>
          <div class="sum-row"><span>วันที่</span><strong>{{ formatDateThai(successRes.reserved_at.split('T')[0]) }}</strong></div>
          <div class="sum-row"><span>สถานะ</span><span class="badge-pending">⏳ รอยืนยัน</span></div>
        </div>
        <div class="modal-btns">
          <router-link to="/my-reservations" class="btn-primary">ดูการจองของฉัน</router-link>
          <button @click="resetBook" class="btn-secondary">จองอีกโต๊ะ</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { tableApi } from '../api/auth'

const today = new Date().toISOString().split('T')[0]
const selectedDate = ref('')
const selectedTableId = ref(null)
const note = ref('')
const tables = ref([])
const floorItems = ref([])
const loadingFloor = ref(false)
const booking = ref(false)
const bookError = ref(null)
const successRes = ref(null)
const showConfirm = ref(false)
const selectedTable = computed(() => tables.value.find(t => t.id === selectedTableId.value))

let lastDistance = 0
const zoom = ref(1)
const minZoom = 0.5
const maxZoom = 3

const panX = ref(0)
const panY = ref(0)

const isPanning = ref(false)
let startX = 0
let startY = 0

const padding = 50
const viewportRef = ref(null)

function clamp(value, min, max) {
  return Math.min(Math.max(value, min), max)
}
const panBounds = computed(() => {
  if (!viewportRef.value) return null

  const vw = viewportRef.value.clientWidth
  const vh = viewportRef.value.clientHeight

  const scaledWidth = canvasWidth.value * zoom.value
  const scaledHeight = canvasHeight.value * zoom.value

  return {
    minX: Math.min(0, vw - scaledWidth),
    maxX: 0,
    minY: Math.min(0, vh - scaledHeight),
    maxY: 0
  }
})

const canvasWidth = computed(() => {
  const allItems = [
    ...floorItems.value,
    ...tables.value
  ]

  if (!allItems.length) return 900

  const maxX = Math.max(
    ...allItems.map(item =>
      (item.pos_x || 0) + (item.width || 50)
    )
  )

  return maxX + padding
})

const canvasHeight = computed(() => {
  const allItems = [
    ...floorItems.value,
    ...tables.value
  ]

  if (!allItems.length) return 580

  const maxY = Math.max(
    ...allItems.map(item =>
      (item.pos_y || 0) + (item.height || 50)
    )
  )

  return maxY + padding
})

const canvasStyle = computed(() => ({
  width: `${canvasWidth.value}px`,
  height: `${canvasHeight.value}px`,
  transform: `translate(${panX.value}px, ${panY.value}px) scale(${zoom.value})`,
  transformOrigin: '0 0'
}))
function onWheel(e) {
  const zoomSpeed = 0.001
  const newZoom = zoom.value - e.deltaY * zoomSpeed

  zoom.value = Math.min(maxZoom, Math.max(minZoom, newZoom))

  // 👇 clamp ใหม่หลัง zoom
  if (panBounds.value) {
    panX.value = clamp(panX.value, panBounds.value.minX, panBounds.value.maxX)
    panY.value = clamp(panY.value, panBounds.value.minY, panBounds.value.maxY)
  }
}
function startPan(e) {
  if (e.target.closest('.table-node')) return

  isPanning.value = true
  startX = e.clientX - panX.value
  startY = e.clientY - panY.value

  window.addEventListener('mousemove', onPanMove)
  window.addEventListener('mouseup', stopPan)
}

function onPanMove(e) {
  if (!isPanning.value || !panBounds.value) return

  const newX = e.clientX - startX
  const newY = e.clientY - startY

  panX.value = clamp(
    newX,
    panBounds.value.minX,
    panBounds.value.maxX
  )

  panY.value = clamp(
    newY,
    panBounds.value.minY,
    panBounds.value.maxY
  )
}

function stopPan() {
  isPanning.value = false
  window.removeEventListener('mousemove', onPanMove)
  window.removeEventListener('mouseup', stopPan)
}
function getDistance(touches) {
  const dx = touches[0].clientX - touches[1].clientX
  const dy = touches[0].clientY - touches[1].clientY
  return Math.sqrt(dx * dx + dy * dy)
}

function onTouchStart(e) {
  if (e.touches.length === 1) {
    startX = e.touches[0].clientX - panX.value
    startY = e.touches[0].clientY - panY.value
  }

  if (e.touches.length === 2) {
    lastDistance = getDistance(e.touches)
  }
}

function onTouchMove(e) {
  if (!panBounds.value) return

  if (e.touches.length === 1) {
    const newX = e.touches[0].clientX - startX
    const newY = e.touches[0].clientY - startY

    panX.value = clamp(
      newX,
      panBounds.value.minX,
      panBounds.value.maxX
    )

    panY.value = clamp(
      newY,
      panBounds.value.minY,
      panBounds.value.maxY
    )
  }

  if (e.touches.length === 2) {
    const newDistance = getDistance(e.touches)
    const scaleFactor = newDistance / lastDistance

    zoom.value = Math.min(
      maxZoom,
      Math.max(minZoom, zoom.value * scaleFactor)
    )

    lastDistance = newDistance
  }
}

function onTouchEnd() {
  lastDistance = 0
}

async function loadFloorPlan() {
  if (!selectedDate.value) return
  loadingFloor.value = true
  selectedTableId.value = null
  try {
    const res = await tableApi.getFloorPlan(selectedDate.value)
    tables.value = res.data.tables || []
    floorItems.value = res.data.floor_items || []
  } catch {
    tables.value = []
    floorItems.value = []
  } finally {
    loadingFloor.value = false
  }
}

function selectTable(t) {
  if (t.is_booked) return
  selectedTableId.value = t.id
}

function formatDateThai(dateStr) {
  if (!dateStr) return ''
  const d = dateStr.length > 10 ? new Date(dateStr) : new Date(dateStr + 'T00:00:00')
  return d.toLocaleDateString('th-TH', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })
}
async function handleBook() {
  if (!selectedTableId.value) return
  booking.value = true
  bookError.value = null
  try {
    const res = await tableApi.createReservation({
      table_id: selectedTableId.value,
      date: selectedDate.value,
      note: note.value,
    })
    successRes.value = res.data
  } catch (e) {
    bookError.value = e.response?.data?.error || 'เกิดข้อผิดพลาด กรุณาลองใหม่'
  } finally {
    booking.value = false
  }
}

function resetBook() {
  successRes.value = null
  selectedTableId.value = null
  note.value = ''
  bookError.value = null
  loadFloorPlan()
}
</script>

<style scoped>
.page { min-height: 100vh; padding: 24px 16px; }
.container { max-width: 800px; margin: 0 auto; }

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 12px;
}
.btn-back { color: #888; text-decoration: none; font-size: 0.9rem; display: block; margin-bottom: 4px; }
h1 { font-size: 1.8rem; color: #1a1a2e; margin: 0; }
.btn-my {
  padding: 10px 18px;
  background: white;
  border: 1.5px solid #e2e8f0;
  border-radius: 10px;
  text-decoration: none;
  color: #555;
  font-weight: 600;
  font-size: 0.9rem;
}

.card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  margin-bottom: 20px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.07);
}
.step-title { font-size: 1rem; font-weight: 700; color: #1a1a2e; margin: 0 0 16px 0; }

.date-row { display: flex; gap: 20px; align-items: center; flex-wrap: wrap; }
.form-group { display: flex; flex-direction: column; gap: 6px; }
label { font-weight: 600; font-size: 0.85rem; color: #555; }
input {
  padding: 10px 14px;
  border: 2px solid #e8e8e8;
  border-radius: 10px;
  font-size: 0.95rem;
  outline: none;
}
input:focus { border-color: #667eea; }

.time-notice {
  display: flex;
  align-items: center;
  gap: 12px;
  background: #ffc558;
  border: 1.5px solid #f3e5b9;
  border-radius: 12px;
  padding: 12px 16px;
}
.notice-icon { font-size: 1.6rem; }
.notice-title { font-weight: 700; color: #92400e; font-size: 0.9rem; }
.notice-sub { color: #92400e; font-size: 0.85rem; }

/* Floor Plan */
.floor-card { padding: 24px; }
.floor-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; flex-wrap: wrap; gap: 8px; }
.legend { display: flex; gap: 14px; flex-wrap: wrap; }
.leg { display: flex; align-items: center; gap: 5px; font-size: 0.8rem; color: #0a0a0a; }
.dot { width: 10px; height: 10px; border-radius: 50%; }
.dot-free { background: #2855d3; }
.dot-booked { background: #d41b1b; }
.dot-sel { background: #f39c29; }

.loading-floor { text-align: center; color: #888; padding: 60px; }

.floor-viewport {
  width: 100%;
  height: 600px;
  overflow: hidden;
  position: relative;
  border-radius: 12px;
  border: 2px solid #ffffff;
  touch-action: none;
  cursor: grab;
  background: #f9fafd;
}

.floor-viewport:active {
  cursor: grabbing;
}
.floor-canvas {
  position: relative;
  background:
    radial-gradient(circle, #fcfcfc 1px, transparent 1px);
  background-size: 24px 24px;
  background-color: hsl(55, 53%, 92%);
}

.floor-item {
  position: absolute;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.85rem;
  font-weight: 700;
  color: #1a1a2e;
  pointer-events: none;
}
.fi-stage    { background: rgba(237,137,54,0.18); border: 2px solid rgba(237,137,54,0.4); }
.fi-bar      { background: rgba(159,122,234,0.18); border: 2px solid rgba(159,122,234,0.4); }
.fi-restroom { background: rgba(72,187,120,0.18); border: 2px solid rgba(72,187,120,0.4); }
.fi-entrance { background: rgba(66,153,225,0.18); border: 2px solid rgba(66,153,225,0.4); }

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
  cursor: pointer;
  box-shadow: 0 2px 10px rgba(102,126,234,0.15);
  transition: all 0.15s;
  z-index: 10;
  padding: 2px;
}
.table-node:hover:not(.is-booked) {
  transform: scale(1.08);
  box-shadow: 0 4px 18px rgba(102,126,234,0.3);
}
.table-node.is-booked {
  background: #f7f7f7;
  border-color: #cbd5e0;
  cursor: not-allowed;
  opacity: 0.6;
}
.table-node.is-selected {
  background: #f6ad55;
  border-color: #f6ad55;
  color: white;
  transform: scale(1.1);
  box-shadow: 0 6px 20px rgba(246,173,85,0.4);
}
.t-num { font-size: 1.2rem; font-weight: 800; color: #1a1a2e; line-height: 1; }
.is-selected .t-num { color: white; }
.t-name { font-size: 0.9rem; color: #070707; font-weight: 600; max-width: 70px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; text-align: center; line-height: 1; }
.is-selected .t-name { color: rgba(255,255,255,0.85); }
.t-price { font-size: 0.7rem; color: #667eea; font-weight: 600; margin-top: 10%; line-height: 1; }
.is-selected .t-price { color: rgba(255,255,255,0.9); }
.t-seats { font-size: 0.62rem; color: #aaa; }
.is-selected .t-seats { color: rgba(255,255,255,0.7); }
.t-booked-tag { font-size: 0.5rem; background: #ff0000; color: #ffffff; padding: 1px 4px; border-radius: 4px; margin-top: 2px; }

.empty-canvas {
  position: absolute;
  top: 50%; left: 50%;
  transform: translate(-50%, -50%);
  text-align: center;
  color: #b9b7b7;
  font-size: 0.95rem;
  line-height: 1.8;
}

.selected-info {
  margin-top: 14px;
  padding: 12px 16px;
  background: #fff8ed;
  border-radius: 10px;
  border: 1.5px solid #fcd34d;
}
.sel-detail { display: flex; align-items: center; gap: 10px; }
.sel-icon { font-size: 1.4rem; }
.sel-detail strong { display: block; color: #1a1a2e; }
.sel-sub-row { display: flex; align-items: center; gap: 8px; margin-top: 2px; }
.sel-sub { font-size: 0.85rem; color: #d97706; }
.sel-price { font-size: 0.9rem; color: #667eea; font-weight: 600; }

/* Summary */
.summary-box {
  background: #f8f8ff;
  border-radius: 10px;
  padding: 14px 16px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.sum-row { display: flex; justify-content: space-between; align-items: center; font-size: 0.9rem; color: #555; }
.sum-row strong { color: #1a1a2e; }
.price-highlight { color: #667eea; font-size: 1rem; }

.alert-error {
  background: #fff5f5;
  color: #c53030;
  border: 1px solid #fed7d7;
  padding: 10px 14px;
  border-radius: 8px;
  font-size: 0.9rem;
  margin-top: 12px;
}

.btn-book {
  width: 100%;
  margin-top: 16px;
  padding: 14px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 1rem;
  font-weight: 700;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition: opacity 0.2s;
}
.btn-book:hover:not(:disabled) { opacity: 0.9; }
.btn-book:disabled { opacity: 0.6; cursor: not-allowed; }

/* Modal */
.modal-overlay {
  position: fixed; inset: 0;
  background: rgba(0,0,0,0.55);
  display: flex; align-items: center; justify-content: center;
  z-index: 100; padding: 16px;
}
.modal {
  background: white;
  border-radius: 20px;
  padding: 36px;
  width: 100%;
  max-width: 400px;
  text-align: center;
  box-shadow: 0 20px 60px rgba(0,0,0,0.3);
}
.success-icon { font-size: 3.5rem; margin-bottom: 10px; }
.modal h2 { font-size: 1.5rem; color: #1a1a2e; margin-bottom: 6px; }
.success-sub { color: #888; font-size: 0.9rem; line-height: 1.6; margin-bottom: 16px; }
.success-info { background: #f8f8ff; border-radius: 10px; padding: 14px; margin-bottom: 20px; text-align: left; display: flex; flex-direction: column; gap: 10px; }
.badge-pending { background: #fef3c7; color: #92400e; padding: 3px 10px; border-radius: 20px; font-size: 0.8rem; font-weight: 600; }
.modal-btns { display: flex; flex-direction: column; gap: 10px; }
.btn-primary { display: block; padding: 12px; background: linear-gradient(135deg, #667eea, #764ba2); color: white; border-radius: 10px; text-decoration: none; font-weight: 600; }
.btn-secondary { padding: 12px; background: #f0f0f0; color: #555; border: none; border-radius: 10px; font-weight: 600; cursor: pointer; }

.spinner {
  width: 16px; height: 16px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }
</style>
