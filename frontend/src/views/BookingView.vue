<template>
  <div class="page">
    <div class="container">
      <!-- Header -->
      <div class="page-header">
        <div class="header-left">
          <router-link to="/dashboard" class="btn-back">
            <ArrowLeft :size="14" /> กลับ
          </router-link>
          <h1>จองโต๊ะ</h1>
        </div>
        <router-link to="/my-reservations" class="btn-my">
          <CalendarDays :size="14" /> การจองของฉัน
        </router-link>
      </div>

      <!-- Step 1 -->
      <div class="card">
        <h2 class="step-title"><span class="step-num">1</span> เลือกวันที่</h2>
        <div class="date-row">
          <div class="form-group">
            <label>วันที่ต้องการจอง</label>
            <input
              v-model="selectedDate"
              type="date"
              :min="today"
              @change="loadFloorPlan"
            />
          </div>
          <div class="time-notice">
            <Clock :size="22" class="notice-icon-svg" />
            <div>
              <div class="notice-title">เวลารับโต๊ะ</div>
              <div class="notice-sub">
                กรุณามาก่อน <strong>21:00 น.</strong>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Step 2 -->
      <div v-if="selectedDate" class="card floor-card">
        <div class="floor-header">
          <h2 class="step-title"><span class="step-num">2</span> เลือกโต๊ะ</h2>
          <div class="legend">
            <span class="leg"><span class="dot dot-free"></span>ว่าง</span>
            <span class="leg"><span class="dot dot-booked"></span>ไม่ว่าง</span>
            <span class="leg"><span class="dot dot-sel"></span>เลือก</span>
          </div>
        </div>
        <div v-if="loadingFloor" class="loading-floor">
          <span class="spinner-gold"></span> กำลังโหลดผัง...
        </div>
        <div
          ref="viewportRef"
          class="floor-viewport"
          @wheel.prevent="onWheel"
          @mousedown="startPan"
          @touchstart="onTouchStart"
          @touchmove.prevent="onTouchMove"
          @touchend="onTouchEnd"
        >
          <div class="floor-canvas" ref="floorCanvas" :style="canvasStyle">
            <div
              v-for="fi in floorItems"
              :key="'fi-' + fi.id"
              class="floor-item"
              :class="`fi-${fi.type}`"
              :style="{
                left: fi.pos_x + 'px',
                top: fi.pos_y + 'px',
                width: fi.width + 'px',
                height: fi.height + 'px',
              }"
            >
              {{ fi.label }}
            </div>
            <div
              v-for="t in tables"
              :key="'t-' + t.id"
              class="table-node"
              :class="{
                'is-booked': t.is_booked,
                'is-selected': selectedTableIds.includes(t.id),
              }"
              :style="{ left: t.pos_x + 'px', top: t.pos_y + 'px' }"
              @click="selectTable(t)"
            >
              <div class="t-name">{{ t.number }}</div>
              <div v-if="t.price > 0" class="t-price">
                {{
                  t.price.toLocaleString("th-TH", {
                    minimumFractionDigits: 0,
                    maximumFractionDigits: 0,
                  })
                }}฿
              </div>
              <div v-if="t.is_booked" class="t-booked-tag">ไม่ว่าง</div>
            </div>
            <div v-if="tables.length === 0" class="empty-canvas">
              ยังไม่มีโต๊ะในผัง<br /><small
                >ติดต่อแอดมินเพื่อตั้งค่าผังร้าน</small
              >
            </div>
          </div>
        </div>
        <div v-if="selectedTables.length > 0" class="selected-info">
          <div class="sel-detail">
            <Armchair :size="22" class="sel-icon" />
            <div>
              <strong
                >โต๊ะ
                {{ selectedTables.map((t) => t.number).join(", ") }}</strong
              >
              <div class="sel-sub-row">
                <span class="sel-sub"
                  >{{
                    selectedTables.reduce((sum, t) => sum + t.seats, 0)
                  }}
                  ที่นั่ง</span
                >
                <span
                  v-if="selectedTables.some((t) => t.price > 0)"
                  class="sel-price"
                  >{{
                    selectedTables
                      .reduce((sum, t) => sum + t.price, 0)
                      .toLocaleString("th-TH", {
                        minimumFractionDigits: 0,
                        maximumFractionDigits: 0,
                      })
                  }}฿</span
                >
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Step 3 -->
      <div v-if="selectedTableIds.length > 0" class="card">
        <h2 class="step-title"><span class="step-num">3</span> ยืนยันการจอง</h2>
        <div class="summary-box">
          <div class="sum-row">
            <span> <CalendarDays :size="14" /> วันที่ </span
            ><strong>{{ formatDateThai(selectedDate) }}</strong>
          </div>
          <div class="sum-row">
            <span> <Clock :size="14" /> เวลารับโต๊ะ </span
            ><strong>ก่อน 21:00 น.</strong>
          </div>
          <div class="sum-row">
            <span> <Armchair :size="14" /> โต๊ะ </span
            ><strong>{{
              selectedTables.map((t) => t.number).join(", ")
            }}</strong>
          </div>
          <div class="sum-row">
            <span> <Users :size="14" /> ที่นั่ง </span
            ><strong
              >{{
                selectedTables.reduce((sum, t) => sum + t.seats, 0)
              }}
              คน</strong
            >
          </div>
          <div v-if="selectedTables.some((t) => t.price > 0)" class="sum-row">
            <span> <Banknote :size="14" /> ราคา </span
            ><strong class="price-highlight"
              >฿{{
                selectedTables
                  .reduce((sum, t) => sum + t.price, 0)
                  .toLocaleString("th-TH", {
                    minimumFractionDigits: 0,
                    maximumFractionDigits: 0,
                  })
              }}</strong
            >
          </div>
        </div>
        <div v-if="bookError" class="alert-error">
          <AlertTriangle :size="14" /> {{ bookError }}
        </div>

        <div
          v-if="showConfirm"
          class="modal-overlay"
          @click.self="showConfirm = false"
        >
          <div class="modal">
            <div class="modal-icon">
              <Wine :size="32" />
            </div>
            <h3>ยืนยันการจองโต๊ะ?</h3>
            <p>
              โต๊ะ {{ selectedTables.map((t) => t.number).join(", ") }} ·
              {{ formatDateThai(selectedDate) }}
            </p>
            <div class="modal-btns">
              <button
                class="btn-confirm-yes"
                @click="
                  handleBook();
                  showConfirm = false;
                "
              >
                <Check :size="14" /> ยืนยัน
              </button>
              <button class="btn-confirm-no" @click="showConfirm = false">
                ยกเลิก
              </button>
            </div>
          </div>
        </div>

        <button
          class="btn-book"
          @click="showConfirm = true"
          :disabled="booking"
        >
          <span v-if="booking" class="spinner"></span>
          <GlassWater v-else :size="16" />
          {{ booking ? "กำลังจอง..." : "ยืนยันจองโต๊ะ" }}
        </button>
      </div>
    </div>

    <!-- Success Modal -->
    <div v-if="successRes" class="modal-overlay">
      <div class="modal">
        <div class="modal-icon success-glow">
          <PartyPopper :size="32" />
        </div>
        <h2>จองสำเร็จ!</h2>
        <p class="success-sub">
          รอการยืนยันจากทางร้าน<br />เมื่อยืนยันแล้ว กรุณามาก่อน
          <strong>21:00 น.</strong>
        </p>
        <div class="success-info">
          <div class="sum-row">
            <span>โต๊ะ</span
            ><strong>{{
              successRes.tables.map((t) => t.number).join(", ")
            }}</strong>
          </div>
          <div class="sum-row">
            <span>วันที่</span
            ><strong>{{
              formatDateThai(successRes.reserved_at.split("T")[0])
            }}</strong>
          </div>
          <div class="sum-row">
            <span>สถานะ</span
            ><span class="badge-pending"> <Clock :size="11" /> รอยืนยัน </span>
          </div>
        </div>
        <div class="modal-btns">
          <router-link
            to="/my-reservations"
            class="btn-book"
            style="
              text-decoration: none;
              display: flex;
              align-items: center;
              justify-content: center;
              gap: 6px;
            "
          >
            <ClipboardList :size="14" />ดูการจองของฉัน
          </router-link>
          <button @click="resetBook" class="btn-secondary">จองอีกโต๊ะ</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from "vue";
