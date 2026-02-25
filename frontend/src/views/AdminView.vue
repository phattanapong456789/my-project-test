<template>
  <div class="admin-page">
    <div class="admin-card">
      <div class="admin-header">
        <div>
          <h1>จัดการผู้ใช้งาน</h1>
          <p>{{ users.length }} บัญชีทั้งหมด · Admin: {{ adminUsers.length }} · ผู้ใช้: {{ regularUsers.length }}</p>
        </div>
        <div class="header-actions">
          <button @click="showAddModal = true" class="btn-add">
            <UserPlus :size="14" /> เพิ่ม User
          </button>
          <router-link to="/dashboard" class="btn-back-nav">
            <ArrowLeft :size="14" /> กลับ
          </router-link>
        </div>
      </div>

      <div v-if="loading" class="loading"><span class="spinner-gold"></span> กำลังโหลด...</div>
      <div v-if="error" class="alert-error">
        <AlertTriangle :size="14" /> {{ error }}
      </div>

      <div class="tabs-header">
        <button @click="activeTab = 'admin'" :class="['tab-button', { active: activeTab === 'admin' }]">
          <Crown :size="13" /> Admin ({{ adminUsers.length }})
        </button>
        <button @click="activeTab = 'user'" :class="['tab-button', { active: activeTab === 'user' }]">
          <User :size="13" /> ผู้ใช้งาน ({{ regularUsers.length }})
        </button>
      </div>

      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>#</th>
              <th>ชื่อ</th>
              <th>อีเมล</th>
              <th>บทบาท</th>
              <th>สมัครเมื่อ</th>
              <th>จัดการ</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(u, index) in displayedUsers" :key="u.id">
              <td class="td-id">{{ index + 1 }}</td>
              <td class="td-name">
                <div class="user-avatar">{{ u.name[0]?.toUpperCase() }}</div>
                {{ u.name }}
              </td>
              <td class="td-email">{{ u.email }}</td>
              <td>
                <span class="role-badge" :class="u.role === 'admin' ? 'badge-admin' : 'badge-user'">
                  <Crown v-if="u.role === 'admin'" :size="10" />
                  <User v-else :size="10" />
                  {{ u.role === 'admin' ? 'Admin' : 'User' }}
                </span>
              </td>
              <td class="td-date">{{ formatDate(u.created_at) }}</td>
              <td class="td-actions">
                <button v-if="u.id !== currentUserID" @click="toggleRole(u)" class="btn-action btn-role"
                  :title="u.role === 'admin' ? 'เปลี่ยนเป็น User' : 'เปลี่ยนเป็น Admin'">
                  <User v-if="u.role === 'admin'" :size="13" />
                  <Crown v-else :size="13" />
                </button>
                <button @click="openChangePasswordModal(u)" class="btn-action btn-pw" title="เปลี่ยนรหัสผ่าน">
                  <KeyRound :size="13" />
                </button>
                <button v-if="u.id !== currentUserID" @click="confirmDelete(u)" class="btn-action btn-del" title="ลบ">
                  <Trash2 :size="13" />
                </button>
                <span v-if="u.id === currentUserID" class="self-label">คุณ</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal เพิ่ม -->
    <div v-if="showAddModal" class="modal-overlay" @click.self="showAddModal = false">
      <div class="modal">
        <h2>เพิ่มผู้ใช้ใหม่</h2>
        <form @submit.prevent="handleAddUser">
          <div class="form-group"><label>ชื่อ</label><input v-model="newUser.name" placeholder="ชื่อ" required
              minlength="2" /></div>
          <div class="form-group"><label>อีเมล</label><input v-model="newUser.email" type="email"
              placeholder="email@example.com" required /></div>
          <div class="form-group"><label>รหัสผ่าน</label><input v-model="newUser.password" type="password"
              placeholder="อย่างน้อย 6 ตัว" required minlength="6" /></div>
          <div class="form-group">
            <label>Role</label>
            <select v-model="newUser.role">
              <option value="user">User</option>
              <option value="admin">Admin</option>
            </select>
          </div>
          <div v-if="addError" class="alert-error">
            <AlertTriangle :size="13" /> {{ addError }}
          </div>
          <div class="modal-actions">
            <button type="button" @click="showAddModal = false" class="btn-cancel">ยกเลิก</button>
            <button type="submit" class="btn-confirm" :disabled="adding">{{ adding ? 'กำลังเพิ่ม...' : 'เพิ่ม User'
              }}</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Modal ยืนยันลบ -->
    <div v-if="deleteTarget" class="modal-overlay" @click.self="deleteTarget = null">
      <div class="modal">
        <div class="modal-warn-icon">
          <AlertTriangle :size="28" />
        </div>
        <h2>ยืนยันการลบ</h2>
        <p>ต้องการลบ <strong>{{ deleteTarget.name }}</strong> ออกจากระบบ?</p>
        <p class="warning-text">การกระทำนี้ไม่สามารถยกเลิกได้</p>
        <div class="modal-actions">
          <button @click="deleteTarget = null" class="btn-cancel">ยกเลิก</button>
          <button @click="handleDelete" class="btn-danger">
            <Trash2 :size="13" /> ลบ
          </button>
        </div>
      </div>
    </div>

    <!-- Modal เปลี่ยน Role -->
    <div v-if="roleChangeTarget" class="modal-overlay" @click.self="roleChangeTarget = null">
      <div class="modal">
        <h2>ยืนยันการเปลี่ยน Role</h2>
        <p>เปลี่ยนบทบาทของ <strong>{{ roleChangeTarget.name }}</strong> เป็น <strong>{{
          roleChangeTarget.role === 'admin' ? 'User' :'Admin' }}</strong>?</p>
        <div v-if="roleChangeError" class="alert-error">
          <AlertTriangle :size="13" /> {{ roleChangeError }}
        </div>
        <div v-if="roleChangeSuccess" class="alert-success">
          <CheckCircle2 :size="13" /> {{ roleChangeSuccess }}
        </div>
        <div class="modal-actions">
          <button @click="roleChangeTarget = null" class="btn-cancel" :disabled="roleChangeLoading">ยกเลิก</button>
          <button @click="confirmChangeRole" class="btn-confirm" :disabled="roleChangeLoading">{{ roleChangeLoading ?
            'กำลังบันทึก...' : 'ยืนยัน' }}</button>
        </div>
      </div>
    </div>

    <!-- Modal เปลี่ยนรหัสผ่าน -->
    <div v-if="changePasswordTarget" class="modal-overlay" @click.self="changePasswordTarget = null">
      <div class="modal">
        <h2>เปลี่ยนรหัสผ่าน</h2>
        <p class="modal-sub">ผู้ใช้: <strong>{{ changePasswordTarget.name }}</strong></p>
        <form @submit.prevent="confirmChangePassword">
          <div class="form-group"><label>รหัสผ่านใหม่</label><input v-model="changePasswordForm.newPassword"
              type="password" placeholder="อย่างน้อย 6 ตัว" required minlength="6" /></div>
          <div class="form-group"><label>ยืนยันรหัสผ่าน</label><input v-model="changePasswordForm.confirmPassword"
              type="password" placeholder="ยืนยันรหัสผ่าน" required minlength="6" /></div>
          <div v-if="changePasswordError" class="alert-error">
            <AlertTriangle :size="13" /> {{ changePasswordError }}
          </div>
          <div v-if="changePasswordSuccess" class="alert-success">
            <CheckCircle2 :size="13" /> {{ changePasswordSuccess }}
          </div>
          <div class="modal-actions">
            <button type="button" @click="changePasswordTarget = null" class="btn-cancel">ยกเลิก</button>
            <button type="submit" class="btn-confirm" :disabled="changePasswordLoading">{{ changePasswordLoading ?
              'กำลังบันทึก...' : 'ยืนยัน' }}</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { adminApi } from '../api/auth'
