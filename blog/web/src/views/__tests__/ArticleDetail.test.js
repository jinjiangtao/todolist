import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import ArticleDetail from '../ArticleDetail.vue'

vi.mock('@/utils/request', () => ({
  default: {
    get: vi.fn().mockResolvedValue({ 
      data: {
        id: 1,
        title: 'Test Article',
        content: '# Hello World',
        created_at: '2024-01-15T10:00:00Z',
        view_count: 10
      }
    })
  }
}))

describe('ArticleDetail Component', () => {
  let router
  let pinia

  beforeEach(async () => {
    pinia = createPinia()
    setActivePinia(pinia)
    
    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', name: 'Home', component: { template: '<div>Home</div>' } },
        { path: '/article/:id', name: 'ArticleDetail', component: ArticleDetail }
      ]
    })
    router.push('/article/1')
    await router.isReady()
  })

  it('renders header with title', () => {
    const wrapper = mount(ArticleDetail, {
      global: {
        plugins: [router, pinia]
      }
    })
    expect(wrapper.find('.header h1').text()).toBe('个人博客')
  })

  it('has back link', () => {
    const wrapper = mount(ArticleDetail, {
      global: {
        plugins: [router, pinia]
      }
    })
    expect(wrapper.find('.back-link').exists()).toBe(true)
    expect(wrapper.find('.back-link').text()).toBe('返回首页')
  })

  it('formatDate function works correctly', () => {
    const wrapper = mount(ArticleDetail, {
      global: {
        plugins: [router, pinia]
      }
    })
    const dateStr = '2024-01-15T10:00:00Z'
    const result = wrapper.vm.formatDate(dateStr)
    expect(result).toBeDefined()
  })

  it('initializes with loading state', () => {
    const wrapper = mount(ArticleDetail, {
      global: {
        plugins: [router, pinia]
      }
    })
    expect(wrapper.vm.loading).toBe(true)
    expect(wrapper.vm.article).toBe(null)
  })

  it('renderedContent is empty when article is null', () => {
    const wrapper = mount(ArticleDetail, {
      global: {
        plugins: [router, pinia]
      }
    })
    expect(wrapper.vm.renderedContent).toBe('')
  })

  it('renderedContent renders markdown when article exists', async () => {
    const wrapper = mount(ArticleDetail, {
      global: {
        plugins: [router, pinia]
      }
    })
    wrapper.vm.article = {
      id: 1,
      title: 'Test Article',
      content: '# Hello World',
      created_at: '2024-01-15T10:00:00Z',
      view_count: 10
    }
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.renderedContent).toContain('Hello World')
  })
})