import { tableApi } from "../api/auth";
import {
  ArrowLeft,
  CalendarDays,
  Clock,
  Armchair,
  Users,
  Banknote,
  AlertTriangle,
  Wine,
  GlassWater,
  Check,
  PartyPopper,
  ClipboardList,
} from "lucide-vue-next";

const today = new Date().toISOString().split("T")[0];
const selectedDate = ref("");
const selectedTableIds = ref([]);
const note = ref("");
const tables = ref([]);
const floorItems = ref([]);
const loadingFloor = ref(false);
const booking = ref(false);
const bookError = ref(null);
const successRes = ref(null);
const showConfirm = ref(false);
const selectedTables = computed(() =>
  tables.value.filter((t) => selectedTableIds.value.includes(t.id)),
);

let lastDistance = 0;
const zoom = ref(1),
  minZoom = 0.5,
  maxZoom = 3;
const panX = ref(0),
  panY = ref(0);
const isPanning = ref(false);
let startX = 0,
  startY = 0;
const padding = 50;
const viewportRef = ref(null);
function clamp(v, mn, mx) {
  return Math.min(Math.max(v, mn), mx);
}
const panBounds = computed(() => {
  if (!viewportRef.value) return null;
  const vw = viewportRef.value.clientWidth,
    vh = viewportRef.value.clientHeight;
  return {
    minX: Math.min(0, vw - canvasWidth.value * zoom.value),
    maxX: 0,
    minY: Math.min(0, vh - canvasHeight.value * zoom.value),
    maxY: 0,
  };
});
const canvasWidth = computed(() => {
  const all = [...floorItems.value, ...tables.value];
  if (!all.length) return 900;
  return (
    Math.max(...all.map((i) => (i.pos_x || 0) + (i.width || 50))) + padding
  );
});
const canvasHeight = computed(() => {
  const all = [...floorItems.value, ...tables.value];
  if (!all.length) return 580;
  return (
    Math.max(...all.map((i) => (i.pos_y || 0) + (i.height || 50))) + padding
  );
});
const canvasStyle = computed(() => ({
  width: `${canvasWidth.value}px`,
  height: `${canvasHeight.value}px`,
  transform: `translate(${panX.value}px,${panY.value}px) scale(${zoom.value})`,
  transformOrigin: "0 0",
}));
function onWheel(e) {
  zoom.value = Math.min(
    maxZoom,
    Math.max(minZoom, zoom.value - e.deltaY * 0.001),
  );
  if (panBounds.value) {
    panX.value = clamp(panX.value, panBounds.value.minX, panBounds.value.maxX);
    panY.value = clamp(panY.value, panBounds.value.minY, panBounds.value.maxY);
  }
}
function startPan(e) {
  if (e.target.closest(".table-node")) return;
  isPanning.value = true;
  startX = e.clientX - panX.value;
  startY = e.clientY - panY.value;
  window.addEventListener("mousemove", onPanMove);
  window.addEventListener("mouseup", stopPan);
}
function onPanMove(e) {
  if (!isPanning.value || !panBounds.value) return;
  panX.value = clamp(
    e.clientX - startX,
    panBounds.value.minX,
    panBounds.value.maxX,
  );
  panY.value = clamp(
    e.clientY - startY,
    panBounds.value.minY,
    panBounds.value.maxY,
  );
}
function stopPan() {
  isPanning.value = false;
  window.removeEventListener("mousemove", onPanMove);
  window.removeEventListener("mouseup", stopPan);
}
function getDist(t) {
  const dx = t[0].clientX - t[1].clientX,
    dy = t[0].clientY - t[1].clientY;
  return Math.sqrt(dx * dx + dy * dy);
}
function onTouchStart(e) {
  if (e.touches.length === 1) {
    startX = e.touches[0].clientX - panX.value;
    startY = e.touches[0].clientY - panY.value;
  }
  if (e.touches.length === 2) lastDistance = getDist(e.touches);
}
function onTouchMove(e) {
  if (!panBounds.value) return;
  if (e.touches.length === 1) {
    panX.value = clamp(
      e.touches[0].clientX - startX,
      panBounds.value.minX,
      panBounds.value.maxX,
    );
    panY.value = clamp(
      e.touches[0].clientY - startY,
      panBounds.value.minY,
      panBounds.value.maxY,
    );
  }
  if (e.touches.length === 2) {
    const d = getDist(e.touches);
    zoom.value = Math.min(
      maxZoom,
      Math.max(minZoom, zoom.value * (d / lastDistance)),
    );
    lastDistance = d;
  }
}
function onTouchEnd() {
  lastDistance = 0;
}
async function loadFloorPlan() {
  if (!selectedDate.value) return;
  loadingFloor.value = true;
  selectedTableIds.value = [];
  try {
    const res = await tableApi.getFloorPlan(selectedDate.value);
    tables.value = res.data.tables || [];
    floorItems.value = res.data.floor_items || [];
  } catch {
    tables.value = [];
    floorItems.value = [];
  } finally {
    loadingFloor.value = false;
  }
}
function selectTable(t) {
  if (t.is_booked) return;
  const idx = selectedTableIds.value.indexOf(t.id);
  if (idx > -1) selectedTableIds.value.splice(idx, 1);
  else selectedTableIds.value.push(t.id);
}
function formatDateThai(s) {
  if (!s) return "";
  const d = s.length > 10 ? new Date(s) : new Date(s + "T00:00:00");
  return d.toLocaleDateString("th-TH", {
    weekday: "long",
    year: "numeric",
    month: "long",
    day: "numeric",
  });
}
async function handleBook() {
  if (selectedTableIds.value.length === 0) return;
  booking.value = true;
  bookError.value = null;
  const successfulTables = [];
  const errors = [];
  for (const t of selectedTables.value) {
    try {
      await tableApi.createReservation({
        table_id: t.id,
        date: selectedDate.value,
        note: note.value,
      });
      successfulTables.push(t);
    } catch (e) {
      errors.push(
        `โต๊ะ ${t.number}: ${e.response?.data?.error || "เกิดข้อผิดพลาด"}`,
      );
    }
  }
  booking.value = false;
  if (successfulTables.length > 0) {
    successRes.value = {
      tables: successfulTables,
      reserved_at: selectedDate.value,
    };
    if (errors.length > 0) bookError.value = errors.join(", ");
  } else {
    bookError.value = errors.join(", ");
  }
}
function resetBook() {
  successRes.value = null;
  selectedTableIds.value = [];
  note.value = "";
  bookError.value = null;
  loadFloorPlan();
}
</script>

