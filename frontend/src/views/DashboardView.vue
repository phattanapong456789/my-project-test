<template>
  <div class="dashboard-page">
    <div class="bg-orb bg-orb-1"></div>
    <div class="bg-orb bg-orb-2"></div>

    <div class="dashboard-container">
      <!-- Navbar -->
      <nav class="navbar">
        <div class="navbar-brand">
          <Wine :size="18" />
          <span class="brand-text">NOIR</span>
        </div>
        <div class="navbar-actions">
          <router-link v-if="user?.role === 'admin'" to="/admin/reservations" class="nav-btn nav-btn-outline">
            <ClipboardList :size="14" /> จัดการจอง
          </router-link>
          <router-link v-if="user?.role === 'admin'" to="/admin" class="nav-btn nav-btn-outline">
            <Users :size="14" /> จัดการผู้ใช้
          </router-link>
          <router-link v-if="user?.role === 'admin'" to="/admin/floor-plan" class="nav-btn nav-btn-outline">
            <LayoutGrid :size="14" /> ผังร้าน
          </router-link>
          <router-link to="/booking" class="nav-btn nav-btn-gold">
            <GlassWater :size="14" /> จองโต๊ะ
          </router-link>
          <router-link to="/my-reservations" class="nav-btn nav-btn-outline">
            <CalendarDays :size="14" /> การจองของฉัน
          </router-link>
          <button @click="store.logout()" class="nav-btn nav-btn-danger">
            <LogOut :size="14" /> ออกจากระบบ
          </button>
        </div>
      </nav>

      <!-- Hero -->
      <div class="hero-section">
        <div class="avatar" :class="user?.role === 'admin' ? 'avatar-admin' : 'avatar-user'">
          {{ initials }}
        </div>
        <div class="hero-text">
          <h1>สวัสดี, {{ user?.name || '...' }}</h1>
          <p>{{ user?.email }}</p>
        </div>
        <span class="role-badge" :class="user?.role === 'admin' ? 'badge-admin' : 'badge-user'">
          <Crown v-if="user?.role === 'admin'" :size="12" />
          <UserIcon v-else :size="12" />
          {{ user?.role === 'admin' ? 'Admin' : 'Member' }}
        </span>
      </div>

      <!-- Info Cards -->
      <div class="info-grid">
        <div class="info-card">
          <div class="info-card-icon" style="color:var(--purple)"><Hash :size="24" /></div>
          <div>
            <div class="info-label">User ID</div>
            <div class="info-value">#{{ user?.id }}</div>
          </div>
          <div class="info-glow" style="background:var(--purple-dim)"></div>
        </div>
        <div class="info-card">
          <div class="info-card-icon" style="color:var(--cyan)"><Mail :size="24" /></div>
          <div>
            <div class="info-label">อีเมล</div>
            <div class="info-value email-val">{{ user?.email }}</div>
          </div>
          <div class="info-glow" style="background:rgba(6,182,212,0.08)"></div>
        </div>
        <div class="info-card">
          <div class="info-card-icon" style="color:var(--green)"><Calendar :size="24" /></div>
          <div>
            <div class="info-label">สมัครเมื่อ</div>
            <div class="info-value">{{ formatDate(user?.created_at) }}</div>
          </div>
          <div class="info-glow" style="background:rgba(16,185,129,0.08)"></div>
        </div>
        <div class="info-card">
          <div class="info-card-icon" style="color:var(--gold)"><Shield :size="24" /></div>
          <div>
            <div class="info-label">บทบาท</div>
            <div class="info-value">
              <span class="role-badge" :class="user?.role === 'admin' ? 'badge-admin' : 'badge-user'">
                <Crown v-if="user?.role === 'admin'" :size="11" />
                <UserIcon v-else :size="11" />
                {{ user?.role === 'admin' ? 'Admin' : 'Member' }}
              </span>
            </div>
          </div>
          <div class="info-glow" style="background:var(--gold-dim)"></div>
        </div>
      </div>

      <!-- Token Box -->
      <div class="token-box">
        <div class="token-header">
          <div class="token-title">
            <span class="token-dot"></span>
            <Key :size="14" /> JWT Token
          </div>
          <button @click="copyToken" class="btn-copy">
            <Check v-if="copied" :size="13" /> <span>{{ copied ? 'Copied!' : 'Copy' }}</span>
            <Copy v-if="!copied" :size="13" />
          </button>
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
import {
  Wine, GlassWater, ClipboardList, Users, LayoutGrid, CalendarDays, LogOut,
  Crown, User as UserIcon, Hash, Mail, Calendar, Shield, Key, Copy, Check
} from 'lucide-vue-next'

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
  return t.length > 80 ? t.slice(0, 80) + '...' : t
})
function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('th-TH', { year: 'numeric', month: 'long', day: 'numeric' })
}
async function copyToken() {
  const token = localStorage.getItem('token') || ''
  await navigator.clipboard.writeText(token)
  copied.value = true
  setTimeout(() => copied.value = false, 2000)
}
</script>

<style scoped>
.dashboard-page { min-height:100vh; background:radial-gradient(ellipse at 10% 20%,rgba(139,92,246,0.08) 0%,transparent 50%),radial-gradient(ellipse at 90% 80%,rgba(201,168,76,0.07) 0%,transparent 50%),var(--bg-base); position:relative; overflow-x:hidden; }
.bg-orb { position:fixed; border-radius:50%; filter:blur(100px); pointer-events:none; z-index:0; animation:floatOrb 12s ease-in-out infinite; }
.bg-orb-1 { width:500px; height:500px; background:radial-gradient(circle,rgba(201,168,76,0.12),transparent 70%); top:-200px; right:-100px; }
.bg-orb-2 { width:600px; height:600px; background:radial-gradient(circle,rgba(139,92,246,0.1),transparent 70%); bottom:-200px; left:-150px; animation-delay:5s; }
@keyframes floatOrb { 0%,100%{transform:scale(1) translate(0,0)} 50%{transform:scale(1.08) translate(15px,-15px)} }

