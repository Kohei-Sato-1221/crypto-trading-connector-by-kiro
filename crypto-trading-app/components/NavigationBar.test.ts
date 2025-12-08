import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import * as fc from 'fast-check'
import NavigationBar from './NavigationBar.vue'

// Feature: crypto-market-page, Property 5: Navigation active state
describe('NavigationBar - Property Based Tests', () => {
  it('Property 5: For any current route, exactly one navigation item should be highlighted', () => {
    // Define all valid navigation paths
    const validPaths = ['/market', '/trade', '/history', '/portfolio']
    
    fc.assert(
      fc.property(
        fc.constantFrom(...validPaths),
        (currentPath) => {
          // Mock useRoute to return the current path
          vi.doMock('#app', () => ({
            useRoute: () => ({
              path: currentPath
            })
          }))

          // Mount the component
          const wrapper = mount(NavigationBar, {
            global: {
              stubs: {
                NuxtLink: {
                  template: '<a :class="$attrs.class"><slot /></a>',
                  props: ['to']
                }
              },
              mocks: {
                useRoute: () => ({
                  path: currentPath
                })
              }
            }
          })

          // Find all navigation links
          const links = wrapper.findAll('a')
          
          // Count how many links have the active class (blue color)
          const activeLinks = links.filter(link => 
            link.classes().includes('text-[#137fec]')
          )
          
          // Exactly one link should be active
          expect(activeLinks.length).toBe(1)
          
          wrapper.unmount()
        }
      ),
      { numRuns: 100 } // Run 100 iterations as specified in design doc
    )
  })
})