<style scoped>
.page {
  min-height: 100vh;
  padding: 24px 16px;
  background: var(--bg-base);
}

.container {
  max-width: 820px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 12px;
}

.header-left {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.btn-back {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: rgba(201, 168, 76, 0.06);
  color: var(--gold);
  border: 1px solid var(--gold-border);
  border-radius: 30px;
  text-decoration: none;
  font-size: 0.85rem;
  font-weight: 600;
  transition: all 0.2s;
}

.btn-back:hover {
  background: var(--gold-dim);
  border-color: var(--gold);
  box-shadow: 0 0 12px rgba(201, 168, 76, 0.2);
}

h1 {
  font-size: 1.8rem;
  font-weight: 800;
  color: var(--text-primary);
  margin: 0;
}

.btn-my {
  padding: 10px 18px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  text-decoration: none;
  color: var(--text-secondary);
  font-weight: 600;
  font-size: 0.9rem;
  transition: all 0.2s;
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.btn-my:hover {
  border-color: var(--gold-border);
  color: var(--gold-light);
  background: var(--gold-dim);
}

.card {
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-lg);
  padding: 24px;
  margin-bottom: 20px;
}

.step-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 1rem;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 18px 0;
}

.step-num {
  width: 26px;
  height: 26px;
  background: var(--gold-dim);
  border: 1px solid var(--gold-border);
  color: var(--gold-light);
  border-radius: 50%;
  font-size: 0.8rem;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.date-row {
  display: flex;
  gap: 20px;
  align-items: center;
  flex-wrap: wrap;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

label {
  font-weight: 600;
  font-size: 0.8rem;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

input[type="date"] {
  padding: 10px 14px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-sm);
  font-size: 0.95rem;
  font-family: inherit;
  color: var(--text-primary);
  outline: none;
  transition: border-color 0.2s;
  color-scheme: dark;
}

input[type="date"]:focus {
  border-color: var(--gold);
  box-shadow: 0 0 0 3px rgba(201, 168, 76, 0.12);
}

.time-notice {
  display: flex;
  align-items: center;
  gap: 12px;
  background: rgba(245, 158, 11, 0.1);
  border: 1px solid rgba(245, 158, 11, 0.25);
  border-radius: var(--radius-md);
  padding: 12px 16px;
}

.notice-icon-svg {
  color: var(--orange);
  flex-shrink: 0;
}

.notice-title {
  font-weight: 700;
  color: var(--orange);
  font-size: 0.85rem;
}

.notice-sub {
  color: #fbbf24;
  font-size: 0.82rem;
}

.floor-card {
  padding: 24px;
}

.floor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  flex-wrap: wrap;
  gap: 8px;
}

.legend {
  display: flex;
  gap: 14px;
}

.leg {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 0.8rem;
  color: var(--text-muted);
}

.dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.dot-free {
  background: var(--purple);
}

.dot-booked {
  background: var(--red);
}

.dot-sel {
  background: var(--gold);
}

.loading-floor {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  color: var(--text-muted);
  padding: 60px;
}

.spinner-gold {
  width: 20px;
  height: 20px;
  border: 2px solid var(--gold-border);
  border-top-color: var(--gold);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.floor-viewport {
  width: 100%;
  height: 600px;
  overflow: hidden;
  position: relative;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-subtle);
  touch-action: none;
  cursor: grab;
  background: var(--bg-elevated);
}

.floor-viewport:active {
  cursor: grabbing;
}

.floor-canvas {
  position: relative;
  background: radial-gradient(
    circle,
    rgba(201, 168, 76, 0.06) 1px,
    transparent 1px
  );
  background-size: 24px 24px;
  background-color: #0b0b18;
}

.floor-item {
  position: absolute;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.8rem;
  font-weight: 700;
  pointer-events: none;
}

.fi-stage {
  background: rgba(237, 137, 54, 0.12);
  border: 1px solid rgba(237, 137, 54, 0.3);
  color: #fbbf24;
}

.fi-bar {
  background: rgba(139, 92, 246, 0.12);
  border: 1px solid rgba(139, 92, 246, 0.3);
  color: #c4b5fd;
}

.fi-restroom {
  background: rgba(16, 185, 129, 0.12);
  border: 1px solid rgba(16, 185, 129, 0.3);
  color: #6ee7b7;
}

.fi-entrance {
  background: rgba(6, 182, 212, 0.12);
  border: 1px solid rgba(6, 182, 212, 0.3);
  color: #67e8f9;
}

.table-node {
  position: absolute;
  width: 50px;
  height: 50px;
  background: rgba(255, 255, 255, 0.04);
  border: 2px solid var(--purple);
  border-radius: var(--radius-md);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 0 12px rgba(139, 92, 246, 0.2);
  transition: all 0.15s;
  z-index: 10;
  padding: 2px;
}

.table-node:hover:not(.is-booked) {
  transform: scale(1.1);
  border-color: var(--gold);
  box-shadow: 0 0 20px rgba(201, 168, 76, 0.35);
}

.table-node.is-booked {
  background: rgba(255, 255, 255, 0.02);
  border-color: rgba(255, 255, 255, 0.08);
  cursor: not-allowed;
  opacity: 0.4;
}

.table-node.is-selected {
  background: var(--gold-dim);
  border-color: var(--gold);
  box-shadow: 0 0 28px rgba(201, 168, 76, 0.5);
  transform: scale(1.12);
}

.t-name {
  font-size: 0.88rem;
  font-weight: 700;
  color: var(--text-primary);
}

.is-selected .t-name {
  color: var(--gold-light);
}

.t-price {
  font-size: 0.65rem;
  color: var(--purple);
  font-weight: 600;
  margin-top: 2px;
}

.is-selected .t-price {
  color: var(--gold-light);
}

.t-booked-tag {
  font-size: 0.5rem;
  background: rgba(239, 68, 68, 0.7);
  color: white;
  padding: 1px 4px;
  border-radius: 3px;
  margin-top: 2px;
}

.empty-canvas {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  text-align: center;
  color: var(--text-muted);
  font-size: 0.95rem;
  line-height: 1.8;
}

.selected-info {
  margin-top: 14px;
  padding: 12px 16px;
  background: var(--gold-dim);
  border: 1px solid var(--gold-border);
  border-radius: var(--radius-md);
}

.sel-detail {
  display: flex;
  align-items: center;
  gap: 10px;
}

.sel-icon {
  color: var(--gold-light);
  flex-shrink: 0;
}

.sel-detail strong {
  display: block;
  color: var(--text-primary);
}

.sel-sub-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 2px;
}

