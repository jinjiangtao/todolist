import { describe, it, expect, beforeEach } from 'vitest'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import { useAuthStore } from '@/store/auth'

const routes = [
  { path: '/', name: 'Home', component: { template: '<div>Home</div>' } },
  { path: '/article/:id', name: 'ArticleDetail', component: { template: '<div>Article</div>' } },
  { path: '/login', name: 'Login', component: { template: '<div>Login</div>' } },
  { path: '/admin', name: 'Admin', component: { template: '<div>Admin</div>' }, meta: { requiresAuth: true } },
  { path: '/admin/article/edit/:id?', name: 'ArticleEdit', component: { template: '<div>Edit</div>' }, meta: { requiresAuth: true } },
  { path: '/admin/settings', name: 'Settings', component: { template: '<div>Settings</div>' }, meta: { requiresAuth: true } }
]

describe('Router', () => {
  let router
  let pinia

  beforeEach(async () => {
    pinia = createPinia()
    setActivePinia(pinia)
    
    router = createRouter({
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
  })

  it('allows access to home page without authentication', async () => {
    router.push('/')
    await router.isReady()
    expect(router.currentRoute.value.name).toBe('Home')
  })

  it('allows access to article detail page without authentication', async () => {
    router.push('/article/1')
    await router.isReady()
    expect(router.currentRoute.value.name).toBe('ArticleDetail')
  })

  it('allows access to login page without authentication', async () => {
    router.push('/login')
    await router.isReady()
    expect(router.currentRoute.value.name).toBe('Login')
  })

  it('redirects to login when accessing admin without authentication', async () => {
    const authStore = useAuthStore()
    authStore.logout()
    
    router.push('/admin')
    await router.isReady()
    await router.currentRoute.value
    expect(router.currentRoute.value.path).toBe('/login')
  })

  it('redirects to login when accessing admin/article/edit without authentication', async () => {
    const authStore = useAuthStore()
    authStore.logout()
    
    router.push('/admin/article/edit')
    await router.isReady()
    expect(router.currentRoute.value.path).toBe('/login')
  })

  it('redirects to login when accessing admin/settings without authentication', async () => {
    const authStore = useAuthStore()
    authStore.logout()
    
    router.push('/admin/settings')
    await router.isReady()
    expect(router.currentRoute.value.path).toBe('/login')
  })

  it('allows access to admin when authenticated', async () => {
    const authStore = useAuthStore()
    authStore.setToken('test-token')
    
    router.push('/admin')
    await router.isReady()
    expect(router.currentRoute.value.name).toBe('Admin')
  })

  it('redirects to admin when accessing login while authenticated', async () => {
    const authStore = useAuthStore()
    authStore.setToken('test-token')
    
    router.push('/login')
    await router.isReady()
    expect(router.currentRoute.value.path).toBe('/admin')
  })

  it('allows access to admin/article/edit when authenticated', async () => {
    const authStore = useAuthStore()
    authStore.setToken('test-token')
    
    router.push('/admin/article/edit')
    await router.isReady()
    expect(router.currentRoute.value.name).toBe('ArticleEdit')
  })

  it('allows access to admin/article/edit/:id when authenticated', async () => {
    const authStore = useAuthStore()
    authStore.setToken('test-token')
    
    router.push('/admin/article/edit/1')
    await router.isReady()
    expect(router.currentRoute.value.name).toBe('ArticleEdit')
    expect(router.currentRoute.value.params.id).toBe('1')
  })

  it('allows access to admin/settings when authenticated', async () => {
    const authStore = useAuthStore()
    authStore.setToken('test-token')
    
    router.push('/admin/settings')
    await router.isReady()
    expect(router.currentRoute.value.name).toBe('Settings')
  })
})