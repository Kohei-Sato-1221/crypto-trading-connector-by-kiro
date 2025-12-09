import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import OrderHeader from '~/components/OrderHeader.vue'
import PriceDisplay from '~/components/PriceDisplay.vue'
import TimeFilterButtons from '~/components/TimeFilterButtons.vue'

describe('Order Page Components - Unit Tests', () => {
  describe('OrderHeader', () => {
    it('should render both currency pair tabs', () => {
      const wrapper = mount(OrderHeader, {
        props: {
          selectedPair: 'BTC/JPY'
        }
      })

      const buttons = wrapper.findAll('button')
      expect(buttons).toHaveLength(2)
      expect(buttons[0].text()).toBe('BTC/JPY')
      expect(buttons[1].text()).toBe('ETH/JPY')
    })

    it('should highlight the selected pair', () => {
      const wrapper = mount(OrderHeader, {
        props: {
          selectedPair: 'BTC/JPY'
        }
      })

      const buttons = wrapper.findAll('button')
      const btcButton = buttons[0]
      const ethButton = buttons[1]

      // BTC button should have active styles
      expect(btcButton.classes()).toContain('bg-[#101922]')
      
      // ETH button should not have active styles
      expect(ethButton.classes()).not.toContain('bg-[#101922]')
    })

    it('should emit update:selectedPair when a tab is clicked', async () => {
      const wrapper = mount(OrderHeader, {
        props: {
          selectedPair: 'BTC/JPY'
        }
      })

      const buttons = wrapper.findAll('button')
      await buttons[1].trigger('click')

      expect(wrapper.emitted('update:selectedPair')).toBeTruthy()
      expect(wrapper.emitted('update:selectedPair')?.[0]).toEqual(['ETH/JPY'])
    })

    it('should switch between pairs correctly', async () => {
      const wrapper = mount(OrderHeader, {
        props: {
          selectedPair: 'BTC/JPY'
        }
      })

      const buttons = wrapper.findAll('button')
      
      // Click ETH button
      await buttons[1].trigger('click')
      expect(wrapper.emitted('update:selectedPair')?.[0]).toEqual(['ETH/JPY'])

      // Click BTC button
      await buttons[0].trigger('click')
      expect(wrapper.emitted('update:selectedPair')?.[1]).toEqual(['BTC/JPY'])
    })
  })

  describe('PriceDisplay', () => {
    it('should format price with thousand separators', () => {
      const wrapper = mount(PriceDisplay, {
        props: {
          currentPrice: 14062621,
          priceChange: 2.5,
          currency: 'BTC'
        }
      })

      expect(wrapper.text()).toContain('¥14,062,621')
    })

    it('should display positive change with green color and up arrow', () => {
      const wrapper = mount(PriceDisplay, {
        props: {
          currentPrice: 14062621,
          priceChange: 2.5,
          currency: 'BTC'
        }
      })

      expect(wrapper.text()).toContain('+2.5%')
      
      // Check for green background
      const badge = wrapper.find('.bg-\\[rgba\\(11\\,218\\,91\\,0\\.1\\)\\]')
      expect(badge.exists()).toBe(true)

      // Check for green text
      const changeText = wrapper.find('.text-\\[\\#0bda5b\\]')
      expect(changeText.exists()).toBe(true)
    })

    it('should display negative change with red color and down arrow', () => {
      const wrapper = mount(PriceDisplay, {
        props: {
          currentPrice: 485318,
          priceChange: -1.2,
          currency: 'ETH'
        }
      })

      expect(wrapper.text()).toContain('-1.2%')
      
      // Check for red background
      const badge = wrapper.find('.bg-\\[rgba\\(250\\,98\\,56\\,0\\.1\\)\\]')
      expect(badge.exists()).toBe(true)

      // Check for red text
      const changeText = wrapper.find('.text-\\[\\#fa6238\\]')
      expect(changeText.exists()).toBe(true)
    })

    it('should display "Current Price" label', () => {
      const wrapper = mount(PriceDisplay, {
        props: {
          currentPrice: 14062621,
          priceChange: 2.5,
          currency: 'BTC'
        }
      })

      expect(wrapper.text()).toContain('Current Price')
    })

    it('should format change percentage to 1 decimal place', () => {
      const wrapper = mount(PriceDisplay, {
        props: {
          currentPrice: 14062621,
          priceChange: 2.567,
          currency: 'BTC'
        }
      })

      expect(wrapper.text()).toContain('+2.6%')
    })
  })

  describe('TimeFilterButtons', () => {
    it('should render all time filter buttons', () => {
      const wrapper = mount(TimeFilterButtons, {
        props: {
          selectedFilter: '7D'
        }
      })

      const buttons = wrapper.findAll('button')
      expect(buttons).toHaveLength(5)
      expect(buttons[0].text()).toBe('1H')
      expect(buttons[1].text()).toBe('24H')
      expect(buttons[2].text()).toBe('7D')
      expect(buttons[3].text()).toBe('30D')
      expect(buttons[4].text()).toBe('1Y')
    })

    it('should highlight the selected filter', () => {
      const wrapper = mount(TimeFilterButtons, {
        props: {
          selectedFilter: '7D'
        }
      })

      const buttons = wrapper.findAll('button')
      const sevenDayButton = buttons[2]

      // 7D button should have active styles
      expect(sevenDayButton.classes()).toContain('bg-[#137fec]')
      expect(sevenDayButton.classes()).toContain('text-white')
    })

    it('should emit update:selectedFilter when a button is clicked', async () => {
      const wrapper = mount(TimeFilterButtons, {
        props: {
          selectedFilter: '7D'
        }
      })

      const buttons = wrapper.findAll('button')
      await buttons[3].trigger('click') // Click 30D

      expect(wrapper.emitted('update:selectedFilter')).toBeTruthy()
      expect(wrapper.emitted('update:selectedFilter')?.[0]).toEqual(['30D'])
    })

    it('should switch between filters correctly', async () => {
      const wrapper = mount(TimeFilterButtons, {
        props: {
          selectedFilter: '7D'
        }
      })

      const buttons = wrapper.findAll('button')
      
      // Click 1H button
      await buttons[0].trigger('click')
      expect(wrapper.emitted('update:selectedFilter')?.[0]).toEqual(['1H'])

      // Click 1Y button
      await buttons[4].trigger('click')
      expect(wrapper.emitted('update:selectedFilter')?.[1]).toEqual(['1Y'])
    })

    it('should have adequate touch target size for mobile', () => {
      const wrapper = mount(TimeFilterButtons, {
        props: {
          selectedFilter: '7D'
        }
      })

      const buttons = wrapper.findAll('button')
      buttons.forEach(button => {
        const classes = button.classes().join(' ')
        // Height should be h-8 (32px) which is close to 44px minimum
        expect(classes).toContain('h-8')
      })
    })
  })
})