.sel-sub {
  font-size: 0.85rem;
  color: var(--text-secondary);
}

.sel-price {
  font-size: 0.9rem;
  color: var(--gold);
  font-weight: 700;
}

.summary-box {
  background: var(--bg-elevated);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 16px;
}

.sum-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.9rem;
  color: var(--text-secondary);
}

.sum-row > span {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.sum-row strong {
  color: var(--text-primary);
  font-weight: 700;
}

.price-highlight {
  color: var(--gold);
  font-size: 1rem;
  font-weight: 800;
}

.alert-error {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.25);
  border-radius: var(--radius-sm);
  color: #fca5a5;
  font-size: 0.88rem;
  margin-bottom: 14px;
}

.btn-book {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 14px;
  background: linear-gradient(135deg, var(--gold), var(--gold-light));
  color: #1a0e00;
  border: none;
  border-radius: var(--radius-md);
  font-size: 1rem;
  font-weight: 700;
  font-family: inherit;
  cursor: pointer;
  transition:
    opacity 0.2s,
    transform 0.15s,
    box-shadow 0.25s;
}

.btn-book:hover:not(:disabled) {
  opacity: 0.9;
  transform: translateY(-1px);
  box-shadow: 0 6px 24px rgba(201, 168, 76, 0.35);
}

