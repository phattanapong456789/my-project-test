<template>
  <div class="page">
    <div class="container">
      <div class="header">
        <div>
          <router-link to="/dashboard" class="btn-back"><ArrowLeft :size="14" /> กลับ</router-link>
          <h1>การจองของฉัน</h1>
        </div>
        <router-link to="/booking" class="btn-new"><Plus :size="14" /> จองโต๊ะใหม่</router-link>
      </div>

      <div v-if="loading" class="loading">
        <span class="spinner-gold"></span> กำลังโหลด...
      </div>

      <div v-else-if="reservations.length === 0" class="empty-state">
        <div class="empty-icon"><CalendarOff :size="48" /></div>
        <h2>ยังไม่มีการจอง</h2>
        <p>จองโต๊ะและมาสนุกกับเราได้เลย</p>
        <router-link to="/booking" class="btn-primary">จองโต๊ะเลย</router-link>
      </div>

      <div v-else class="res-list">
        <div v-for="r in reservations" :key="r.id" class="res-card" :class="`s-${r.status}`">
          <div class="res-header">
            <div class="res-table">
              <div class="table-icon"><Armchair :size="20" /></div>
              <div>
                <strong>โต๊ะ {{ r.table.number }}</strong>
                <span class="seats-tag"><Users :size="10" /> {{ r.table.seats }} ที่นั่ง</span>
              </div>
            </div>
            <span class="status-badge" :class="`badge-${r.status}`">
              <Clock v-if="r.status==='pending'" :size="11" />
              <CheckCircle2 v-if="r.status==='approved'" :size="11" />
              <XCircle v-if="r.status==='rejected'" :size="11" />
              <Ban v-if="r.status==='cancelled'" :size="11" />
              {{ statusLabel(r.status) }}
            </span>
          </div>

          <div class="res-details">
            <div class="detail-row"><CalendarDays :size="13" /><span>{{ formatDate(r.reserved_at) }}</span></div>
            <div class="detail-row time-notice" v-if="r.status === 'approved'">
              <Clock :size="13" /><span>กรุณามารับโต๊ะ<strong> ก่อน 21:00 น.</strong></span>
            </div>
            <div v-if="r.note" class="detail-row"><FileText :size="13" /><span>{{ r.note }}</span></div>
            <div v-if="r.admin_note" class="admin-note"><MessageSquare :size="13" /><span>ทางร้าน: {{ r.admin_note }}</span></div>
          </div>

          <div class="res-footer">
            <span class="created-at">จองเมื่อ {{ formatCreated(r.created_at) }}</span>
            <button v-if="r.status === 'pending'" @click="cancelRes(r)" class="btn-cancel" :disabled="cancelling === r.id">
              <X :size="12" /> {{ cancelling === r.id ? '...' : 'ยกเลิก' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { tableApi } from '../api/auth'
import { ArrowLeft, Plus, CalendarOff, Armchair, Users, Clock, CheckCircle2, XCircle, Ban, CalendarDays, FileText, MessageSquare, X } from 'lucide-vue-next'

const reservations = ref([])
const loading = ref(false)
const cancelling = ref(null)
onMounted(load)

async function load() {
  loading.value = true
  try { const res = await tableApi.getMyReservations(); reservations.value = res.data || [] }
  catch { reservations.value = [] }
  finally { loading.value = false }
}
async function cancelRes(r) {
  if (!confirm(`ยืนยันยกเลิกการจอง "โต๊ะ ${r.table.number}" ?`)) return
  cancelling.value = r.id
  try { await tableApi.cancelReservation(r.id); r.status = 'cancelled' }
  catch (e) { alert(e.response?.data?.error || 'เกิดข้อผิดพลาด') }
  finally { cancelling.value = null }
}
function statusLabel(s) { return { pending:'รอยืนยัน', approved:'ยืนยันแล้ว', rejected:'ปฏิเสธ', cancelled:'ยกเลิก' }[s] || s }
function formatDate(dt) { return new Date(dt).toLocaleDateString('th-TH', { weekday:'short', year:'numeric', month:'short', day:'numeric' }) }
function formatCreated(dt) { return new Date(dt).toLocaleDateString('th-TH', { month:'short', day:'numeric', year:'numeric' }) }
</script>

<style scoped>
.page { min-height:100vh; padding:24px 16px; background:var(--bg-base); }
.container { max-width:680px; margin:0 auto; }
.header { display:flex; justify-content:space-between; align-items:flex-start; margin-bottom:28px; flex-wrap:wrap; gap:12px; }
.btn-back { display:inline-flex; align-items:center; gap:6px; padding:8px 16px; background:rgba(201,168,76,0.06); color:var(--gold); border:1px solid var(--gold-border); border-radius:30px; text-decoration:none; font-size:0.85rem; font-weight:600; margin-bottom:6px; transition:all 0.2s; }
.btn-back:hover { background:var(--gold-dim); border-color:var(--gold); box-shadow:0 0 12px rgba(201,168,76,0.2); }
h1 { font-size:1.8rem; font-weight:800; color:var(--text-primary); margin:0; }
.btn-new { display:inline-flex; align-items:center; gap:6px; padding:10px 20px; background:linear-gradient(135deg,var(--gold),var(--gold-light)); color:#1a0e00; border-radius:var(--radius-md); text-decoration:none; font-weight:700; font-size:0.9rem; transition:opacity 0.2s,transform 0.15s; }
.btn-new:hover { opacity:0.88; transform:translateY(-1px); }

.loading { display:flex; align-items:center; justify-content:center; gap:10px; color:var(--text-muted); padding:60px; }
.spinner-gold { width:20px; height:20px; border:2px solid var(--gold-border); border-top-color:var(--gold); border-radius:50%; animation:spin 0.8s linear infinite; }
@keyframes spin { to { transform:rotate(360deg); } }

.empty-state { text-align:center; padding:60px 20px; }
.empty-icon { color:var(--text-muted); margin-bottom:16px; display:flex; justify-content:center; }
.empty-state h2 { color:var(--text-primary); margin-bottom:8px; font-size:1.3rem; }
.empty-state p { color:var(--text-muted); margin-bottom:24px; }
.btn-primary { display:inline-flex; align-items:center; gap:6px; padding:12px 28px; background:linear-gradient(135deg,var(--gold),var(--gold-light)); color:#1a0e00; border-radius:var(--radius-md); text-decoration:none; font-weight:700; transition:opacity 0.2s; }
.btn-primary:hover { opacity:0.88; }

.res-list { display:flex; flex-direction:column; gap:14px; }
.res-card { background:var(--bg-surface); border:1px solid var(--border-subtle); border-radius:var(--radius-lg); padding:20px; border-left:4px solid var(--border-subtle); transition:border-color 0.2s; }
.s-pending { border-left-color:var(--orange); } .s-approved { border-left-color:var(--green); } .s-rejected { border-left-color:var(--red); } .s-cancelled { opacity:0.6; }

.res-header { display:flex; justify-content:space-between; align-items:center; margin-bottom:14px; gap:8px; flex-wrap:wrap; }
.res-table { display:flex; align-items:center; gap:12px; }
.table-icon { width:40px; height:40px; border-radius:var(--radius-sm); background:var(--bg-elevated); border:1px solid var(--border-subtle); display:flex; align-items:center; justify-content:center; color:var(--text-muted); flex-shrink:0; }
.res-table strong { display:block; color:var(--text-primary); font-weight:700; }
.seats-tag { display:inline-flex; align-items:center; gap:3px; font-size:0.75rem; color:var(--text-muted); background:rgba(255,255,255,0.06); border:1px solid var(--border-subtle); padding:2px 8px; border-radius:10px; margin-top:2px; }

.status-badge { display:inline-flex; align-items:center; gap:4px; padding:4px 12px; border-radius:20px; font-size:0.78rem; font-weight:600; white-space:nowrap; }
.badge-pending  { background:rgba(245,158,11,0.15); color:#fbbf24; border:1px solid rgba(245,158,11,0.3); }
.badge-approved { background:rgba(16,185,129,0.15); color:#6ee7b7; border:1px solid rgba(16,185,129,0.3); }
.badge-rejected { background:rgba(239,68,68,0.15);  color:#fca5a5; border:1px solid rgba(239,68,68,0.3); }
.badge-cancelled{ background:rgba(255,255,255,0.06); color:var(--text-muted); border:1px solid var(--border-subtle); }

.res-details { display:flex; flex-direction:column; gap:6px; margin-bottom:14px; }
.detail-row { display:flex; gap:8px; align-items:flex-start; font-size:0.88rem; color:var(--text-secondary); }
.detail-row svg { flex-shrink:0; margin-top:2px; color:var(--text-muted); }
.time-notice { background:rgba(245,158,11,0.08); border:1px solid rgba(245,158,11,0.2); padding:6px 10px; border-radius:var(--radius-sm); color:#fbbf24; }
.time-notice svg { color:#fbbf24; }
.admin-note { display:flex; gap:8px; align-items:flex-start; background:rgba(139,92,246,0.08); border:1px solid rgba(139,92,246,0.2); border-radius:var(--radius-sm); padding:6px 10px; font-size:0.85rem; color:#c4b5fd; }
.admin-note svg { flex-shrink:0; margin-top:2px; }

.res-footer { display:flex; justify-content:space-between; align-items:center; }
.created-at { font-size:0.78rem; color:var(--text-muted); }
.btn-cancel { display:inline-flex; align-items:center; gap:4px; padding:5px 14px; background:rgba(239,68,68,0.08); color:#fca5a5; border:1px solid rgba(239,68,68,0.2); border-radius:var(--radius-sm); cursor:pointer; font-size:0.82rem; font-weight:600; font-family:inherit; transition:all 0.2s; }
.btn-cancel:hover:not(:disabled) { background:rgba(239,68,68,0.15); }
.btn-cancel:disabled { opacity:0.5; cursor:not-allowed; }
</style>
