<template>
  <div class="page">
    <div class="container">
      <div class="header">
        <div>
          <router-link to="/dashboard" class="btn-back">← กลับ</router-link>
          <h1>📋 การจองของฉัน</h1>
        </div>
        <router-link to="/booking" class="btn-new">+ จองโต๊ะใหม่</router-link>
      </div>

      <div v-if="loading" class="loading">กำลังโหลด...</div>

      <div v-else-if="reservations.length === 0" class="empty-state">
        <div class="empty-icon">🍺</div>
        <h2>ยังไม่มีการจอง</h2>
        <p>จองโต๊ะและมาสนุกกับเราได้เลย!</p>
        <router-link to="/booking" class="btn-primary">จองโต๊ะเลย</router-link>
      </div>

      <div v-else class="res-list">
        <div v-for="r in reservations" :key="r.id" class="res-card" :class="`s-${r.status}`">
          <div class="res-header">
            <div class="res-table">
              <span class="t-icon">🪑</span>
              <div>
                <strong>โต๊ะ {{ r.table.number }}</strong>
                <span class="seats-tag">{{ r.table.seats }} ที่นั่ง</span>
              </div>
            </div>
            <span class="status-badge" :class="`badge-${r.status}`">{{ statusLabel(r.status) }}</span>
          </div>

          <div class="res-details">
            <div class="detail-row">
              <span>📅</span><span>{{ formatDate(r.reserved_at) }}</span>
            </div>
            <div class="detail-row time-notice" v-if="r.status === 'approved'">
              <span>🕘</span><span>กรุณามารับโต๊ะ<strong>ก่อน 21:00 น.</strong></span>
            </div>
            <div v-if="r.note" class="detail-row">
              <span>📝</span><span>{{ r.note }}</span>
            </div>
            <div v-if="r.admin_note" class="admin-note">
              <span>💬</span><span>ทางร้าน: {{ r.admin_note }}</span>
            </div>
          </div>

          <div class="res-footer">
            <span class="created-at">จองเมื่อ {{ formatCreated(r.created_at) }}</span>
            <button
              v-if="r.status === 'pending'"
              @click="cancelRes(r)"
              class="btn-cancel"
              :disabled="cancelling === r.id"
            >{{ cancelling === r.id ? '...' : '✕ ยกเลิก' }}</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { tableApi } from '../api/auth'

const reservations = ref([])
const loading = ref(false)
const cancelling = ref(null)

onMounted(load)

async function load() {
  loading.value = true
  try {
    const res = await tableApi.getMyReservations()
    reservations.value = res.data || []
  } catch (e) {
    reservations.value = []
  } finally {
    loading.value = false
  }
}

async function cancelRes(r) {
  if (!confirm(`ยืนยันยกเลิกการจอง "โต๊ะ ${r.table.number}" ?`)) return
  cancelling.value = r.id
  try {
    await tableApi.cancelReservation(r.id)
    r.status = 'cancelled'
  } catch (e) {
    alert(e.response?.data?.error || 'เกิดข้อผิดพลาด')
  } finally {
    cancelling.value = null
  }
}

function statusLabel(s) {
  return { pending: '⏳ รอยืนยัน', approved: '✅ ยืนยันแล้ว', rejected: '❌ ปฏิเสธ', cancelled: '🚫 ยกเลิก' }[s] || s
}
function formatDate(dt) {
  return new Date(dt).toLocaleDateString('th-TH', { weekday: 'short', year: 'numeric', month: 'short', day: 'numeric' })
}
function formatCreated(dt) {
  return new Date(dt).toLocaleDateString('th-TH', { month: 'short', day: 'numeric', year: 'numeric' })
}
</script>

<style scoped>
.page { min-height: 100vh; padding: 24px 16px; }
.container { max-width: 680px; margin: 0 auto; }
.header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 24px; flex-wrap: wrap; gap: 12px; }
.btn-back { color: #888; text-decoration: none; font-size: 0.9rem; display: block; margin-bottom: 4px; }
h1 { font-size: 1.8rem; color: #1a1a2e; margin: 0; }
.btn-new { padding: 10px 18px; background: linear-gradient(135deg, #667eea, #764ba2); color: white; border-radius: 10px; text-decoration: none; font-weight: 600; font-size: 0.9rem; }
.loading { text-align: center; color: #888; padding: 60px; }
.empty-state { text-align: center; padding: 60px 20px; }
.empty-icon { font-size: 4rem; margin-bottom: 16px; }
.empty-state h2 { color: #1a1a2e; margin-bottom: 8px; }
.empty-state p { color: #888; margin-bottom: 24px; }
.btn-primary { display: inline-block; padding: 12px 28px; background: linear-gradient(135deg, #667eea, #764ba2); color: white; border-radius: 10px; text-decoration: none; font-weight: 600; }
.res-list { display: flex; flex-direction: column; gap: 16px; }
.res-card { background: white; border-radius: 16px; padding: 20px; box-shadow: 0 4px 16px rgba(0,0,0,0.07); border-left: 4px solid #e2e8f0; }
.s-pending { border-left-color: #f6ad55; }
.s-approved { border-left-color: #48bb78; }
.s-rejected { border-left-color: #fc8181; }
.s-cancelled { border-left-color: #cbd5e0; opacity: 0.7; }
.res-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 14px; gap: 8px; }
.res-table { display: flex; align-items: center; gap: 10px; }
.t-icon { font-size: 1.5rem; }
.res-table strong { display: block; color: #1a1a2e; }
.seats-tag { font-size: 0.75rem; color: #888; background: #f0f0f0; padding: 2px 8px; border-radius: 10px; }
.status-badge { padding: 4px 12px; border-radius: 20px; font-size: 0.8rem; font-weight: 600; white-space: nowrap; }
.badge-pending { background: #fef3c7; color: #92400e; }
.badge-approved { background: #f0fff4; color: #276749; }
.badge-rejected { background: #fff5f5; color: #c53030; }
.badge-cancelled { background: #f0f0f0; color: #888; }
.res-details { display: flex; flex-direction: column; gap: 6px; margin-bottom: 14px; }
.detail-row { display: flex; gap: 8px; font-size: 0.9rem; color: #555; }
.time-notice { background: #fffbeb; padding: 6px 10px; border-radius: 8px; color: #92400e; }
.time-notice strong { margin-left: 4px; }
.admin-note { display: flex; gap: 8px; background: #f0f4ff; border-radius: 8px; padding: 6px 10px; font-size: 0.85rem; color: #3730a3; }
.res-footer { display: flex; justify-content: space-between; align-items: center; }
.created-at { font-size: 0.8rem; color: #aaa; }
.btn-cancel { padding: 6px 14px; background: #fff5f5; color: #c53030; border: 1px solid #fed7d7; border-radius: 8px; cursor: pointer; font-size: 0.85rem; font-weight: 600; }
.btn-cancel:hover:not(:disabled) { background: #fed7d7; }
.btn-cancel:disabled { opacity: 0.6; cursor: not-allowed; }
</style>
