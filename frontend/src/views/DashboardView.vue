<template>
  <div class="dashboard-page">
    <div class="dashboard-card">
      <!-- Header -->
      <div class="dash-header">
        <div class="welcome">
          <div class="avatar" :class="user?.role === 'admin' ? 'avatar-admin' : ''">{{ initials }}</div>
          <div>
            <h1>สวัสดี, {{ user?.name || '...' }}! 👋</h1>
            <p>{{ user?.email }}</p>
          </div>
        </div>
        <div class="header-actions">
          <router-link v-if="user?.role === 'admin'" to="/admin/reservations" class="btn-admin">📋 จัดการจอง</router-link>
          <router-link v-if="user?.role === 'admin'" to="/admin" class="btn-users">👥 จัดการผู้ใช้</router-link>
          <router-link v-if="user?.role === 'admin'" to="/admin/floor-plan" class="btn-floor">🗺️ ผังร้าน</router-link>
          <router-link to="/booking" class="btn-booking">🍺 จองโต๊ะ</router-link>
          <router-link to="/my-reservations" class="btn-my-res">📋 การจองของฉัน</router-link>
          <button @click="store.logout()" class="btn-logout">ออกจากระบบ</button>
        </div>
      </div>

      <!-- Info Cards -->
      <div class="info-grid">
        <div class="info-card purple">
          <div class="info-icon">🆔</div>
          <div>
            <div class="info-label">User ID</div>
            <div class="info-value">#{{ user?.id }}</div>
          </div>
        </div>

        <div class="info-card blue">
          <div class="info-icon">📧</div>
          <div>
            <div class="info-label">อีเมล</div>
            <div class="info-value">{{ user?.email }}</div>
          </div>
        </div>

        <div class="info-card green">
          <div class="info-icon">📅</div>
          <div>
            <div class="info-label">สมัครเมื่อ</div>
            <div class="info-value">{{ formatDate(user?.created_at) }}</div>
          </div>
        </div>

        <div class="info-card orange">
          <div class="info-icon">🎭</div>
          <div>
            <div class="info-label">Role</div>
            <div class="info-value">
              <span class="role-badge" :class="user?.role === 'admin' ? 'role-admin' : 'role-user'">
                {{ user?.role === 'admin' ? '👑 Admin' : '👤 User' }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Token Info -->
      <div class="token-box">
        <div class="token-header">
          <span>🔑 JWT Token (ใช้สำหรับเรียก API)</span>
          <button @click="copyToken" class="btn-copy">{{ copied ? '✅ Copied!' : '📋 Copy' }}</button>
        </div>
        <code class="token-text">{{ shortToken }}</code>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { storeToRefs } from 'pinia'
import { authApi } from '../api/auth'

const store = useAuthStore()
const { user } = storeToRefs(store)
const copied = ref(false)

onMounted(() => store.fetchMe())

const initials = computed(() => {
  if (!user.value?.name) return '?'
  return user.value.name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2)
})

const shortToken = computed(() => {
  const t = localStorage.getItem('token') || ''
  return t.length > 60 ? t.slice(0, 60) + '...' : t
})

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('th-TH', {
    year: 'numeric', month: 'long', day: 'numeric'
  })
}

async function copyToken() {
  const token = localStorage.getItem('token') || ''
  await navigator.clipboard.writeText(token)
  copied.value = true
  setTimeout(() => copied.value = false, 2000)
}
</script>

<style scoped>
.dashboard-page {
  min-height: 100vh;
  padding: 32px 16px;
  display: flex;
  justify-content: center;
  align-items: flex-start;
}

.dashboard-card {
  background: white;
  border-radius: 20px;
  padding: 36px;
  width: 100%;
  max-width: 680px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.2);
}

.dash-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
  flex-wrap: wrap;
  gap: 16px;
}

.welcome {
  display: flex;
  align-items: center;
  gap: 16px;
}