import { UserPlus, ArrowLeft, Crown, User, KeyRound, Trash2, AlertTriangle, CheckCircle2 } from 'lucide-vue-next'

const users = ref([])
const loading = ref(false)
const error = ref(null)
const currentUserID = ref(Number(localStorage.getItem('userID')))
const showAddModal = ref(false)
const adding = ref(false)
const addError = ref(null)
const newUser = ref({ name: '', email: '', password: '', role: 'user' })
const deleteTarget = ref(null)
const roleChangeTarget = ref(null)
const roleChangeError = ref('')
const roleChangeSuccess = ref('')
const roleChangeLoading = ref(false)
const activeTab = ref('admin')
const changePasswordTarget = ref(null)
const changePasswordForm = ref({ newPassword: '', confirmPassword: '' })
const changePasswordError = ref('')
const changePasswordSuccess = ref('')
const changePasswordLoading = ref(false)

const adminUsers = computed(() => users.value.filter(u => u.role === 'admin').sort((a, b) => new Date(a.created_at) - new Date(b.created_at)))
const regularUsers = computed(() => users.value.filter(u => u.role === 'user').sort((a, b) => new Date(a.created_at) - new Date(b.created_at)))
const displayedUsers = computed(() => activeTab.value === 'admin' ? adminUsers.value : regularUsers.value)