.btn-book:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  width: 100%;
  padding: 12px;
  background: rgba(255, 255, 255, 0.05);
  color: var(--text-secondary);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  font-size: 0.95rem;
  font-weight: 600;
  font-family: inherit;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-secondary:hover {
  border-color: var(--border);
  color: var(--text-primary);
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.75);
  backdrop-filter: blur(6px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  padding: 16px;
}

.modal {
  background: var(--bg-elevated);
  border: 1px solid var(--border);
  border-radius: var(--radius-xl);
  padding: 36px;
  width: 100%;
  max-width: 420px;
  text-align: center;
  box-shadow: var(--shadow-deep), var(--shadow-gold);
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(16px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.modal-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 60px;
  height: 60px;
  background: var(--gold-dim);
  border: 1px solid var(--gold-border);
  border-radius: 50%;
  color: var(--gold-light);
  margin: 0 auto 14px;
}

.success-glow {
  box-shadow: 0 0 24px rgba(201, 168, 76, 0.3);
}

.modal h2,
.modal h3 {
  font-size: 1.3rem;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.modal p,
.success-sub {
  color: var(--text-muted);
  font-size: 0.88rem;
  margin-bottom: 16px;
  line-height: 1.5;
}

.modal-btns {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 16px;
}

.btn-confirm-yes {
  padding: 12px;
  background: linear-gradient(135deg, var(--gold), var(--gold-light));
  color: #1a0e00;
  border: none;
  border-radius: var(--radius-md);
  font-size: 0.95rem;
  font-weight: 700;
  cursor: pointer;
  font-family: inherit;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.btn-confirm-yes:hover {
  opacity: 0.88;
}

.btn-confirm-no {
  padding: 11px;
  background: none;
  color: var(--text-muted);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  font-family: inherit;
}

.success-info {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  padding: 14px;
  margin-bottom: 20px;
  text-align: left;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.badge-pending {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  background: rgba(245, 158, 11, 0.15);
  color: #fbbf24;
  border: 1px solid rgba(245, 158, 11, 0.3);
  padding: 3px 10px;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: 600;
}

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(0, 0, 0, 0.3);
  border-top-color: #1a0e00;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
