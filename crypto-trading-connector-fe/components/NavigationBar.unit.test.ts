import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import NavigationBar from './NavigationBar.vue'

// Mock Nuxt's useRoute composable
vi.mock('#app', () => ({
  useRoute: () => ({
    path: '/market'
  })
}))

describe('NavigationBar - Unit Tests', () => {
  it('should display 3 navigation links', () => {
    const wrapper = mount(NavigationBar, {
      global: {
        stubs: {
          NuxtLink: {
            template: '<a :to="to" :class="$attrs.class"><slot /></a>',
            props: ['to']
          }
        }
      }
    })

    const links = wrapper.findAll('a')
    expect(links.length).toBe(3)
    
    wrapper.unmount()
  })

  it('should display correct labels for all navigation items', () => {
    const wrapper = mount(NavigationBar, {
      global: {
        stubs: {
          NuxtLink: {
            template: '<a><slot /></a>',
            props: ['to']
          }
        }
      }
    })

    const text = wrapper.text()
    expect(text).toContain('Market')
    expect(text).toContain('Trade')
    expect(text).toContain('History')
    
    wrapper.unmount()
  })

  it('should have correct navigation paths', () => {
    const wrapper = mount(NavigationBar, {
      global: {
        stubs: {
          NuxtLink: true // Use default stub
        }
      }
    })

    // Find all NuxtLink components
    const links = wrapper.findAllComponents({ name: 'NuxtLink' })
    
    // Verify we have 3 links
    expect(links.length).toBe(3)
    
    // Verify the paths
    const paths = links.map(link => link.props('to'))
    expect(paths).toEqual(['/market', '/trade', '/history'])
    
    wrapper.unmount()
  })

  it('should highlight the active navigation item based on current route', () => {
    const wrapper = mount(NavigationBar, {
      global: {
        stubs: {
          NuxtLink: {
            template: '<a :class="$attrs.class"><slot /></a>',
            props: ['to']
          }
        }
      }
    })

    const links = wrapper.findAll('a')
    const activeLinks = links.filter(link => 
      link.classes().includes('text-[#137fec]')
    )
    
    // Exactly one link should be active (the /market link since that's our mock)
    expect(activeLinks.length).toBe(1)
    
    wrapper.unmount()
  })
})