onMounted(fetchUsers)

async function fetchUsers() {
  loading.value = true
  try { const res = await adminApi.getUsers(); users.value = res.data; const me = res.data.find(u => u.email === localStorage.getItem('email')); if (me) currentUserID.value = me.id }
  catch { error.value = 'ไม่สามารถโหลดข้อมูลได้' }
  finally { loading.value = false }
}
function toggleRole(u) { roleChangeTarget.value = u; roleChangeError.value = ''; roleChangeSuccess.value = '' }
async function confirmChangeRole() {
  roleChangeLoading.value = true; roleChangeError.value = ''; roleChangeSuccess.value = ''
  try { const nr = roleChangeTarget.value.role === 'admin' ? 'user' : 'admin'; const res = await adminApi.updateRole(roleChangeTarget.value.id, nr); const idx = users.value.findIndex(u => u.id === roleChangeTarget.value.id); if (idx !== -1) users.value[idx] = res.data.user; roleChangeSuccess.value = 'เปลี่ยนบทบาทสำเร็จ'; setTimeout(() => { roleChangeTarget.value = null }, 1500) }
  catch (e) { roleChangeError.value = e.response?.data?.error || 'เกิดข้อผิดพลาด' }
  finally { roleChangeLoading.value = false }
}
function confirmDelete(u) { deleteTarget.value = u }
async function handleDelete() { try { await adminApi.deleteUser(deleteTarget.value.id); users.value = users.value.filter(u => u.id !== deleteTarget.value.id); deleteTarget.value = null } catch (e) { alert(e.response?.data?.error || 'เกิดข้อผิดพลาด') } }
async function handleAddUser() { adding.value = true; addError.value = null; try { const res = await adminApi.createUser(newUser.value); users.value.unshift(res.data); showAddModal.value = false; newUser.value = { name: '', email: '', password: '', role: 'user' } } catch (e) { addError.value = e.response?.data?.error || 'เกิดข้อผิดพลาด' } finally { adding.value = false } }
function formatDate(d) { if (!d) return '-'; return new Date(d).toLocaleDateString('th-TH', { year: 'numeric', month: 'short', day: 'numeric' }) }
function openChangePasswordModal(u) { changePasswordTarget.value = u; changePasswordForm.value = { newPassword: '', confirmPassword: '' }; changePasswordError.value = ''; changePasswordSuccess.value = '' }
async function confirmChangePassword() {
  changePasswordError.value = ''; changePasswordSuccess.value = ''
  if (changePasswordForm.value.newPassword !== changePasswordForm.value.confirmPassword) { changePasswordError.value = 'รหัสผ่านไม่ตรงกัน'; return }
  changePasswordLoading.value = true
  try { await adminApi.updateUserPassword(changePasswordTarget.value.id, changePasswordForm.value.newPassword); changePasswordSuccess.value = 'เปลี่ยนรหัสผ่านสำเร็จ'; setTimeout(() => { changePasswordTarget.value = null }, 1500) } catch (e) { changePasswordError.value = e.response?.data?.error || 'เกิดข้อผิดพลาด' } finally { changePasswordLoading.value = false }
}
</script>

