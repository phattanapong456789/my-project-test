<template>
  <div class="admin-page">
    <div class="admin-card">
      <div class="admin-header">
        <div>
          <h1>จัดการการจองโต๊ะ</h1>
          <p>อนุมัติหรือปฏิเสธการจองจากลูกค้า</p>
        </div>
        <div class="header-actions">
          <router-link to="/admin/floor-plan" class="btn-outline"><LayoutGrid :size="14" /> ผังร้าน</router-link>
          <router-link to="/dashboard" class="btn-back-nav"><ArrowLeft :size="14" /> กลับหน้าหลัก</router-link>
        </div>
      </div>

      <!-- Summary -->
      <div class="summary-grid">
        <div class="summary-card s-total"><div class="s-icon"><CalendarDays :size="24" /></div><div class="s-num">{{ shownTotal }}</div><div class="s-label">รายการทั้งหมด</div></div>
        <div class="summary-card s-pending"><div class="s-icon" style="color:var(--orange)"><Clock :size="24" /></div><div class="s-num">{{ shownPending }}</div><div class="s-label">รอยืนยัน</div></div>
        <div class="summary-card s-approved"><div class="s-icon" style="color:var(--green)"><CheckCircle2 :size="24" /></div><div class="s-num">{{ shownApproved }}</div><div class="s-label">ยืนยันแล้ว</div></div>
        <div class="summary-card s-tables"><div class="s-icon" style="color:var(--gold)"><Armchair :size="24" /></div><div class="s-num">{{ activeTables }}</div><div class="s-label">โต๊ะ</div></div>
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
            <option value="pending">รอยืนยัน</option>
            <option value="approved">ยืนยันแล้ว</option>
            <option value="rejected">ปฏิเสธ</option>
            <option value="cancelled">ยกเลิก</option>
          </select>
        </div>
        <button @click="clearFilter" class="btn-clear"><X :size="13" /> ล้าง</button>
      </div>

      <div v-if="loading" class="loading"><span class="spinner-gold"></span> กำลังโหลด...</div>
      <div v-else-if="reservations.length===0" class="empty">ไม่มีการจองในช่วงนี้</div>

      <template v-else>
        <div v-if="selectedIds.length>0" class="bulk-actions">
          <span class="selected-count">เลือก {{ selectedIds.length }} รายการ</span>
          <button @click="deleteSelected" class="btn-delete-bulk" :disabled="deleting"><Trash2 :size="13" /> ลบรายการที่เลือก</button>
        </div>

        <div class="res-list">
          <div class="res-list-header">
            <label class="checkbox-wrap">
              <input type="checkbox" :checked="isAllSelected" @change="toggleSelectAll" />
              <span>เลือกทั้งหมด</span>
            </label>
          </div>
          <div v-for="r in groupedReservations" :key="r.id" class="res-row" :class="`status-${r.status}`">
            <label class="checkbox-wrap res-checkbox">
              <input type="checkbox" :checked="isGroupSelected(r)" @change="toggleGroupSelection(r)" />
            </label>
            <div class="res-info">
              <div class="res-main">
                <strong>โต๊ะ {{ r.tables.map(t => t.number).join(', ') }}</strong>
                <span class="status-badge" :class="`badge-${r.status}`">
                  <Clock v-if="r.status==='pending'" :size="10" />
                  <CheckCircle2 v-if="r.status==='approved'" :size="10" />
                  <XCircle v-if="r.status==='rejected'" :size="10" />
                  <Ban v-if="r.status==='cancelled'" :size="10" />
                  {{ statusLabel(r.status) }}
                </span>
              </div>
              <div class="res-meta">
                <span><User :size="12" /> {{ r.user.name }} ({{ r.user.email }})</span>
                <span><CalendarDays :size="12" /> {{ formatDT(r.reserved_at) }}</span>
              </div>
              <div v-if="r.note" class="res-note"><FileText :size="12" /> {{ r.note }}</div>
              <div v-if="r.admin_note" class="res-admin-note"><MessageSquare :size="12" /> ทางร้าน: {{ r.admin_note }}</div>
            </div>
            <div class="res-actions" v-if="r.status==='pending'">
              <button @click="openAction(r,'approved')" class="btn-approve"><CheckCircle2 :size="13" /> อนุมัติ</button>
              <button @click="openAction(r,'rejected')" class="btn-reject"><XCircle :size="13" /> ปฏิเสธ</button>
            </div>
          </div>
        </div>
      </template>
    </div>

    <!-- Modal ลบ -->
    <div v-if="deleteConfirm" class="modal-overlay" @click.self="deleteConfirm=false">
      <div class="modal">
        <div class="modal-icon"><Trash2 :size="28" /></div>
        <h2>ยืนยันลบรายการจอง</h2>
        <p>ต้องการลบ {{ selectedIds.length }} รายการที่เลือก?</p>
        <div class="modal-actions">
          <button @click="deleteConfirm=false" class="btn-cancel">ยกเลิก</button>
          <button @click="confirmDelete" class="btn-danger" :disabled="deleting">{{ deleting?'กำลังลบ...':'ยืนยันลบ' }}</button>
        </div>
      </div>
    </div>

    <!-- Modal อนุมัติ/ปฏิเสธ -->
    <div v-if="actionTarget" class="modal-overlay" @click.self="actionTarget=null">
      <div class="modal">
        <div class="modal-icon" :class="actionStatus==='approved'?'icon-approve':'icon-reject'">
          <CheckCircle2 v-if="actionStatus==='approved'" :size="28" />
          <XCircle v-else :size="28" />
        </div>
        <h2>{{ actionStatus==='approved'?'อนุมัติการจอง':'ปฏิเสธการจอง' }}</h2>
        <p><strong>{{ actionTarget.user.name }}</strong> จอง <strong>โต๊ะ {{ actionTarget.tables.map(t => t.number).join(', ') }}</strong><br>{{ formatDT(actionTarget.reserved_at) }}</p>
        <div class="form-group" style="margin-top:14px;">
          <label>หมายเหตุถึงลูกค้า (ไม่บังคับ)</label>
          <input v-model="adminNote" type="text" placeholder="เช่น กรุณามาก่อนเวลา 15 นาที..." />
        </div>
        <div v-if="actionError" class="alert-error"><AlertTriangle :size="13" /> {{ actionError }}</div>
        <div class="modal-actions">
          <button @click="actionTarget=null" class="btn-cancel">ยกเลิก</button>
          <button @click="confirmAction" :class="actionStatus==='approved'?'btn-approve-confirm':'btn-danger'" :disabled="actioning">{{ actioning?'กำลังบันทึก...':(actionStatus==='approved'?'ยืนยันอนุมัติ':'ยืนยันปฏิเสธ') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { tableApi } from '../api/auth'
import { LayoutGrid, ArrowLeft, CalendarDays, Clock, CheckCircle2, Armchair, X, Trash2, User, XCircle, Ban, FileText, MessageSquare, AlertTriangle } from 'lucide-vue-next'

const reservations = ref([])
const summary = ref(null)
const loading = ref(false)
const filterDate = ref('')
const filterStatus = ref('')
const shownTotal = computed(()=>reservations.value.length)
const shownPending = computed(()=>reservations.value.filter(r=>r.status==='pending').length)
const shownApproved = computed(()=>reservations.value.filter(r=>r.status==='approved').length)
const activeTables = computed(()=>summary.value?.tables?.active??0)
const actionTarget = ref(null)
const actionStatus = ref('')
const adminNote = ref('')
const actioning = ref(false)
const actionError = ref(null)
const selectedIds = ref([])
const deleteConfirm = ref(false)

const groupedReservations = computed(() => {
  const groups = {}
  reservations.value.forEach(r => {
    // Generate grouping key based on user_id, reserved_at, status, and note
    const key = `${r.user.id}_${r.reserved_at}_${r.status}_${r.note || ''}`
    if (!groups[key]) {
      groups[key] = {
        id: r.id,
        ids: [r.id],
        user: r.user,
        reserved_at: r.reserved_at,
        status: r.status,
        note: r.note,
        admin_note: r.admin_note,
        tables: [r.table]
      }
    } else {
      groups[key].ids.push(r.id)
      groups[key].tables.push(r.table)
    }
  })
  return Object.values(groups).sort((a, b) => new Date(b.reserved_at) - new Date(a.reserved_at))
})
const deleting = ref(false)
const isAllSelected = computed(()=>reservations.value.length>0&&selectedIds.value.length===reservations.value.length)
function toggleSelectAll() { isAllSelected.value?selectedIds.value=[]:selectedIds.value=reservations.value.map(r=>r.id) }
function isGroupSelected(group) { return group.ids.every(id => selectedIds.value.includes(id)) }
function toggleGroupSelection(group) {
  const allSelected = isGroupSelected(group)
  if (allSelected) {
    selectedIds.value = selectedIds.value.filter(id => !group.ids.includes(id))
  } else {
    const newIds = [...selectedIds.value]
    group.ids.forEach(id => {
      if (!newIds.includes(id)) newIds.push(id)
    })
    selectedIds.value = newIds
  }
}
function deleteSelected() { if(selectedIds.value.length===0)return; deleteConfirm.value=true }
async function confirmDelete() { deleting.value=true; try{ await tableApi.adminDeleteReservationsBulk(selectedIds.value); selectedIds.value=[]; deleteConfirm.value=false; load(); loadSummary() }catch(e){ actionError.value=e.response?.data?.error||'เกิดข้อผิดพลาด' }finally{ deleting.value=false } }
onMounted(()=>{ loadSummary(); load() })
async function loadSummary() { try{ const res=await tableApi.adminGetSummary(); summary.value=res.data }catch{} }
async function load() { loading.value=true; selectedIds.value=[]; try{ const res=await tableApi.adminGetReservations(filterDate.value,filterStatus.value); reservations.value=res.data }finally{ loading.value=false } }
function clearFilter() { filterDate.value=''; filterStatus.value=''; load() }
function openAction(group,status) { actionTarget.value=group; actionStatus.value=status; adminNote.value=''; actionError.value=null }
async function confirmAction() {
  actioning.value=true; actionError.value=null;
  try {
    await Promise.all(actionTarget.value.ids.map(id => tableApi.adminUpdateStatus(id, actionStatus.value, adminNote.value)))
    
    // Update local state
    reservations.value.forEach(r => {
      if (actionTarget.value.ids.includes(r.id)) {
        r.status = actionStatus.value
        r.admin_note = adminNote.value
      }
    })
    
    actionTarget.value=null; loadSummary()
  } catch(e) { actionError.value=e.response?.data?.error||'เกิดข้อผิดพลาด' }
  finally { actioning.value=false }
}
function statusLabel(s) { return{pending:'รอยืนยัน',approved:'ยืนยันแล้ว',rejected:'ปฏิเสธ',cancelled:'ยกเลิก'}[s]||s }
function formatDT(dt) { return new Date(dt).toLocaleString('th-TH',{weekday:'short',month:'short',day:'numeric',hour:'2-digit',minute:'2-digit'}) }
</script>

<style scoped>
.admin-page { min-height:100vh; padding:32px 16px; display:flex; justify-content:center; background:var(--bg-base); }
.admin-card { background:var(--bg-surface); border:1px solid var(--border-subtle); border-radius:var(--radius-xl); padding:36px; width:100%; max-width:960px; box-shadow:var(--shadow-card); }
.admin-header { display:flex; justify-content:space-between; align-items:flex-start; margin-bottom:24px; flex-wrap:wrap; gap:16px; }
h1 { font-size:1.5rem; font-weight:800; color:var(--text-primary); }
.admin-header p { color:var(--text-muted); margin-top:4px; font-size:0.88rem; }
.header-actions { display:flex; gap:10px; flex-wrap:wrap; }
.btn-outline { display:inline-flex; align-items:center; gap:6px; padding:9px 18px; background:rgba(201,168,76,0.06); color:var(--gold); border:1px solid var(--gold-border); border-radius:30px; text-decoration:none; font-weight:600; font-size:0.88rem; transition:all 0.2s; }
.btn-outline:hover { background:var(--gold-dim); border-color:var(--gold); box-shadow:0 0 12px rgba(201,168,76,0.2); }
.btn-back-nav { display:inline-flex; align-items:center; gap:6px; padding:9px 18px; background:rgba(201,168,76,0.06); color:var(--gold); border:1px solid var(--gold-border); border-radius:30px; text-decoration:none; font-weight:600; font-size:0.88rem; transition:all 0.2s; }
.btn-back-nav:hover { background:var(--gold-dim); border-color:var(--gold); box-shadow:0 0 12px rgba(201,168,76,0.2); }

.summary-grid { display:grid; grid-template-columns:repeat(4,1fr); gap:14px; margin-bottom:24px; }
@media(max-width:600px) { .summary-grid { grid-template-columns:repeat(2,1fr); } }
.summary-card { background:var(--bg-elevated); border:1px solid var(--border-subtle); border-radius:var(--radius-md); padding:16px; text-align:center; }
.s-pending { border-color:rgba(245,158,11,0.3); } .s-approved { border-color:rgba(16,185,129,0.3); } .s-tables { border-color:var(--gold-border); }
.s-icon { display:flex; justify-content:center; margin-bottom:6px; color:var(--text-muted); }
.s-num { font-size:1.8rem; font-weight:800; color:var(--text-primary); }
.s-label { font-size:0.75rem; color:var(--text-muted); margin-top:2px; }

.filters { display:flex; gap:12px; align-items:flex-end; margin-bottom:20px; flex-wrap:wrap; }
.form-group { display:flex; flex-direction:column; gap:4px; }
label { font-size:0.78rem; font-weight:600; color:var(--text-muted); text-transform:uppercase; letter-spacing:0.05em; }
input[type="date"], select { padding:8px 12px; background:rgba(255,255,255,0.04); border:1px solid var(--border-subtle); border-radius:var(--radius-sm); font-size:0.88rem; font-family:inherit; color:var(--text-primary); outline:none; color-scheme:dark; }
input[type="date"]:focus, select:focus { border-color:var(--gold); }
select option { background:var(--bg-elevated); }
.btn-clear { display:inline-flex; align-items:center; gap:5px; padding:8px 16px; background:rgba(255,255,255,0.05); border:1px solid var(--border-subtle); border-radius:var(--radius-sm); cursor:pointer; font-size:0.85rem; color:var(--text-secondary); font-family:inherit; transition:all 0.2s; }
.btn-clear:hover { border-color:var(--border); color:var(--text-primary); }

.loading { display:flex; align-items:center; gap:10px; color:var(--text-muted); padding:40px; justify-content:center; }
.spinner-gold { width:18px; height:18px; border:2px solid var(--gold-border); border-top-color:var(--gold); border-radius:50%; animation:spin 0.8s linear infinite; }
@keyframes spin { to { transform:rotate(360deg); } }
.empty { text-align:center; color:var(--text-muted); padding:40px; }

.bulk-actions { display:flex; align-items:center; gap:16px; margin-bottom:14px; padding:12px 16px; background:rgba(239,68,68,0.08); border:1px solid rgba(239,68,68,0.2); border-radius:var(--radius-md); }
.selected-count { font-weight:700; color:#fca5a5; }
.btn-delete-bulk { display:inline-flex; align-items:center; gap:6px; padding:7px 16px; background:rgba(239,68,68,0.7); color:white; border:none; border-radius:var(--radius-sm); cursor:pointer; font-weight:600; font-size:0.88rem; font-family:inherit; }
.btn-delete-bulk:hover:not(:disabled) { background:rgba(239,68,68,0.9); }
.btn-delete-bulk:disabled { opacity:0.5; cursor:not-allowed; }

.res-list { display:flex; flex-direction:column; gap:8px; }
.res-list-header { padding:10px 16px; background:var(--bg-elevated); border-radius:var(--radius-sm); border:1px solid var(--border-subtle); }
.checkbox-wrap { display:flex; align-items:center; gap:8px; cursor:pointer; font-size:0.88rem; color:var(--text-secondary); }
.checkbox-wrap input[type="checkbox"] { width:16px; height:16px; cursor:pointer; accent-color:var(--gold); }
.res-checkbox { flex-shrink:0; }
.res-row { display:flex; justify-content:space-between; align-items:flex-start; padding:16px; gap:12px; border-radius:var(--radius-md); border:1px solid var(--border-subtle); background:var(--bg-elevated); flex-wrap:wrap; border-left-width:4px; }
.status-pending { border-left-color:var(--orange); } .status-approved { border-left-color:var(--green); } .status-rejected { border-left-color:var(--red); opacity:0.75; } .status-cancelled { opacity:0.5; }
.res-info { flex:1; min-width:0; }
.res-main { display:flex; align-items:center; gap:8px; flex-wrap:wrap; margin-bottom:6px; }
.res-main strong { color:var(--text-primary); font-size:0.95rem; font-weight:700; }
.status-badge { display:inline-flex; align-items:center; gap:4px; padding:3px 10px; border-radius:20px; font-size:0.75rem; font-weight:600; }
.badge-pending  { background:rgba(245,158,11,0.15); color:#fbbf24; border:1px solid rgba(245,158,11,0.3); }
.badge-approved { background:rgba(16,185,129,0.15); color:#6ee7b7; border:1px solid rgba(16,185,129,0.3); }
.badge-rejected { background:rgba(239,68,68,0.12); color:#fca5a5; border:1px solid rgba(239,68,68,0.3); }
.badge-cancelled{ background:rgba(255,255,255,0.06); color:var(--text-muted); border:1px solid var(--border-subtle); }
.res-meta { display:flex; gap:14px; flex-wrap:wrap; font-size:0.82rem; color:var(--text-muted); }
.res-meta span { display:inline-flex; align-items:center; gap:4px; }
.res-note { display:inline-flex; align-items:center; gap:4px; font-size:0.82rem; color:var(--text-muted); margin-top:5px; }
.res-admin-note { display:inline-flex; align-items:center; gap:4px; font-size:0.82rem; color:#c4b5fd; background:rgba(139,92,246,0.08); padding:4px 10px; border-radius:var(--radius-sm); margin-top:5px; border:1px solid rgba(139,92,246,0.2); }
.res-actions { display:flex; gap:8px; flex-shrink:0; }
.btn-approve { display:inline-flex; align-items:center; gap:5px; padding:7px 14px; background:rgba(16,185,129,0.1); color:#6ee7b7; border:1px solid rgba(16,185,129,0.3); border-radius:var(--radius-sm); cursor:pointer; font-weight:600; font-size:0.82rem; font-family:inherit; }
.btn-approve:hover { background:rgba(16,185,129,0.2); }
.btn-reject { display:inline-flex; align-items:center; gap:5px; padding:7px 14px; background:rgba(239,68,68,0.1); color:#fca5a5; border:1px solid rgba(239,68,68,0.25); border-radius:var(--radius-sm); cursor:pointer; font-weight:600; font-size:0.82rem; font-family:inherit; }
.btn-reject:hover { background:rgba(239,68,68,0.18); }

.modal-overlay { position:fixed; inset:0; background:rgba(0,0,0,0.75); backdrop-filter:blur(6px); display:flex; align-items:center; justify-content:center; z-index:100; padding:16px; }
.modal { background:var(--bg-elevated); border:1px solid var(--border); border-radius:var(--radius-xl); padding:32px; width:100%; max-width:440px; box-shadow:var(--shadow-deep),var(--shadow-gold); animation:fadeIn 0.3s ease; }
@keyframes fadeIn { from{opacity:0;transform:translateY(14px)} to{opacity:1;transform:translateY(0)} }
.modal-icon { display:inline-flex; align-items:center; justify-content:center; width:56px; height:56px; border-radius:50%; margin-bottom:12px; }
.icon-approve { background:rgba(16,185,129,0.15); color:#6ee7b7; border:1px solid rgba(16,185,129,0.3); }
.icon-reject  { background:rgba(239,68,68,0.12); color:#fca5a5; border:1px solid rgba(239,68,68,0.3); }
.modal h2 { font-size:1.2rem; color:var(--text-primary); font-weight:700; margin-bottom:10px; }
.modal p { color:var(--text-secondary); font-size:0.88rem; margin-bottom:8px; line-height:1.5; }
.alert-error { display:flex; align-items:center; gap:6px; padding:10px 14px; background:rgba(239,68,68,0.1); border:1px solid rgba(239,68,68,0.25); border-radius:var(--radius-sm); color:#fca5a5; font-size:0.88rem; margin-top:10px; }
.modal-actions { display:flex; gap:10px; justify-content:flex-end; margin-top:20px; }
.btn-cancel { padding:10px 20px; background:rgba(255,255,255,0.05); color:var(--text-secondary); border:1px solid var(--border-subtle); border-radius:var(--radius-sm); cursor:pointer; font-weight:600; font-family:inherit; }
.btn-cancel:hover { border-color:var(--border); color:var(--text-primary); }
.btn-approve-confirm { padding:10px 20px; background:rgba(16,185,129,0.8); color:white; border:none; border-radius:var(--radius-sm); cursor:pointer; font-weight:700; font-family:inherit; }
.btn-approve-confirm:hover { background:rgba(16,185,129,1); }
.btn-danger { display:inline-flex; align-items:center; gap:6px; padding:10px 20px; background:rgba(239,68,68,0.8); color:white; border:none; border-radius:var(--radius-sm); cursor:pointer; font-weight:700; font-family:inherit; }
.btn-danger:hover { background:rgba(239,68,68,1); }
.btn-approve-confirm:disabled, .btn-danger:disabled { opacity:0.5; cursor:not-allowed; }
input[type="text"] { width:100%; padding:10px 14px; background:rgba(255,255,255,0.04); border:1px solid var(--border-subtle); border-radius:var(--radius-sm); font-size:0.9rem; font-family:inherit; color:var(--text-primary); outline:none; }
input[type="text"]:focus { border-color:var(--gold); }
input[type="text"]::placeholder { color:var(--text-muted); }
</style>
