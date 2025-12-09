<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useOrderForm } from '~/composables/useOrderForm'
import PriceInput from './PriceInput.vue'
import AmountInput from './AmountInput.vue'

interface Props {
  selectedPair: 'BTC/JPY' | 'ETH/JPY'
  currentPrice: number
  availableBalance: number
  onSubmitOrder?: (order: any) => Promise<void>
}

interface Emits {
  (e: 'submit-order', order: {
    pair: string
    orderType: string
    price: number
    amount: number
    estimatedTotal: number
  }): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// State for submission
const isSubmitting = ref(false)
const submitError = ref<string | null>(null)
const submitSuccess = ref(false)

const pair = computed(() => props.selectedPair)
const currentPriceRef = ref(props.currentPrice)
const availableBalanceRef = ref(props.availableBalance)

// Watch for prop changes
watch(() => props.currentPrice, (newPrice) => {
  currentPriceRef.value = newPrice
})

watch(() => props.availableBalance, (newBalance) => {
  availableBalanceRef.value = newBalance
})

const {
  price,
  amount,
  estimatedTotal,
  minAmount,
  step,
  isBalanceSufficient,
  isValidOrder,
  initializeDefaults,
  incrementAmount,
  decrementAmount,
  setDiscountPrice,
  validateOrder
} = useOrderForm(pair, currentPriceRef, availableBalanceRef)

// Initialize on mount
onMounted(() => {
  initializeDefaults()
})

// Note: Pair change handling is done in useOrderForm composable

const symbol = computed(() => {
  return props.selectedPair === 'BTC/JPY' ? 'BTC' : 'ETH'
})

const formattedBalance = computed(() => {
  return `¥ ${props.availableBalance.toLocaleString()}`
})

const formattedTotal = computed(() => {
  return `¥ ${estimatedTotal.value.toLocaleString()}`
})

const handleSubmit = async () => {
  const validation = validateOrder()
  
  if (!validation.valid) {
    submitError.value = validation.error || 'Invalid order'
    return
  }

  const orderData = {
    pair: props.selectedPair,
    orderType: 'limit' as const,
    price: price.value,
    amount: amount.value,
    estimatedTotal: estimatedTotal.value
  }

  // Reset states
  submitError.value = null
  submitSuccess.value = false
  isSubmitting.value = true

  try {
    // If parent provides onSubmitOrder, use it (for API integration)
    if (props.onSubmitOrder) {
      await props.onSubmitOrder(orderData)
    } else {
      // Otherwise emit event (for backward compatibility)
      emit('submit-order', orderData)
    }
    
    submitSuccess.value = true
    
    // Clear success message after 3 seconds
    setTimeout(() => {
      submitSuccess.value = false
    }, 3000)
  } catch (error: any) {
    console.error('Order submission failed:', error)
    
    // Extract error message
    if (error.data && error.data.message) {
      submitError.value = error.data.message
    } else if (error.message) {
      submitError.value = error.message
    } else {
      submitError.value = 'Failed to submit order. Please try again.'
    }
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="bg-[#1c2936] rounded-t-3xl shadow-[0px_-4px_20px_0px_rgba(0,0,0,0.3)] border-t border-gray-800">
    <div class="px-4 sm:px-6 py-6 space-y-6">
      <!-- Header -->
      <div class="flex items-center justify-between">
        <h2 class="text-white text-lg sm:text-xl font-bold">
          Order Details
        </h2>
        
        <!-- Order Type (Limit only) -->
        <div class="bg-[#233648] rounded-lg px-1 py-1">
          <div class="bg-[#101922] rounded-md px-3 py-1.5 shadow-sm">
            <span class="text-white text-xs font-bold">Limit</span>
          </div>
        </div>
      </div>

      <!-- Price Input -->
      <PriceInput
        v-model="price"
        :current-price="currentPrice"
        :currency="'JPY'"
        @update:model-value="(val) => price = val"
      />

      <!-- Amount Input -->
      <AmountInput
        v-model="amount"
        :min-amount="minAmount"
        :step="step"
        :symbol="symbol"
        @increment="incrementAmount"
        @decrement="decrementAmount"
        @update:model-value="(val) => amount = val"
      />

      <!-- Available Balance -->
      <div class="bg-[rgba(19,127,236,0.1)] border border-[rgba(19,127,236,0.2)] rounded-lg p-3.5 sm:p-3 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <svg
            class="w-4 h-4 text-blue-200"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            aria-hidden="true"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <span class="text-blue-200 text-xs font-semibold">Available Balance</span>
        </div>
        <span class="text-white text-sm font-bold">{{ formattedBalance }}</span>
      </div>

      <!-- Estimated Total -->
      <div class="border-t border-gray-800 pt-4">
        <div class="flex items-center justify-between mb-4">
          <span class="text-slate-400 text-sm font-medium">Estimated Total</span>
          <span class="text-white text-2xl sm:text-3xl font-extrabold tracking-tight">{{ formattedTotal }}</span>
        </div>

        <!-- Submit Button -->
        <button
          @click="handleSubmit"
          :disabled="!isValidOrder || isSubmitting"
          class="w-full bg-[#137fec] hover:bg-[#1068c4] active:bg-[#0e5aa8] disabled:bg-gray-600 disabled:cursor-not-allowed text-white font-bold py-4 sm:py-3.5 rounded-xl shadow-[0px_10px_15px_-3px_rgba(19,127,236,0.3),0px_4px_6px_-4px_rgba(19,127,236,0.3)] transition-all duration-200 flex items-center justify-center gap-2 touch-manipulation min-h-[52px]"
          aria-label="Place buy order"
        >
          <span v-if="!isSubmitting" class="text-base sm:text-sm">Place Buy Order</span>
          <span v-else class="text-base sm:text-sm">Submitting...</span>
          <svg
            v-if="!isSubmitting"
            class="w-6 h-6 sm:w-5 sm:h-5"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            aria-hidden="true"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M13 7l5 5m0 0l-5 5m5-5H6"
            />
          </svg>
          <!-- Loading spinner -->
          <svg
            v-else
            class="animate-spin w-5 h-5"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              class="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              stroke-width="4"
            ></circle>
            <path
              class="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            ></path>
          </svg>
        </button>

        <!-- Success Message -->
        <div
          v-if="submitSuccess"
          class="mt-3 bg-green-900/30 border border-green-500/50 rounded-lg p-3 text-green-400 text-sm text-center font-medium"
          role="alert"
        >
          ✓ Order submitted successfully!
        </div>

        <!-- Error Message -->
        <div
          v-if="submitError"
          class="mt-3 bg-red-900/30 border border-red-500/50 rounded-lg p-3 text-red-400 text-sm text-center font-medium"
          role="alert"
        >
          {{ submitError }}
        </div>

        <!-- Warning if balance insufficient -->
        <div
          v-if="!isBalanceSufficient && estimatedTotal > 0 && !submitError"
          class="mt-3 text-[#fa6238] text-xs text-center font-medium"
          role="alert"
        >
          Insufficient balance for this order
        </div>
      </div>

      <!-- Terms -->
      <div class="text-slate-500 text-[10px] sm:text-xs text-center leading-relaxed">
        Market prices are volatile. By placing this order, you agree to our
        <span class="underline cursor-pointer hover:text-slate-400">Terms of Service</span>.
      </div>
    </div>
  </div>
</template>
