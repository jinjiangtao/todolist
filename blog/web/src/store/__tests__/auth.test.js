import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '../auth'

describe('Auth Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it('initializes with empty token', () => {
    const store = useAuthStore()
    expect(store.token).toBe('')
    expect(store.isAuthenticated).toBe(false)
  })

  it('sets token and saves to localStorage', () => {
    const store = useAuthStore()
    const token = 'test-token-123'
    store.setToken(token)
    expect(store.token).toBe(token)
    expect(store.isAuthenticated).toBe(true)
    expect(localStorage.getItem('token')).toBe(token)
  })

  it('loads token from localStorage on initialization', () => {
    const token = 'persisted-token'
    localStorage.setItem('token', token)
    const store = useAuthStore()
    expect(store.token).toBe(token)
    expect(store.isAuthenticated).toBe(true)
  })

  it('logs out by clearing token and localStorage', () => {
    const store = useAuthStore()
    store.setToken('some-token')
    store.logout()
    expect(store.token).toBe('')
    expect(store.isAuthenticated).toBe(false)
    expect(localStorage.getItem('token')).toBeNull()
  })
})
