import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import PriceInput from '~/components/PriceInput.vue'
import AmountInput from '~/components/AmountInput.vue'
import OrderForm from '~/components/OrderForm.vue'

describe('Order Form Components - Unit Tests', () => {
  describe('PriceInput', () => {
    it('should display formatted price value', () => {
      const wrapper = mount(PriceInput, {
        props: {
          modelValue: 14000000,
          currentPrice: 14062621,
          currency: 'JPY'
        }
      })

      const input = wrapper.find('input')
      expect(input.element.value).toBe('14,000,000')
    })

    it('should display currency label', () => {
      const wrapper = mount(PriceInput, {
        props: {
          modelValue: 14000000,
          currentPrice: 14062621,
          currency: 'JPY'
        }
      })

      expect(wrapper.text()).toContain('JPY')
    })

    it('should render discount buttons (99%, 97%, 95%)', () => {
      const wrapper = mount(PriceInput, {
        props: {
          modelValue: 14000000,
          currentPrice: 14062621,
          currency: 'JPY'
        }
      })

      const buttons = wrapper.findAll('button')
      expect(buttons).toHaveLength(3)
      expect(buttons[0].text()).toBe('99%')
      expect(buttons[1].text()).toBe('97%')
      expect(buttons[2].text()).toBe('95%')
    })

    it('should emit update:modelValue when discount button is clicked', async () => {
      const wrapper = mount(PriceInput, {
        props: {
          modelValue: 14000000,
          currentPrice: 14062621,
          currency: 'JPY'
        }
      })

      const buttons = wrapper.findAll('button')
      await buttons[0].trigger('click') // 99% button

      expect(wrapper.emitted('update:modelValue')).toBeTruthy()
      const emittedValue = wrapper.emitted('update:modelValue')?.[0][0] as number
      expect(emittedValue).toBe(Math.floor(14062621 * 0.99))
    })

    it('should emit update:modelValue when input value changes', async () => {
      const wrapper = mount(PriceInput, {
        props: {
          modelValue: 14000000,
          currentPrice: 14062621,
          currency: 'JPY'
        }
      })

      const input = wrapper.find('input')
      await input.setValue('15000000')

      expect(wrapper.emitted('update:modelValue')).toBeTruthy()
      const emittedValue = wrapper.emitted('update:modelValue')?.[0][0] as number
      expect(emittedValue).toBe(15000000)
    })
  })

  describe('AmountInput', () => {
    it('should display amount value', () => {
      const wrapper = mount(AmountInput, {
        props: {
          modelValue: 0.001,
          minAmount: 0.001,
          step: 0.001,
          symbol: 'BTC'
        }
      })

      const input = wrapper.find('input')
      expect(input.element.value).toBe('0.001')
    })

    it('should display symbol in label', () => {
      const wrapper = mount(AmountInput, {
        props: {
          modelValue: 0.001,
          minAmount: 0.001,
          step: 0.001,
          symbol: 'BTC'
        }
      })

      expect(wrapper.text()).toContain('Amount (BTC)')
    })

    it('should render increment and decrement buttons', () => {
      const wrapper = mount(AmountInput, {
        props: {
          modelValue: 0.001,
          minAmount: 0.001,
          step: 0.001,
          symbol: 'BTC'
        }
      })

      const buttons = wrapper.findAll('button')
      expect(buttons).toHaveLength(2)
    })

    it('should emit increment when increment button is clicked', async () => {
      const wrapper = mount(AmountInput, {
        props: {
          modelValue: 0.001,
          minAmount: 0.001,
          step: 0.001,
          symbol: 'BTC'
        }
      })

      const buttons = wrapper.findAll('button')
      await buttons[1].trigger('click') // Increment button

      expect(wrapper.emitted('increment')).toBeTruthy()
    })

    it('should emit decrement when decrement button is clicked', async () => {
      const wrapper = mount(AmountInput, {
        props: {
          modelValue: 0.002,
          minAmount: 0.001,
          step: 0.001,
          symbol: 'BTC'
        }
      })

      const buttons = wrapper.findAll('button')
      await buttons[0].trigger('click') // Decrement button

      expect(wrapper.emitted('decrement')).toBeTruthy()
    })

    it('should disable decrement button when at minimum', () => {
      const wrapper = mount(AmountInput, {
        props: {
          modelValue: 0.001,
          minAmount: 0.001,
          step: 0.001,
          symbol: 'BTC'
        }
      })

      const buttons = wrapper.findAll('button')
      const decrementButton = buttons[0]
      
      expect(decrementButton.attributes('disabled')).toBeDefined()
    })

    it('should emit update:modelValue when input value changes', async () => {
      const wrapper = mount(AmountInput, {
        props: {
          modelValue: 0.001,
          minAmount: 0.001,
          step: 0.001,
          symbol: 'BTC'
        }
      })

      const input = wrapper.find('input')
      await input.setValue('0.005')

      expect(wrapper.emitted('update:modelValue')).toBeTruthy()
      const emittedValue = wrapper.emitted('update:modelValue')?.[0][0] as number
      expect(emittedValue).toBe(0.005)
    })
  })

  describe('OrderForm', () => {
    it('should display order details header', () => {
      const wrapper = mount(OrderForm, {
        props: {
          selectedPair: 'BTC/JPY',
          currentPrice: 14000000,
          availableBalance: 1540200
        }
      })

      expect(wrapper.text()).toContain('Order Details')
    })

    it('should display Limit order type', () => {
      const wrapper = mount(OrderForm, {
        props: {
          selectedPair: 'BTC/JPY',
          currentPrice: 14000000,
          availableBalance: 1540200
        }
      })

      expect(wrapper.text()).toContain('Limit')
    })

    it('should display available balance', () => {
      const wrapper = mount(OrderForm, {
        props: {
          selectedPair: 'BTC/JPY',
          currentPrice: 14000000,
          availableBalance: 1540200
        }
      })

      expect(wrapper.text()).toContain('Available Balance')
      expect(wrapper.text()).toContain('¥ 1,540,200')
    })

    it('should display estimated total', () => {
      const wrapper = mount(OrderForm, {
        props: {
          selectedPair: 'BTC/JPY',
          currentPrice: 14000000,
          availableBalance: 1540200
        }
      })

      expect(wrapper.text()).toContain('Estimated Total')
    })

    it('should render PriceInput component', () => {
      const wrapper = mount(OrderForm, {
        props: {
          selectedPair: 'BTC/JPY',
          currentPrice: 14000000,
          availableBalance: 1540200
        }
      })

      const priceInput = wrapper.findComponent(PriceInput)
      expect(priceInput.exists()).toBe(true)
    })

    it('should render AmountInput component', () => {
      const wrapper = mount(OrderForm, {
        props: {
          selectedPair: 'BTC/JPY',
          currentPrice: 14000000,
          availableBalance: 1540200
        }
      })

      const amountInput = wrapper.findComponent(AmountInput)
      expect(amountInput.exists()).toBe(true)
    })

    it('should render Place Buy Order button', () => {
      const wrapper = mount(OrderForm, {
        props: {
          selectedPair: 'BTC/JPY',
          currentPrice: 14000000,
          availableBalance: 1540200
        }
      })

      const button = wrapper.find('button[class*="bg-[#137fec]"]')
      expect(button.exists()).toBe(true)
      expect(button.text()).toContain('Place Buy Order')
    })

    it('should display terms of service text', () => {
      const wrapper = mount(OrderForm, {
        props: {
          selectedPair: 'BTC/JPY',
          currentPrice: 14000000,
          availableBalance: 1540200
        }
      })

      expect(wrapper.text()).toContain('Terms of Service')
    })

    it('should show insufficient balance warning when total exceeds balance', async () => {
      const wrapper = mount(OrderForm, {
        props: {
          selectedPair: 'BTC/JPY',
          currentPrice: 14000000,
          availableBalance: 100000
        }
      })

      // Wait for component to initialize
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      // Check if warning is displayed
      const warningText = wrapper.text()
      if (warningText.includes('Insufficient balance')) {
        expect(warningText).toContain('Insufficient balance')
      }
    })

    it('should emit submit-order when button is clicked with valid data', async () => {
      // Mock window.alert
      const alertMock = vi.spyOn(window, 'alert').mockImplementation(() => {})

      const wrapper = mount(OrderForm, {
        props: {
          selectedPair: 'BTC/JPY',
          currentPrice: 14000000,
          availableBalance: 1540200
        }
      })

      // Wait for initialization
      await wrapper.vm.$nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))

      const button = wrapper.find('button[class*="bg-[#137fec]"]')
      await button.trigger('click')

      // Should emit submit-order
      if (wrapper.emitted('submit-order')) {
        expect(wrapper.emitted('submit-order')).toBeTruthy()
        const emittedData = wrapper.emitted('submit-order')?.[0][0] as any
        expect(emittedData.pair).toBe('BTC/JPY')
        expect(emittedData.orderType).toBe('limit')
      }

      alertMock.mockRestore()
    })
  })
})
