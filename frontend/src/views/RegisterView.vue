<template>
  <div class="auth-page">
    <div class="auth-card">
      <div class="auth-header">
        <div class="logo">✨</div>
        <h1>สมัครสมาชิก</h1>
        <p>สร้างบัญชีใหม่ได้เลย</p>
      </div>

      <form @submit.prevent="handleRegister">
        <div class="form-group">
          <label>ชื่อ</label>
          <input
            v-model="name"
            type="text"
            placeholder="ชื่อของคุณ"
            required
            minlength="2"
            :disabled="loading"
          />
        </div>

        <div class="form-group">
          <label>อีเมล</label>
          <input
            v-model="email"
            type="email"
            placeholder="example@email.com"
            required
            :disabled="loading"
          />
        </div>

        <div class="form-group">
          <label>รหัสผ่าน</label>
          <div class="input-wrapper">
            <input
              v-model="password"
              :type="showPassword ? 'text' : 'password'"
              placeholder="อย่างน้อย 6 ตัวอักษร"
              required
              minlength="6"
              :disabled="loading"
            />
            <button type="button" class="toggle-pw" @click="showPassword = !showPassword">
              {{ showPassword ? '🙈' : '👁️' }}
            </button>
          </div>
        </div>

        <div class="form-group">
          <label>ยืนยันรหัสผ่าน</label>
          <input
            v-model="confirmPassword"
            :type="showPassword ? 'text' : 'password'"
            placeholder="พิมพ์รหัสผ่านอีกครั้ง"
            required
            :disabled="loading"
          />
          <p v-if="confirmPassword && password !== confirmPassword" class="field-error">
            รหัสผ่านไม่ตรงกัน
          </p>
        </div>

        <div v-if="error" class="alert error">⚠️ {{ error }}</div>

        <button
          type="submit"
          class="btn-primary"
          :disabled="loading || (confirmPassword && password !== confirmPassword)"
        >
          <span v-if="loading" class="spinner"></span>
          {{ loading ? 'กำลังสมัคร...' : 'สมัครสมาชิก' }}
        </button>
      </form>

      <p class="auth-footer">
        มีบัญชีแล้ว?
        <router-link to="/login">เข้าสู่ระบบ</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import { storeToRefs } from 'pinia'

const store = useAuthStore()
const { loading, error } = storeToRefs(store)

const name = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const showPassword = ref(false)

async function handleRegister() {
  if (password.value !== confirmPassword.value) return
  await store.register(name.value, email.value, password.value)
}
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
}

.auth-card {
  background: white;
  border-radius: 20px;
  padding: 40px;
  width: 100%;
  max-width: 420px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.2);
}

.auth-header {
  text-align: center;
  margin-bottom: 32px;
}

.logo { font-size: 3rem; margin-bottom: 12px; }
h1 { font-size: 1.8rem; color: #1a1a2e; margin-bottom: 6px; }
.auth-header p { color: #888; }

.form-group { margin-bottom: 20px; }

label {
  display: block;
  font-weight: 600;
  margin-bottom: 6px;
  color: #444;
  font-size: 0.9rem;
}

input {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid #e8e8e8;
  border-radius: 10px;
  font-size: 1rem;
  transition: border-color 0.2s;
  outline: none;
}
input:focus { border-color: #667eea; }
input:disabled { background: #f5f5f5; }

.input-wrapper { position: relative; }
.input-wrapper input { padding-right: 48px; }
.toggle-pw {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  cursor: pointer;
  font-size: 1.1rem;
}

.field-error {
  color: #c53030;
  font-size: 0.8rem;
  margin-top: 4px;
}

.alert {
  padding: 12px 16px;
  border-radius: 8px;
  margin-bottom: 16px;
  font-size: 0.9rem;
}
.alert.error {
  background: #fff5f5;
  color: #c53030;
  border: 1px solid #fed7d7;
}

.btn-primary {
  width: 100%;
  padding: 14px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: opacity 0.2s, transform 0.1s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}
.btn-primary:hover:not(:disabled) { opacity: 0.9; transform: translateY(-1px); }
.btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }

.auth-footer {
  text-align: center;
  margin-top: 24px;
  color: #888;
  font-size: 0.9rem;
}
.auth-footer a {
  color: #667eea;
  font-weight: 600;
  text-decoration: none;
}
.auth-footer a:hover { text-decoration: underline; }
</style>
