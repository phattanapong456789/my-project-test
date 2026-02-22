import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import DashboardView from '../views/DashboardView.vue'
import AdminView from '../views/AdminView.vue'
import BookingView from '../views/BookingView.vue'
import MyReservationsView from '../views/MyReservationsView.vue'
import AdminReservationsView from '../views/AdminReservationsView.vue'
import AdminFloorPlanView from '../views/AdminFloorPlanView.vue'

const routes = [
  { path: '/',                  redirect: '/login' },
  { path: '/login',             component: LoginView,             meta: { guest: true } },
  { path: '/register',          component: RegisterView,          meta: { guest: true } },
  { path: '/dashboard',         component: DashboardView,         meta: { requiresAuth: true } },
  { path: '/booking',           component: BookingView,           meta: { requiresAuth: true } },
  { path: '/my-reservations',   component: MyReservationsView,    meta: { requiresAuth: true } },
  { path: '/admin',             component: AdminView,             meta: { requiresAuth: true, adminOnly: true } },
  { path: '/admin/reservations',component: AdminReservationsView, meta: { requiresAuth: true, adminOnly: true } },
  { path: '/admin/floor-plan',  component: AdminFloorPlanView,    meta: { requiresAuth: true, adminOnly: true } },
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, _, next) => {
  const token = localStorage.getItem('token')
  const role = localStorage.getItem('role')
  if (to.meta.requiresAuth && !token) next('/login')
  else if (to.meta.guest && token) next('/dashboard')
  else if (to.meta.adminOnly && role !== 'admin') next('/dashboard')
  else next()
})

export default router