<style scoped>
.admin-page {
  min-height: 100vh;
  padding: 32px 16px;
  display: flex;
  justify-content: center;
  background: var(--bg-base);
}

.admin-card {
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-xl);
  padding: 36px;
  width: 100%;
  max-width: 960px;
  box-shadow: var(--shadow-card);
}

.admin-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 28px;
  flex-wrap: wrap;
  gap: 16px;
}

h1 {
  font-size: 1.6rem;
  font-weight: 800;
  color: var(--text-primary);
}

.admin-header p {
  color: var(--text-muted);
  margin-top: 4px;
  font-size: 0.88rem;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.btn-add {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 10px 20px;
  background: linear-gradient(135deg, var(--gold), var(--gold-light));
  color: #1a0e00;
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  font-weight: 700;
  font-family: inherit;
  transition: opacity 0.2s;
}

.btn-add:hover {
  opacity: 0.88;
}

.btn-back-nav {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 10px 20px;
  background: rgba(201, 168, 76, 0.06);
  color: var(--gold);
  border: 1px solid var(--gold-border);
  border-radius: 30px;
  text-decoration: none;
  font-weight: 600;
  font-size: 0.88rem;
  transition: all 0.2s;
}

.btn-back-nav:hover {
  background: var(--gold-dim);
  border-color: var(--gold);
  box-shadow: 0 0 12px rgba(201, 168, 76, 0.2);
}

.loading {
  display: flex;
  align-items: center;
  gap: 10px;
  color: var(--text-muted);
  padding: 20px;
}

.spinner-gold {
  width: 18px;
  height: 18px;
  border: 2px solid var(--gold-border);
  border-top-color: var(--gold);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.alert-error {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.25);
  border-radius: var(--radius-sm);
  color: #fca5a5;
  font-size: 0.88rem;
  margin-bottom: 16px;
}

.alert-success {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: rgba(16, 185, 129, 0.1);
  border: 1px solid rgba(16, 185, 129, 0.25);
  border-radius: var(--radius-sm);
  color: #6ee7b7;
  font-size: 0.88rem;
  margin-bottom: 16px;
}

.tabs-header {
  display: flex;
  gap: 0;
  margin-bottom: 20px;
  border-bottom: 1px solid var(--border-subtle);
}

.tab-button {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 12px 24px;
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  cursor: pointer;
  font-weight: 600;
  font-size: 0.9rem;
  color: var(--text-muted);
  transition: all 0.2s;
  font-family: inherit;
}

.tab-button:hover {
  color: var(--gold-light);
}

.tab-button.active {
  color: var(--gold);
  border-bottom-color: var(--gold);
}

.table-wrapper {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.88rem;
}

thead tr {
  background: var(--bg-elevated);
}

th {
  text-align: left;
  padding: 12px 16px;
  color: var(--text-muted);
  font-size: 0.72rem;
  text-transform: uppercase;
  font-weight: 700;
  letter-spacing: 0.06em;
  border-bottom: 1px solid var(--border-subtle);
  white-space: nowrap;
}

td {
  padding: 14px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
  vertical-align: middle;
  color: var(--text-secondary);
}

tr:hover td {
  background: rgba(255, 255, 255, 0.02);
}

tr:last-child td {
  border-bottom: none;
}

.td-id {
  color: var(--text-muted);
  font-size: 0.82rem;
}

.td-name {
  display: flex;
  align-items: center;
  gap: 10px;
  font-weight: 600;
  color: var(--text-primary);
}

.td-email {
  color: var(--text-muted);
  font-size: 0.85rem;
}

.td-date {
  color: var(--text-muted);
  font-size: 0.82rem;
  white-space: nowrap;
}

.td-actions {
  display: flex;
  gap: 6px;
  align-items: center;
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--gold), var(--gold-light));
  color: #1a0e00;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.82rem;
  font-weight: 800;
  flex-shrink: 0;
}

