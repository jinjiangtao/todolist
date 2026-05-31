import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import Home from '../Home.vue'

vi.mock('@/utils/request', () => ({
  default: {
    get: vi.fn().mockResolvedValue({ 
      data: { 
        list: [], 
        total: 0, 
        page: 1, 
        page_size: 10, 
        total_page: 0 
      } 
    })
  }
}))

describe('Home Component', () => {
  let router
  let pinia

  beforeEach(async () => {
    pinia = createPinia()
    setActivePinia(pinia)
    
    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', name: 'Home', component: Home },
        { path: '/article/:id', name: 'ArticleDetail', component: { template: '<div>Article</div>' } },
        { path: '/login', name: 'Login', component: { template: '<div>Login</div>' } }
      ]
    })
    router.push('/')
    await router.isReady()
  })

  it('renders header with title', () => {
    const wrapper = mount(Home, {
      global: {
        plugins: [router, pinia]
      }
    })
    expect(wrapper.find('.header h1').text()).toBe('个人博客')
  })

  it('has login link', () => {
    const wrapper = mount(Home, {
      global: {
        plugins: [router, pinia]
      }
    })
    expect(wrapper.find('.login-link').exists()).toBe(true)
    expect(wrapper.find('.login-link').text()).toBe('登录')
  })

  it('has pagination component', () => {
    const wrapper = mount(Home, {
      global: {
        plugins: [router, pinia]
      }
    })
    expect(wrapper.vm.page).toBe(1)
    expect(wrapper.vm.pageSize).toBe(10)
  })

  it('formatDate function works correctly', () => {
    const wrapper = mount(Home, {
      global: {
        plugins: [router, pinia]
      }
    })
    const dateStr = '2024-01-15T10:00:00Z'
    const result = wrapper.vm.formatDate(dateStr)
    expect(result).toBeDefined()
  })

  it('getSummary function truncates long content', () => {
    const wrapper = mount(Home, {
      global: {
        plugins: [router, pinia]
      }
    })
    const longContent = 'This is a very long content that should be truncated because it exceeds 150 characters. We need to make sure that the summary function works correctly and adds ellipsis at the end of the truncated text.'
    const result = wrapper.vm.getSummary(longContent)
    expect(result.length).toBeLessThanOrEqual(153)
    expect(result.endsWith('...')).toBe(true)
  })

  it('getSummary function handles short content', () => {
    const wrapper = mount(Home, {
      global: {
        plugins: [router, pinia]
      }
    })
    const shortContent = 'Short content'
    const result = wrapper.vm.getSummary(shortContent)
    expect(result).toBe(shortContent)
    expect(result.endsWith('...')).toBe(false)
  })

  it('getSummary function handles empty content', () => {
    const wrapper = mount(Home, {
      global: {
        plugins: [router, pinia]
      }
    })
    const result = wrapper.vm.getSummary('')
    expect(result).toBe('')
  })
})