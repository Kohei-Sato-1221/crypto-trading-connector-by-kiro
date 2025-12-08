import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import MarketHeader from './MarketHeader.vue'

// Mock the composable
vi.mock('~/composables/useTimeFilter', () => ({
  useTimeFilter: () => ({
    selectedFilter: { value: '24h' },
    setFilter: vi.fn(),
    isSelected: (filter: string) => filter === '24h'
  })
}))

describe('MarketHeader - Unit Tests', () => {
  it('should display market status indicator', () => {
    const wrapper = mount(MarketHeader)
    
    const statusText = wrapper.text()
    expect(statusText).toContain('Market is Open')
    expect(statusText).toContain('Real-time Data')
  })

  it('should display 5 time filter buttons', () => {
    const wrapper = mount(MarketHeader)
    
    const buttons = wrapper.findAll('button')
    expect(buttons.length).toBe(5)
  })

  it('should display correct time filter labels', () => {
    const wrapper = mount(MarketHeader)
    
    const text = wrapper.text()
    expect(text).toContain('24h')
    expect(text).toContain('7d')
    expect(text).toContain('30d')
    expect(text).toContain('1y')
    expect(text).toContain('All')
  })

  it('should have 24h button selected by default', () => {
    const wrapper = mount(MarketHeader)
    
    const buttons = wrapper.findAll('button')
    const activeButton = buttons.find(btn => 
      btn.classes().includes('bg-[#137fec]')
    )
    
    expect(activeButton).toBeDefined()
    expect(activeButton?.text()).toBe('24h')
  })

  it('should display app title', () => {
    const wrapper = mount(MarketHeader)
    
    const title = wrapper.find('h1')
    expect(title.text()).toBe('Crypto Trading Connector')
  })

  it('should have sticky positioning', () => {
    const wrapper = mount(MarketHeader)
    
    const header = wrapper.find('header')
    expect(header.classes()).toContain('sticky')
    expect(header.classes()).toContain('top-0')
  })

  it('should have backdrop blur effect', () => {
    const wrapper = mount(MarketHeader)
    
    const header = wrapper.find('header')
    expect(header.classes()).toContain('backdrop-blur-md')
  })

  it('should have minimum touch target size for buttons', () => {
    const wrapper = mount(MarketHeader)
    
    const buttons = wrapper.findAll('button')
    buttons.forEach(button => {
      expect(button.classes()).toContain('min-w-[44px]')
    })
  })
})
