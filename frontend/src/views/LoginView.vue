<template>
  <div class="auth-page">
    <div class="orb orb-1"></div>
    <div class="orb orb-2"></div>

    <div class="auth-card">
      <div class="auth-header">
        <div class="brand">
          <Wine :size="16" />
          <span class="brand-name">NOIR</span>
        </div>
        <h1>ยินดีต้อนรับกลับ</h1>
        <p>เข้าสู่ประสบการณ์พิเศษของคุณ</p>
      </div>

      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label>อีเมล</label>
          <div class="input-wrapper">
            <Mail :size="15" class="input-icon" />
            <input v-model="email" type="email" placeholder="example@email.com" required :disabled="loading" />
          </div>
        </div>

        <div class="form-group">
          <label>รหัสผ่าน</label>
          <div class="input-wrapper">
            <Lock :size="15" class="input-icon" />
            <input v-model="password" :type="showPassword ? 'text' : 'password'" placeholder="รหัสผ่านของคุณ" required :disabled="loading" />
            <button type="button" class="toggle-pw" @click="showPassword = !showPassword">
              <EyeOff v-if="showPassword" :size="16" />
              <Eye v-else :size="16" />
            </button>
          </div>
        </div>

        <div v-if="error" class="alert-error">
          <AlertTriangle :size="15" />
          {{ error }}
        </div>

        <button type="submit" class="btn-primary" :disabled="loading">
          <span v-if="loading" class="spinner"></span>
          <LogIn v-else :size="16" />
          {{ loading ? 'กำลังเข้าสู่ระบบ...' : 'เข้าสู่ระบบ' }}
        </button>
      </form>

      <div class="divider"><span>หรือ</span></div>

      <p class="auth-footer">
        ยังไม่มีบัญชี?
        <router-link to="/register">สมัครสมาชิก <ArrowRight :size="12" /></router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import { storeToRefs } from 'pinia'
import { Wine, Mail, Lock, Eye, EyeOff, AlertTriangle, LogIn, ArrowRight } from 'lucide-vue-next'

const store = useAuthStore()
const { loading, error } = storeToRefs(store)
const email = ref('')
const password = ref('')
const showPassword = ref(false)

async function handleLogin() {
  await store.login(email.value, password.value)
}
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  background:
    radial-gradient(ellipse at 20% 50%, rgba(139,92,246,0.12) 0%, transparent 60%),
    radial-gradient(ellipse at 80% 20%, rgba(201,168,76,0.10) 0%, transparent 50%),
    var(--bg-base);
  position: relative;
  overflow: hidden;
}
.orb { position:fixed; border-radius:50%; filter:blur(80px); pointer-events:none; z-index:0; animation:floatOrb 8s ease-in-out infinite; }
.orb-1 { width:350px; height:350px; background:radial-gradient(circle,rgba(201,168,76,0.18),transparent 70%); top:-80px; left:-80px; }
.orb-2 { width:400px; height:400px; background:radial-gradient(circle,rgba(139,92,246,0.15),transparent 70%); bottom:-100px; right:-100px; animation-delay:3s; animation-direction:reverse; }
@keyframes floatOrb { 0%,100%{transform:translate(0,0) scale(1)} 50%{transform:translate(20px,-20px) scale(1.05)} }

.auth-card {
  position: relative; z-index: 1;
  background: var(--bg-glass);
  backdrop-filter: blur(24px); -webkit-backdrop-filter: blur(24px);
  border: 1px solid var(--border);
  border-radius: var(--radius-xl);
  padding: 44px 40px;
  width: 100%; max-width: 420px;
  box-shadow: var(--shadow-deep), var(--shadow-gold);
  animation: fadeIn 0.5s ease both;
}
@keyframes fadeIn { from{opacity:0;transform:translateY(20px)} to{opacity:1;transform:translateY(0)} }

.auth-header { text-align:center; margin-bottom:36px; }

.brand {
  display: inline-flex; align-items:center; gap:7px;
  margin-bottom: 20px;
  padding: 7px 18px;
  background: var(--gold-dim);
  border: 1px solid var(--gold-border);
  border-radius: 30px;
  color: var(--gold-light);
}
.brand-name { font-size:0.9rem; font-weight:800; letter-spacing:0.22em; }

h1 { font-size:1.7rem; font-weight:700; color:var(--text-primary); margin-bottom:8px; }
.auth-header p { color:var(--text-secondary); font-size:0.9rem; }

.form-group { margin-bottom:20px; }
label { display:block; font-weight:600; margin-bottom:8px; color:var(--text-secondary); font-size:0.78rem; letter-spacing:0.06em; text-transform:uppercase; }

.input-wrapper { position:relative; display:flex; align-items:center; }
.input-icon { position:absolute; left:14px; color:var(--text-muted); pointer-events:none; }

input {
  width:100%; padding:13px 16px 13px 42px;
  background: rgba(255,255,255,0.04);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  font-size:0.95rem; font-family:inherit;
  color: var(--text-primary);
  transition: border-color 0.25s, box-shadow 0.25s;
  outline: none;
}
input::placeholder { color:var(--text-muted); }
input:focus { border-color:var(--gold); box-shadow:0 0 0 3px rgba(201,168,76,0.12); }
input:disabled { opacity:0.5; cursor:not-allowed; }

.input-wrapper input { padding-right:48px; }
.toggle-pw {
  position:absolute; right:12px;
  background:none; border:none; cursor:pointer;
  color:var(--text-muted); display:flex; align-items:center;
  transition:color 0.2s;
}
.toggle-pw:hover { color:var(--gold); }

.alert-error {
  display:flex; gap:8px; align-items:center;
  padding:12px 16px;
  background:rgba(239,68,68,0.1); border:1px solid rgba(239,68,68,0.3);
  border-radius:var(--radius-sm);
  color:#fca5a5; font-size:0.88rem; margin-bottom:16px;
}

.btn-primary {
  width:100%; padding:14px;
  background: linear-gradient(135deg, var(--gold), var(--gold-light));
  color: #1a0e00;
  border:none; border-radius:var(--radius-md);
  font-size:0.95rem; font-weight:700; font-family:inherit;
  cursor:pointer;
  display:flex; align-items:center; justify-content:center; gap:8px;
  transition:opacity 0.2s, transform 0.15s, box-shadow 0.25s;
}
.btn-primary:hover:not(:disabled) { opacity:0.9; transform:translateY(-1px); box-shadow:0 6px 24px rgba(201,168,76,0.4); }
.btn-primary:disabled { opacity:0.5; cursor:not-allowed; }

.spinner { width:16px; height:16px; border:2px solid rgba(0,0,0,0.3); border-top-color:#1a0e00; border-radius:50%; animation:spin 0.8s linear infinite; }
@keyframes spin { to { transform:rotate(360deg); } }

.divider { display:flex; align-items:center; gap:12px; margin:24px 0; color:var(--text-muted); font-size:0.8rem; }
.divider::before, .divider::after { content:''; flex:1; height:1px; background:var(--border-subtle); }

.auth-footer { text-align:center; color:var(--text-muted); font-size:0.9rem; }
.auth-footer a { color:var(--gold); font-weight:600; text-decoration:none; display:inline-flex; align-items:center; gap:3px; transition:color 0.2s; }
.auth-footer a:hover { color:var(--gold-light); }

@media (max-width:480px) { .auth-card { padding:32px 24px; } }
</style>