.role-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 10px;
  border-radius: 20px;
  font-size: 0.78rem;
  font-weight: 600;
}

.badge-admin {
  background: rgba(201, 168, 76, 0.15);
  color: var(--gold-light);
  border: 1px solid var(--gold-border);
}

.badge-user {
  background: rgba(139, 92, 246, 0.12);
  color: #c4b5fd;
  border: 1px solid rgba(139, 92, 246, 0.25);
}

.btn-action {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: none;
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-sm);
  padding: 6px 9px;
  cursor: pointer;
  color: var(--text-muted);
  transition: all 0.2s;
}

.btn-role:hover {
  background: var(--gold-dim);
  border-color: var(--gold-border);
  color: var(--gold-light);
}

.btn-pw:hover {
  background: rgba(139, 92, 246, 0.12);
  border-color: rgba(139, 92, 246, 0.3);
  color: #c4b5fd;
}

.btn-del:hover {
  background: rgba(239, 68, 68, 0.1);
  border-color: rgba(239, 68, 68, 0.3);
  color: #fca5a5;
}

.self-label {
  font-size: 0.72rem;
  color: var(--text-muted);
  font-style: italic;
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
  padding: 32px;
  width: 100%;
  max-width: 440px;
  box-shadow: var(--shadow-deep), var(--shadow-gold);
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(14px)
  }

  to {
    opacity: 1;
    transform: translateY(0)
  }
}

.modal-warn-icon {
  color: var(--orange);
  margin-bottom: 10px;
}

.modal h2 {
  font-size: 1.2rem;
  color: var(--text-primary);
  font-weight: 700;
  margin-bottom: 12px;
}

.modal p {
  color: var(--text-secondary);
  font-size: 0.9rem;
  margin-bottom: 8px;
}

.modal-sub {
  color: var(--text-muted);
  margin-bottom: 16px;
  font-size: 0.88rem;
}

.warning-text {
  color: #fca5a5;
  font-size: 0.85rem;
  margin-bottom: 16px;
}

.form-group {
  margin-bottom: 14px;
}

label {
  display: block;
  font-weight: 600;
  margin-bottom: 6px;
  color: var(--text-muted);
  font-size: 0.8rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

input,
select {
  width: 100%;
  padding: 10px 14px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-sm);
  font-size: 0.92rem;
  font-family: inherit;
  color: var(--text-primary);
  outline: none;
  transition: border-color 0.2s;
}

input::placeholder {
  color: var(--text-muted);
}

input:focus,
select:focus {
  border-color: var(--gold);
  box-shadow: 0 0 0 3px rgba(201, 168, 76, 0.1);
}

select option {
  background: var(--bg-elevated);
}

.modal-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  margin-top: 20px;
}

.btn-cancel {
  padding: 10px 20px;
  background: rgba(255, 255, 255, 0.05);
  color: var(--text-secondary);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-weight: 600;
  font-family: inherit;
}

.btn-cancel:hover {
  border-color: var(--border);
  color: var(--text-primary);
}

.btn-confirm {
  padding: 10px 20px;
  background: linear-gradient(135deg, var(--gold), var(--gold-light));
  color: #1a0e00;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-weight: 700;
  font-family: inherit;
}

.btn-confirm:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-danger {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 10px 20px;
  background: rgba(239, 68, 68, 0.8);
  color: white;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-weight: 700;
  font-family: inherit;
}

.btn-danger:hover {
  background: rgba(239, 68, 68, 1);
}
</style>
