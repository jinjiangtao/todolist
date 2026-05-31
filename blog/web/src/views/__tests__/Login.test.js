import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import Login from '../Login.vue'

vi.mock('@/utils/request', () => ({
  default: {
    post: vi.fn().mockResolvedValue({ data: { token: 'test-token' } })
  }
}))

describe('Login Component', () => {
  let router
  let pinia

  beforeEach(async () => {
    pinia = createPinia()
    setActivePinia(pinia)
    
    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', name: 'Home', component: { template: '<div>Home</div>' } },
        { path: '/login', name: 'Login', component: Login },
        { path: '/admin', name: 'Admin', component: { template: '<div>Admin</div>' } }
      ]
    })
    router.push('/login')
    await router.isReady()
  })

  it('renders login form', () => {
    const wrapper = mount(Login, {
      global: {
        plugins: [router, pinia]
      }
    })
    expect(wrapper.find('h1').text()).toBe('登录')
  })

  it('has username and password fields', () => {
    const wrapper = mount(Login, {
      global: {
        plugins: [router, pinia]
      }
    })
    expect(wrapper.vm.form.username).toBe('')
    expect(wrapper.vm.form.password).toBe('')
  })

  it('validates required fields', () => {
    const wrapper = mount(Login, {
      global: {
        plugins: [router, pinia]
      }
    })
    expect(wrapper.vm.rules.username).toBeDefined()
    expect(wrapper.vm.rules.password).toBeDefined()
    expect(wrapper.vm.rules.username[0].required).toBe(true)
    expect(wrapper.vm.rules.password[0].required).toBe(true)
  })
})