.dashboard-container { position:relative; z-index:1; max-width:860px; margin:0 auto; padding:0 20px 40px; }

.navbar { display:flex; justify-content:space-between; align-items:center; padding:20px 0; margin-bottom:32px; border-bottom:1px solid var(--border-subtle); flex-wrap:wrap; gap:16px; }
.navbar-brand { display:flex; align-items:center; gap:8px; color:var(--gold); }
.brand-text { font-size:1rem; font-weight:800; letter-spacing:0.22em; color:var(--gold); }
.navbar-actions { display:flex; gap:8px; flex-wrap:wrap; }

.nav-btn { padding:8px 14px; border-radius:var(--radius-sm); font-size:0.82rem; font-weight:600; font-family:inherit; text-decoration:none; cursor:pointer; border:none; transition:all 0.2s; display:inline-flex; align-items:center; gap:5px; white-space:nowrap; }
.nav-btn-gold { background:linear-gradient(135deg,var(--gold),var(--gold-light)); color:#1a0e00; }
.nav-btn-gold:hover { opacity:0.88; transform:translateY(-1px); }
.nav-btn-outline { background:rgba(255,255,255,0.04); color:var(--text-secondary); border:1px solid var(--border-subtle); }
.nav-btn-outline:hover { border-color:var(--gold-border); color:var(--gold-light); background:var(--gold-dim); }
.nav-btn-danger { background:rgba(239,68,68,0.08); color:#fca5a5; border:1px solid rgba(239,68,68,0.2); }
.nav-btn-danger:hover { background:rgba(239,68,68,0.15); }

.hero-section { display:flex; align-items:center; gap:20px; margin-bottom:32px; flex-wrap:wrap; }
.avatar { width:64px; height:64px; border-radius:50%; font-size:1.3rem; font-weight:800; display:flex; align-items:center; justify-content:center; flex-shrink:0; border:2px solid var(--border); }
.avatar-admin { background:linear-gradient(135deg,var(--gold),var(--gold-light)); color:#1a0e00; box-shadow:0 0 24px rgba(201,168,76,0.35); }
.avatar-user { background:linear-gradient(135deg,var(--purple),#a78bfa); color:white; box-shadow:0 0 24px rgba(139,92,246,0.3); }
.hero-text { flex:1; }
.hero-text h1 { font-size:1.6rem; font-weight:700; color:var(--text-primary); margin-bottom:4px; }
.hero-text p { color:var(--text-muted); font-size:0.9rem; }

.role-badge { display:inline-flex; align-items:center; gap:4px; padding:5px 12px; border-radius:20px; font-size:0.78rem; font-weight:600; }
.badge-admin { background:var(--gold-dim); color:var(--gold-light); border:1px solid var(--gold-border); }
.badge-user  { background:var(--purple-dim); color:#c4b5fd; border:1px solid rgba(139,92,246,0.3); }

.info-grid { display:grid; grid-template-columns:repeat(2,1fr); gap:16px; margin-bottom:24px; }
.info-card { position:relative; padding:20px; background:var(--bg-surface); border:1px solid var(--border-subtle); border-radius:var(--radius-md); display:flex; align-items:center; gap:16px; overflow:hidden; transition:border-color 0.25s,transform 0.2s; }
.info-card:hover { border-color:var(--border); transform:translateY(-2px); }
.info-glow { position:absolute; inset:0; pointer-events:none; opacity:0; transition:opacity 0.3s; }
.info-card:hover .info-glow { opacity:1; }
.info-card-icon { flex-shrink:0; }
.info-label { font-size:0.7rem; color:var(--text-muted); text-transform:uppercase; letter-spacing:0.08em; font-weight:600; }
.info-value { font-size:0.95rem; font-weight:700; color:var(--text-primary); margin-top:3px; }
.email-val { font-size:0.85rem; word-break:break-all; }

.token-box { background:var(--bg-surface); border:1px solid var(--border-subtle); border-radius:var(--radius-md); padding:20px 24px; }
.token-header { display:flex; justify-content:space-between; align-items:center; margin-bottom:12px; }
.token-title { display:flex; align-items:center; gap:8px; color:var(--text-secondary); font-size:0.85rem; font-weight:600; }
.token-dot { width:8px; height:8px; border-radius:50%; background:var(--green); box-shadow:0 0 8px rgba(16,185,129,0.6); animation:pulse 2s ease-in-out infinite; }
@keyframes pulse { 0%,100%{opacity:1} 50%{opacity:0.4} }
.btn-copy { padding:5px 12px; background:rgba(255,255,255,0.06); color:var(--text-secondary); border:1px solid var(--border-subtle); border-radius:var(--radius-sm); cursor:pointer; font-size:0.8rem; font-family:inherit; font-weight:600; transition:all 0.2s; display:inline-flex; align-items:center; gap:5px; }
.btn-copy:hover { background:var(--gold-dim); color:var(--gold-light); border-color:var(--gold-border); }
.token-text { color:#6ee7b7; font-size:0.78rem; word-break:break-all; line-height:1.6; font-family:'Courier New',monospace; opacity:0.9; }

@media (max-width:600px) { .info-grid{grid-template-columns:1fr} .nav-btn{padding:7px 10px;font-size:0.78rem} .hero-text h1{font-size:1.3rem} }
</style>
