import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/store/auth'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue')
  },
  {
    path: '/article/:id',
    name: 'ArticleDetail',
    component: () => import('@/views/ArticleDetail.vue')
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue')
  },
  {
    path: '/admin',
    name: 'Admin',
    component: () => import('@/views/admin/Admin.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/admin/article/edit/:id?',
    name: 'ArticleEdit',
    component: () => import('@/views/admin/ArticleEdit.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/admin/settings',
    name: 'Settings',
    component: () => import('@/views/admin/Settings.vue'),
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else if (to.path === '/login' && authStore.isAuthenticated) {
    next('/admin')
  } else {
    next()
  }
})

export default router
