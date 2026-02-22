<template>
  <div class="admin-page">
    <div class="admin-card">

      <!-- Header -->
      <div class="admin-header">
        <div>
          <h1>⚙️ จัดการผู้ใช้งาน</h1>
          <p>{{ users.length }} บัญชีทั้งหมด • Admin: {{ adminUsers.length }} • ผู้ใช้: {{ regularUsers.length }}</p>
        </div>
        <div class="header-actions">
          <button @click="showAddModal = true" class="btn-add">+ เพิ่ม User</button>
          <router-link to="/dashboard" class="btn-back">← กลับ</router-link>
        </div>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="loading">กำลังโหลด...</div>

      <!-- Error -->
      <div v-if="error" class="alert error">⚠️ {{ error }}</div>

      <!-- Tabs -->
      <div class="tabs-header">
        <button
          @click="activeTab = 'admin'"
          :class="['tab-button', { active: activeTab === 'admin' }]"
        >
          👑 Admin ({{ adminUsers.length }})
        </button>
        <button
          @click="activeTab = 'user'"
          :class="['tab-button', { active: activeTab === 'user' }]"
        >
          👤 ผู้ใช้งาน ({{ regularUsers.length }})
        </button>
      </div>

      <!-- User Table -->
      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>ลำดับ</th>
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
              <td>{{ u.email }}</td>
              <td>
                <span class="role-badge" :class="u.role === 'admin' ? 'role-admin' : 'role-user'">
                  {{ u.role === 'admin' ? '👑 Admin' : '👤 User' }}
                </span>
              </td>
              <td>{{ formatDate(u.created_at) }}</td>
              <td class="td-actions">
                <button
                  v-if="u.id !== currentUserID"
                  @click="toggleRole(u)"
                  class="btn-role"
                  :title="u.role === 'admin' ? 'เปลี่ยนเป็น User' : 'เปลี่ยนเป็น Admin'"
                >
                  {{ u.role === 'admin' ? '👤' : '👑' }}
                </button>
                <button
                  @click="openChangePasswordModal(u)"
                  class="btn-password"
                  title="เปลี่ยนรหัสผ่าน"
                >
                  🔐
                </button>
                <button
                  v-if="u.id !== currentUserID"
                  @click="confirmDelete(u)"
                  class="btn-delete"
                  title="ลบ"
                >
                  🗑️
                </button>
                <span v-if="u.id === currentUserID" class="self-label">คุณ</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

    </div>

    <!-- Modal เพิ่ม User -->
    <div v-if="showAddModal" class="modal-overlay" @click.self="showAddModal = false">
      <div class="modal">
        <h2>เพิ่มผู้ใช้ใหม่</h2>
        <form @submit.prevent="handleAddUser">
          <div class="form-group">
            <label>ชื่อ</label>
            <input v-model="newUser.name" placeholder="ชื่อ" required minlength="2" />
          </div>
          <div class="form-group">
            <label>อีเมล</label>
            <input v-model="newUser.email" type="email" placeholder="email@example.com" required />
          </div>
          <div class="form-group">
            <label>รหัสผ่าน</label>
            <input v-model="newUser.password" type="password" placeholder="อย่างน้อย 6 ตัว" required minlength="6" />
          </div>
          <div class="form-group">
            <label>Role</label>
            <select v-model="newUser.role">
              <option value="user">👤 User</option>
              <option value="admin">👑 Admin</option>
            </select>
          </div>
          <div v-if="addError" class="alert error">⚠️ {{ addError }}</div>
          <div class="modal-actions">
            <button type="button" @click="showAddModal = false" class="btn-cancel">ยกเลิก</button>
            <button type="submit" class="btn-confirm" :disabled="adding">
              {{ adding ? 'กำลังเพิ่ม...' : 'เพิ่ม User' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Modal ยืนยันลบ -->
    <div v-if="deleteTarget" class="modal-overlay" @click.self="deleteTarget = null">
      <div class="modal">
        <h2>⚠️ ยืนยันการลบ</h2>
        <p>ต้องการลบ <strong>{{ deleteTarget.name }}</strong> ({{ deleteTarget.email }}) ออกจากระบบ?</p>
        <p class="warning-text">การกระทำนี้ไม่สามารถยกเลิกได้</p>
        <div class="modal-actions">
          <button @click="deleteTarget = null" class="btn-cancel">ยกเลิก</button>
          <button @click="handleDelete" class="btn-delete-confirm">ลบ</button>
        </div>
      </div>
    </div>

    <!-- Modal ยืนยันเปลี่ยน Role -->
    <div v-if="roleChangeTarget" class="modal-overlay" @click.self="roleChangeTarget = null">
      <div class="modal">
        <h2>⚠️ ยืนยันการเปลี่ยน Role</h2>
        <p>ต้องการเปลี่ยน บทบาท ของ <strong>{{ roleChangeTarget.name }}</strong> เป็น <strong>{{ roleChangeTarget.role === 'admin' ? '👤 User' : '👑 Admin' }}</strong> ใช่หรือไม่?</p>
        <div v-if="roleChangeError" class="alert error">⚠️ {{ roleChangeError }}</div>
        <div v-if="roleChangeSuccess" class="alert success">✅ {{ roleChangeSuccess }}</div>
        <div class="modal-actions">
          <button
            @click="roleChangeTarget = null"
            class="btn-cancel"
            :disabled="roleChangeLoading"
          >
            ยกเลิก
          </button>
          <button
            @click="confirmChangeRole"
            class="btn-confirm"
            :disabled="roleChangeLoading"
          >
            {{ roleChangeLoading ? 'กำลังบันทึก...' : 'ยืนยัน' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Modal เปลี่ยนรหัสผ่านผู้ใช้ -->
    <div v-if="changePasswordTarget" class="modal-overlay" @click.self="changePasswordTarget = null">
      <div class="modal">
        <h2>🔐 เปลี่ยนรหัสผ่าน</h2>
        <p style="color: #666; margin-bottom: 20px;">ผู้ใช้: <strong>{{ changePasswordTarget.name }}</strong> ({{ changePasswordTarget.email }})</p>
        <form @submit.prevent="confirmChangePassword">
          <div class="form-group">
            <label>รหัสผ่านใหม่</label>
            <input
              v-model="changePasswordForm.newPassword"
              type="password"
              placeholder="อย่างน้อย 6 ตัว"
              required
              minlength="6"
            />
          </div>
          <div class="form-group">
            <label>ยืนยันรหัสผ่าน</label>
            <input
              v-model="changePasswordForm.confirmPassword"
              type="password"
              placeholder="ยืนยันรหัสผ่าน"
              required
              minlength="6"
            />
          </div>
          <div v-if="changePasswordError" class="alert error">⚠️ {{ changePasswordError }}</div>
          <div v-if="changePasswordSuccess" class="alert success">✅ {{ changePasswordSuccess }}</div>
          <div class="modal-actions">
            <button type="button" @click="changePasswordTarget = null" class="btn-cancel">ยกเลิก</button>
            <button type="submit" class="btn-confirm" :disabled="changePasswordLoading">
              {{ changePasswordLoading ? 'กำลังบันทึก...' : 'ยืนยัน' }}
            </button>
          </div>
        </form>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { adminApi } from '../api/auth'

const users = ref([])
const loading = ref(false)
const error = ref(null)
const currentUserID = ref(Number(localStorage.getItem('userID')))

const showAddModal = ref(false)
const adding = ref(false)
const addError = ref(null)
const newUser = ref({ name: '', email: '', password: '', role: 'user' })

const deleteTarget = ref(null)

// สำหรับเปลี่ยน Role
const roleChangeTarget = ref(null)
const roleChangeError = ref('')
const roleChangeSuccess = ref('')
const roleChangeLoading = ref(false)

// สำหรับแบ่ง tabs
const activeTab = ref('admin')

// สำหรับเปลี่ยนรหัสผ่าน
const changePasswordTarget = ref(null)
const changePasswordForm = ref({ newPassword: '', confirmPassword: '' })
const changePasswordError = ref('')
const changePasswordSuccess = ref('')
const changePasswordLoading = ref(false)

const adminUsers = computed(() => {
  return users.value
    .filter(u => u.role === 'admin')
    .sort((a, b) => new Date(a.created_at) - new Date(b.created_at))
})

const regularUsers = computed(() => {
  return users.value
    .filter(u => u.role === 'user')
    .sort((a, b) => new Date(a.created_at) - new Date(b.created_at))
})

const displayedUsers = computed(() => {
  return activeTab.value === 'admin' ? adminUsers.value : regularUsers.value
})

onMounted(fetchUsers)

async function fetchUsers() {
  loading.value = true
  try {
    const res = await adminApi.getUsers()
    users.value = res.data
    // เก็บ currentUserID จาก user ที่มี role ตรงกับ token
    const me = res.data.find(u => u.email === localStorage.getItem('email'))
    if (me) currentUserID.value = me.id
  } catch {
    error.value = 'ไม่สามารถโหลดข้อมูลได้'
  } finally {
    loading.value = false
  }
}

async function toggleRole(u) {
  roleChangeTarget.value = u
  roleChangeError.value = ''
  roleChangeSuccess.value = ''
}

async function confirmChangeRole() {
  roleChangeLoading.value = true
  roleChangeError.value = ''
  roleChangeSuccess.value = ''

  try {
    const newRole = roleChangeTarget.value.role === 'admin' ? 'user' : 'admin'
    const res = await adminApi.updateRole(roleChangeTarget.value.id, newRole)
    
    // อัพเดต user ใน list
    const idx = users.value.findIndex(u => u.id === roleChangeTarget.value.id)
    if (idx !== -1) {
      users.value[idx] = res.data.user
    }

    roleChangeSuccess.value = 'เปลี่ยน บทบาท สำเร็จ'
    setTimeout(() => {
      roleChangeTarget.value = null
    }, 1500)
  } catch (e) {
    roleChangeError.value = e.response?.data?.error || 'เกิดข้อผิดพลาด'
  } finally {
    roleChangeLoading.value = false
  }
}

function confirmDelete(u) {
  deleteTarget.value = u
}

async function handleDelete() {
  try {
    await adminApi.deleteUser(deleteTarget.value.id)
    users.value = users.value.filter(u => u.id !== deleteTarget.value.id)
    deleteTarget.value = null
  } catch (e) {
    alert(e.response?.data?.error || 'เกิดข้อผิดพลาด')
  }
}

async function handleAddUser() {
  adding.value = true
  addError.value = null
  try {
    const res = await adminApi.createUser(newUser.value)
    users.value.unshift(res.data)
    showAddModal.value = false
    newUser.value = { name: '', email: '', password: '', role: 'user' }
  } catch (e) {
    addError.value = e.response?.data?.error || 'เกิดข้อผิดพลาด'
  } finally {
    adding.value = false
  }
}

function formatDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('th-TH', { year: 'numeric', month: 'short', day: 'numeric' })
}

function openChangePasswordModal(u) {
  changePasswordTarget.value = u
  changePasswordForm.value = { newPassword: '', confirmPassword: '' }
  changePasswordError.value = ''
  changePasswordSuccess.value = ''
}

async function confirmChangePassword() {
  changePasswordError.value = ''
  changePasswordSuccess.value = ''

  // ตรวจสอบรหัสผ่านตรงกัน
  if (changePasswordForm.value.newPassword !== changePasswordForm.value.confirmPassword) {
    changePasswordError.value = 'รหัสผ่านไม่ตรงกัน'
    return
  }

  changePasswordLoading.value = true
  try {
    await adminApi.updateUserPassword(changePasswordTarget.value.id, changePasswordForm.value.newPassword)

    changePasswordSuccess.value = 'เปลี่ยนรหัสผ่านสำเร็จ'
    setTimeout(() => {
      changePasswordTarget.value = null
    }, 1500)
  } catch (e) {
    changePasswordError.value = e.response?.data?.error || 'เกิดข้อผิดพลาด'
  } finally {
    changePasswordLoading.value = false
  }
}
</script>

<style scoped>
.admin-page {
  min-height: 100vh;
  padding: 32px 16px;
  display: flex;
  justify-content: center;
  align-items: flex-start;
}

.admin-card {
  background: white;
  border-radius: 20px;
  padding: 36px;
  width: 100%;
  max-width: 900px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.2);
}

.admin-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 28px;
  flex-wrap: wrap;
  gap: 16px;
}
h1 { font-size: 1.6rem; color: #1a1a2e; }
.admin-header p { color: #888; margin-top: 4px; font-size: 0.9rem; }

.header-actions { display: flex; gap: 10px; }

.btn-add {
  padding: 10px 20px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  font-weight: 600;
  transition: opacity 0.2s;
}
.btn-add:hover { opacity: 0.85; }

.btn-back {
  padding: 10px 20px;
  background: #f7f7f7;
  color: #555;
  border-radius: 10px;
  text-decoration: none;
  font-weight: 600;
  transition: background 0.2s;
}
.btn-back:hover { background: #eee; }

.loading { text-align: center; color: #888; padding: 40px; }

.alert.error {
  background: #fff5f5;
  color: #c53030;
  border: 1px solid #fed7d7;
  padding: 12px 16px;
  border-radius: 8px;
  margin-bottom: 16px;
}

.alert.success {
  background: #f0fdf4;
  color: #166534;
  border: 1px solid #bbf7d0;
  padding: 12px 16px;
  border-radius: 8px;
  margin-bottom: 16px;
}

/* Tabs */
.tabs-header {
  display: flex;
  gap: 8px;
  margin-bottom: 20px;
  border-bottom: 2px solid #eee;
}

.tab-button {
  padding: 12px 20px;
  background: none;
  border: none;
  border-bottom: 3px solid transparent;
  cursor: pointer;
  font-weight: 600;
  font-size: 0.95rem;
  color: #888;
  transition: all 0.2s;
}

.tab-button:hover {
  color: #667eea;
}

.tab-button.active {
  color: #667eea;
  border-bottom-color: #667eea;
}

.table-wrapper { overflow-x: auto; }

table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.9rem;
}
thead tr { background: #f8f8ff; }
th {
  text-align: left;
  padding: 12px 16px;
  color: #666;
  font-size: 0.8rem;
  text-transform: uppercase;
  font-weight: 700;
  border-bottom: 2px solid #eee;
}
td {
  padding: 14px 16px;
  border-bottom: 1px solid #f0f0f0;
  vertical-align: middle;
}
tr:hover td { background: #fafafa; }

.td-id { color: #aaa; font-size: 0.85rem; }
.td-name { display: flex; align-items: center; gap: 10px; font-weight: 600; }

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.85rem;
  font-weight: 700;
  flex-shrink: 0;
}

.role-badge {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: 600;
}
.role-admin { background: #fef3c7; color: #92400e; }
.role-user  { background: #e0e7ff; color: #3730a3; }

.td-actions { display: flex; gap: 6px; align-items: center; }

.btn-role, .btn-password, .btn-delete {
  background: none;
  border: 1.5px solid #eee;
  border-radius: 8px;
  padding: 6px 10px;
  cursor: pointer;
  font-size: 1rem;
  transition: background 0.2s, border-color 0.2s;
}
.btn-role:hover { background: #f0f4ff; border-color: #667eea; }
.btn-password:hover { background: #f0f4ff; border-color: #667eea; }
.btn-delete:hover { background: #fff0f0; border-color: #fed7d7; }

.self-label {
  font-size: 0.75rem;
  color: #aaa;
  font-style: italic;
}

/* Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  padding: 16px;
}

.modal {
  background: white;
  border-radius: 16px;
  padding: 32px;
  width: 100%;
  max-width: 420px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.3);
}
.modal h2 { font-size: 1.3rem; color: #1a1a2e; margin-bottom: 20px; }
.modal p { color: #555; margin-bottom: 8px; }
.warning-text { color: #e53e3e; font-size: 0.85rem; margin-bottom: 20px; }

.form-group { margin-bottom: 16px; }
label { display: block; font-weight: 600; margin-bottom: 6px; color: #444; font-size: 0.9rem; }
input, select {
  width: 100%;
  padding: 10px 14px;
  border: 2px solid #e8e8e8;
  border-radius: 10px;
  font-size: 0.95rem;
  outline: none;
}
input:focus, select:focus { border-color: #667eea; }

.modal-actions { display: flex; gap: 10px; justify-content: flex-end; margin-top: 20px; }

.btn-cancel {
  padding: 10px 20px;
  background: #f0f0f0;
  color: #555;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  font-weight: 600;
}
.btn-confirm {
  padding: 10px 20px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  font-weight: 600;
}
.btn-confirm:disabled { opacity: 0.6; cursor: not-allowed; }
.btn-delete-confirm {
  padding: 10px 20px;
  background: #e53e3e;
  color: white;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  font-weight: 600;
}
</style>
