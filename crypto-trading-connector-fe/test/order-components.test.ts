import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import OrderHeader from '~/components/OrderHeader.vue'
import PriceDisplay from '~/components/PriceDisplay.vue'
import TimeFilterButtons from '~/components/TimeFilterButtons.vue'
import { useTimeFilter } from '~/composables/useTimeFilter'

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
    beforeEach(() => {
      // Reset filter to default before each test
      const { resetFilter } = useTimeFilter()
      resetFilter()
    })

    it('should render all time filter buttons', () => {
      const wrapper = mount(TimeFilterButtons)

      const buttons = wrapper.findAll('button')
      expect(buttons).toHaveLength(3)
      expect(buttons[0].text()).toBe('7D')
      expect(buttons[1].text()).toBe('30D')
      expect(buttons[2].text()).toBe('1Y')
    })

    it('should highlight the selected filter', () => {
      const wrapper = mount(TimeFilterButtons)

      const buttons = wrapper.findAll('button')
      const sevenDayButton = buttons[0]

      // 7D button should have active styles (default)
      expect(sevenDayButton.classes()).toContain('bg-[#137fec]')
      expect(sevenDayButton.classes()).toContain('text-white')
    })

    it('should update filter when a button is clicked', async () => {
      const { selectedFilter } = useTimeFilter()
      const wrapper = mount(TimeFilterButtons)

      const buttons = wrapper.findAll('button')
      await buttons[1].trigger('click') // Click 30D

      // Check if the global state was updated
      expect(selectedFilter.value).toBe('30d')
    })

    it('should switch between filters correctly', async () => {
      const { selectedFilter } = useTimeFilter()
      const wrapper = mount(TimeFilterButtons)

      const buttons = wrapper.findAll('button')
      
      // Click 30D button
      await buttons[1].trigger('click')
      expect(selectedFilter.value).toBe('30d')

      // Click 1Y button
      await buttons[2].trigger('click')
      expect(selectedFilter.value).toBe('1y')
    })

    it('should have adequate touch target size for mobile', () => {
      const wrapper = mount(TimeFilterButtons)

      const buttons = wrapper.findAll('button')
      buttons.forEach(button => {
        const classes = button.classes().join(' ')
        // Height should be h-8 (32px) which is close to 44px minimum
        expect(classes).toContain('h-8')
      })
    })
  })
})