.avatar {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  font-size: 1.2rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

h1 { font-size: 1.4rem; color: #1a1a2e; }
.welcome p { color: #888; font-size: 0.9rem; margin-top: 2px; }

.header-actions { display: flex; gap: 8px; align-items: center; flex-wrap: wrap; }

.btn-admin {
  padding: 10px 20px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-radius: 10px;
  text-decoration: none;
  font-weight: 600;
  font-size: 0.9rem;
  transition: opacity 0.2s;
}
.btn-admin:hover { opacity: 0.85; }

.btn-users {
  padding: 10px 20px;
  background: #f0fdf4;
  color: #16a34a;
  border: 1.5px solid #bbf7d0;
  border-radius: 10px;
  text-decoration: none;
  font-weight: 600;
  font-size: 0.9rem;
  transition: background 0.2s;
}
.btn-users:hover { background: #dcfce7; }

.btn-floor {
  padding: 10px 20px;
  background: #f0f4ff;
  color: #667eea;
  border: 1.5px solid #c7d2fe;
  border-radius: 10px;
  text-decoration: none;
  font-weight: 600;
  font-size: 0.9rem;
  transition: background 0.2s;
}
.btn-floor:hover { background: #e0e7ff; }

.btn-booking {
  padding: 10px 20px;
  background: linear-gradient(135deg, #f6ad55, #ed8936);
  color: white;
  border-radius: 10px;
  text-decoration: none;
  font-weight: 600;
  font-size: 0.9rem;
  transition: opacity 0.2s;
}
.btn-booking:hover { opacity: 0.85; }

.btn-my-res {
  padding: 10px 20px;
  background: #f0f4ff;
  color: #667eea;
  border: 1.5px solid #c7d2fe;
  border-radius: 10px;
  text-decoration: none;
  font-weight: 600;
  font-size: 0.9rem;
}
.btn-my-res:hover { background: #e0e7ff; }

.btn-logout {
  padding: 10px 20px;
  background: #fff0f0;
  color: #e53e3e;
  border: 1.5px solid #fed7d7;
  border-radius: 10px;
  cursor: pointer;
  font-weight: 600;
  transition: background 0.2s;
}
.btn-logout:hover { background: #fed7d7; }

.avatar-admin { background: linear-gradient(135deg, #f6ad55, #ed8936) !important; }

.role-badge {
  display: inline-block;
  padding: 3px 10px;
  border-radius: 20px;
  font-size: 0.85rem;
  font-weight: 600;
}
.role-admin { background: #fef3c7; color: #92400e; }
.role-user  { background: #e0e7ff; color: #3730a3; }

.info-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-bottom: 24px;
}

.info-card {
  padding: 20px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  gap: 14px;
}
.info-card.purple { background: #f3f0ff; }
.info-card.blue   { background: #ebf8ff; }
.info-card.green  { background: #f0fff4; }
.info-card.orange { background: #fffaf0; }

.info-icon { font-size: 1.8rem; }
.info-label { font-size: 0.75rem; color: #888; font-weight: 600; text-transform: uppercase; }
.info-value { font-size: 0.95rem; font-weight: 700; color: #1a1a2e; margin-top: 2px; word-break: break-all; }

.token-box {
  background: #1a1a2e;
  border-radius: 14px;
  padding: 20px;
}
.token-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #a0aec0;
  font-size: 0.85rem;
  margin-bottom: 10px;
}
.btn-copy {
  background: rgba(255,255,255,0.1);
  color: white;
  border: none;
  padding: 4px 12px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.8rem;
  transition: background 0.2s;
}
.btn-copy:hover { background: rgba(255,255,255,0.2); }
.token-text {
  color: #68d391;
  font-size: 0.8rem;
  word-break: break-all;
  line-height: 1.5;
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: white;
  border-radius: 14px;
  padding: 32px;
  width: 90%;
  max-width: 450px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.modal h2 {
  margin-bottom: 24px;
  color: #1a1a2e;
  font-size: 1.3rem;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-weight: 600;
  color: #1a1a2e;
  font-size: 0.9rem;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 10px;
  border: 1.5px solid #e2e8f0;
  border-radius: 8px;
  font-size: 0.95rem;
  transition: border-color 0.2s;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #667eea;
}

.alert {
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 16px;
  font-size: 0.9rem;
}

.alert.error {
  background: #fed7d7;
  color: #742a2a;
  border-left: 4px solid #fc8181;
}

.alert.success {
  background: #c6f6d5;
  color: #22543d;
  border-left: 4px solid #9ae6b4;
}

.modal-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 24px;
}

.btn-cancel,
.btn-confirm {
  padding: 10px 24px;
  border-radius: 8px;
  border: none;
  cursor: pointer;
  font-weight: 600;
  font-size: 0.9rem;
  transition: all 0.2s;
}

.btn-cancel {
  background: #f0f0f0;
  color: #1a1a2e;
}

.btn-cancel:hover {
  background: #e2e2e2;
}

.btn-confirm {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
}

.btn-confirm:hover:not(:disabled) {
  opacity: 0.85;
}

.btn-confirm:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

@media (max-width: 480px) {
  .info-grid { grid-template-columns: 1fr; }
}
</style>
