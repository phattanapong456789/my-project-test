<template>
  <div class="admin-page">
    <div class="admin-card">

      <!-- Header -->
      <div class="admin-header">
        <div>
          <h1>🍺 จัดการการจองโต๊ะ</h1>
          <p>อนุมัติหรือปฏิเสธการจองจากลูกค้า</p>
        </div>
        <div class="header-actions">
          <router-link to="/admin/floor-plan" class="btn-tables">🗺️ ผังร้าน</router-link>
          <router-link to="/admin" class="btn-back">← กลับ Admin</router-link>
        </div>
      </div>

      <!-- Summary Cards -->
      <div class="summary-grid" v-if="summary">
        <div class="summary-card">
          <div class="s-icon">📅</div>
          <div class="s-num">{{ summary.today.total }}</div>
          <div class="s-label">จองวันนี้ทั้งหมด</div>
        </div>
        <div class="summary-card pending">
          <div class="s-icon">⏳</div>
          <div class="s-num">{{ summary.today.pending }}</div>
          <div class="s-label">รอยืนยัน</div>
        </div>
        <div class="summary-card approved">
          <div class="s-icon">✅</div>
          <div class="s-num">{{ summary.today.approved }}</div>
          <div class="s-label">ยืนยันแล้ว</div>
        </div>
        <div class="summary-card tables">
          <div class="s-icon">🪑</div>
          <div class="s-num">{{ summary.tables.active }}</div>
          <div class="s-label">โต๊ะที่เปิดให้บริการ</div>
        </div>
      </div>

      <!-- Filters -->
      <div class="filters">
        <div class="form-group">
          <label>กรองตามวันที่</label>
          <input v-model="filterDate" type="date" @change="load" />
        </div>
        <div class="form-group">
          <label>สถานะ</label>
          <select v-model="filterStatus" @change="load">
            <option value="">ทั้งหมด</option>
            <option value="pending">⏳ รอยืนยัน</option>
            <option value="approved">✅ ยืนยันแล้ว</option>
            <option value="rejected">❌ ปฏิเสธ</option>
            <option value="cancelled">🚫 ยกเลิก</option>
          </select>
        </div>
        <button @click="clearFilter" class="btn-clear">ล้าง</button>
      </div>

      <div v-if="loading" class="loading">กำลังโหลด...</div>

      <!-- Reservations List -->
      <div v-else-if="reservations.length === 0" class="empty">ไม่มีการจองในช่วงนี้</div>

      <div v-else class="res-list">
        <div v-for="r in reservations" :key="r.id" class="res-row" :class="`status-${r.status}`">
          <div class="res-info">
            <div class="res-main">
              <strong>{{ r.table.name }}</strong>
              <span class="zone-tag">{{ zoneLabel(r.table.zone) }}</span>
              <span class="status-badge" :class="`badge-${r.status}`">{{ statusLabel(r.status) }}</span>
            </div>
            <div class="res-meta">
              <span>👤 {{ r.user.name }} ({{ r.user.email }})</span>
              <span>📅 {{ formatDT(r.reserved_at) }}</span>
              <span>👥 {{ r.table.seats }} ที่นั่ง</span>
            </div>
            <div v-if="r.note" class="res-note">📝 {{ r.note }}</div>
            <div v-if="r.admin_note" class="res-admin-note">💬 ทางร้าน: {{ r.admin_note }}</div>
          </div>

          <div class="res-actions" v-if="r.status === 'pending'">
            <button @click="openAction(r, 'approved')" class="btn-approve">✅ อนุมัติ</button>
            <button @click="openAction(r, 'rejected')" class="btn-reject">❌ ปฏิเสธ</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal อนุมัติ/ปฏิเสธ -->
    <div v-if="actionTarget" class="modal-overlay" @click.self="actionTarget = null">
      <div class="modal">
        <h2>
          {{ actionStatus === 'approved' ? '✅ อนุมัติการจอง' : '❌ ปฏิเสธการจอง' }}
        </h2>
        <p>
          <strong>{{ actionTarget.user.name }}</strong> จอง
          <strong>{{ actionTarget.table.name }}</strong>
          วันที่ {{ formatDT(actionTarget.reserved_at) }}
        </p>
        <div class="form-group" style="margin-top: 14px;">
          <label>หมายเหตุถึงลูกค้า (ไม่บังคับ)</label>
          <input v-model="adminNote" type="text" placeholder="เช่น กรุณามาก่อนเวลา 15 นาที..." />
        </div>
        <div v-if="actionError" class="alert error">⚠️ {{ actionError }}</div>
        <div class="modal-actions">
          <button @click="actionTarget = null" class="btn-cancel">ยกเลิก</button>
          <button
            @click="confirmAction"
            :class="actionStatus === 'approved' ? 'btn-confirm-approve' : 'btn-confirm-reject'"
            :disabled="actioning"
          >
            {{ actioning ? 'กำลังบันทึก...' : (actionStatus === 'approved' ? 'ยืนยันอนุมัติ' : 'ยืนยันปฏิเสธ') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { tableApi } from '../api/auth'

const reservations = ref([])
const summary = ref(null)
const loading = ref(false)
const filterDate = ref('')
const filterStatus = ref('')

const actionTarget = ref(null)
const actionStatus = ref('')
const adminNote = ref('')
const actioning = ref(false)
const actionError = ref(null)

onMounted(() => { loadSummary(); load() })

async function loadSummary() {
  try {
    const res = await tableApi.adminGetSummary()
    summary.value = res.data
  } catch {}
}

async function load() {
  loading.value = true
  try {
    const res = await tableApi.adminGetReservations(filterDate.value, filterStatus.value)
    reservations.value = res.data
  } finally {
    loading.value = false
  }
}

function clearFilter() {
  filterDate.value = ''
  filterStatus.value = ''
  load()
}

function openAction(r, status) {
  actionTarget.value = r
  actionStatus.value = status
  adminNote.value = ''
  actionError.value = null
}

async function confirmAction() {
  actioning.value = true
  actionError.value = null
  try {
    const res = await tableApi.adminUpdateStatus(actionTarget.value.id, actionStatus.value, adminNote.value)
    const idx = reservations.value.findIndex(r => r.id === actionTarget.value.id)
    if (idx !== -1) reservations.value[idx] = res.data
    actionTarget.value = null
    loadSummary()
  } catch (e) {
    actionError.value = e.response?.data?.error || 'เกิดข้อผิดพลาด'
  } finally {
    actioning.value = false
  }
}

function statusLabel(s) {
  return { pending: '⏳ รอยืนยัน', approved: '✅ ยืนยันแล้ว', rejected: '❌ ปฏิเสธ', cancelled: '🚫 ยกเลิก' }[s] || s
}
function zoneLabel(z) {
  return { main: 'Main', vip: 'VIP', outdoor: 'กลางแจ้ง' }[z] || z
}
function formatDT(dt) {
  return new Date(dt).toLocaleString('th-TH', {
    weekday: 'short', month: 'short', day: 'numeric',
    hour: '2-digit', minute: '2-digit'
  })
}
</script>

<style scoped>
.admin-page { min-height: 100vh; padding: 32px 16px; display: flex; justify-content: center; }
.admin-card { background: white; border-radius: 20px; padding: 36px; width: 100%; max-width: 900px; box-shadow: 0 20px 60px rgba(0,0,0,0.2); }

.admin-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 24px; flex-wrap: wrap; gap: 16px; }
h1 { font-size: 1.6rem; color: #1a1a2e; }
.admin-header p { color: #888; margin-top: 4px; font-size: 0.9rem; }
.header-actions { display: flex; gap: 10px; }
.btn-tables { padding: 10px 18px; background: #f0f4ff; color: #667eea; border: 1.5px solid #c7d2fe; border-radius: 10px; text-decoration: none; font-weight: 600; font-size: 0.9rem; }
.btn-back { padding: 10px 18px; background: #f7f7f7; color: #555; border-radius: 10px; text-decoration: none; font-weight: 600; font-size: 0.9rem; }

.summary-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 14px; margin-bottom: 24px; }
@media(max-width: 600px) { .summary-grid { grid-template-columns: repeat(2,1fr); } }
.summary-card { background: #f8f8ff; border-radius: 12px; padding: 16px; text-align: center; border: 1.5px solid #e8e8e8; }
.summary-card.pending { background: #fffbeb; border-color: #fcd34d; }
.summary-card.approved { background: #f0fff4; border-color: #9ae6b4; }
.summary-card.tables { background: #f0f4ff; border-color: #c7d2fe; }
.s-icon { font-size: 1.4rem; margin-bottom: 4px; }
.s-num { font-size: 1.8rem; font-weight: 800; color: #1a1a2e; }
.s-label { font-size: 0.75rem; color: #888; }

.filters { display: flex; gap: 12px; align-items: flex-end; margin-bottom: 20px; flex-wrap: wrap; }
.form-group { display: flex; flex-direction: column; gap: 4px; }
label { font-size: 0.8rem; font-weight: 600; color: #666; }
input, select { padding: 9px 14px; border: 2px solid #e8e8e8; border-radius: 8px; font-size: 0.9rem; outline: none; }
input:focus, select:focus { border-color: #667eea; }
.btn-clear { padding: 9px 16px; background: #f0f0f0; border: none; border-radius: 8px; cursor: pointer; font-size: 0.85rem; color: #666; }

.loading, .empty { text-align: center; color: #888; padding: 40px; }

.res-list { display: flex; flex-direction: column; gap: 12px; }
.res-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-radius: 12px;
  border: 1.5px solid #e8e8e8;
  gap: 16px;
  flex-wrap: wrap;
}
.res-row.status-pending { border-left: 4px solid #f6ad55; background: #fffef5; }
.res-row.status-approved { border-left: 4px solid #48bb78; background: #f7fff9; }
.res-row.status-rejected { border-left: 4px solid #fc8181; background: #fff9f9; opacity: 0.8; }
.res-row.status-cancelled { opacity: 0.6; }

.res-info { flex: 1; min-width: 0; }
.res-main { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; margin-bottom: 6px; }
.res-main strong { color: #1a1a2e; font-size: 1rem; }
.zone-tag { font-size: 0.72rem; background: #f0f0f0; color: #888; padding: 2px 8px; border-radius: 10px; }
.status-badge { padding: 3px 10px; border-radius: 20px; font-size: 0.75rem; font-weight: 600; }
.badge-pending { background: #fef3c7; color: #92400e; }
.badge-approved { background: #f0fff4; color: #276749; }
.badge-rejected { background: #fff5f5; color: #c53030; }
.badge-cancelled { background: #f0f0f0; color: #888; }

.res-meta { display: flex; gap: 16px; flex-wrap: wrap; font-size: 0.85rem; color: #666; }
.res-note { font-size: 0.85rem; color: #888; margin-top: 6px; }
.res-admin-note { font-size: 0.85rem; color: #3730a3; background: #f0f4ff; padding: 4px 10px; border-radius: 6px; margin-top: 6px; display: inline-block; }

.res-actions { display: flex; gap: 8px; flex-shrink: 0; }
.btn-approve { padding: 8px 16px; background: #f0fff4; color: #276749; border: 1.5px solid #9ae6b4; border-radius: 8px; cursor: pointer; font-weight: 600; font-size: 0.85rem; }
.btn-approve:hover { background: #c6f6d5; }
.btn-reject { padding: 8px 16px; background: #fff5f5; color: #c53030; border: 1.5px solid #fed7d7; border-radius: 8px; cursor: pointer; font-weight: 600; font-size: 0.85rem; }
.btn-reject:hover { background: #fed7d7; }

/* Modal */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 100; padding: 16px; }
.modal { background: white; border-radius: 16px; padding: 32px; width: 100%; max-width: 440px; }
.modal h2 { margin-bottom: 10px; color: #1a1a2e; }
.modal p { color: #555; font-size: 0.9rem; }
.alert.error { background: #fff5f5; color: #c53030; border: 1px solid #fed7d7; padding: 10px 14px; border-radius: 8px; font-size: 0.9rem; margin-top: 12px; }
.modal-actions { display: flex; gap: 10px; justify-content: flex-end; margin-top: 20px; }
.btn-cancel { padding: 10px 20px; background: #f0f0f0; color: #555; border: none; border-radius: 8px; cursor: pointer; font-weight: 600; }
.btn-confirm-approve { padding: 10px 20px; background: #48bb78; color: white; border: none; border-radius: 8px; cursor: pointer; font-weight: 600; }
.btn-confirm-approve:disabled { opacity: 0.6; cursor: not-allowed; }
.btn-confirm-reject { padding: 10px 20px; background: #fc8181; color: white; border: none; border-radius: 8px; cursor: pointer; font-weight: 600; }
.btn-confirm-reject:disabled { opacity: 0.6; cursor: not-allowed; }
</style>
