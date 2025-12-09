<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useOrderForm } from '~/composables/useOrderForm'
import PriceInput from './PriceInput.vue'
import AmountInput from './AmountInput.vue'

interface Props {
  selectedPair: 'BTC/JPY' | 'ETH/JPY'
  currentPrice: number
  availableBalance: number
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

// Reinitialize when pair changes
watch(pair, () => {
  initializeDefaults()
})

const symbol = computed(() => {
  return props.selectedPair === 'BTC/JPY' ? 'BTC' : 'ETH'
})

const formattedBalance = computed(() => {
  return `¥ ${props.availableBalance.toLocaleString()}`
})

const formattedTotal = computed(() => {
  return `¥ ${estimatedTotal.value.toLocaleString()}`
})

const handleSubmit = () => {
  const validation = validateOrder()
  
  if (!validation.valid) {
    alert(validation.error)
    return
  }

  emit('submit-order', {
    pair: props.selectedPair,
    orderType: 'limit',
    price: price.value,
    amount: amount.value,
    estimatedTotal: estimatedTotal.value
  })
}
</script>

<template>
  <div class="bg-[#1c2936] rounded-t-3xl shadow-[0px_-4px_20px_0px_rgba(0,0,0,0.3)] border-t border-gray-800">
    <div class="px-6 py-6 space-y-6">
      <!-- Header -->
      <div class="flex items-center justify-between">
        <h2 class="text-white text-lg font-bold">
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
      <div class="bg-[rgba(19,127,236,0.1)] border border-[rgba(19,127,236,0.2)] rounded-lg p-3 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <svg
            class="w-4 h-4 text-blue-200"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
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
          <span class="text-white text-2xl font-extrabold tracking-tight">{{ formattedTotal }}</span>
        </div>

        <!-- Submit Button -->
        <button
          @click="handleSubmit"
          :disabled="!isValidOrder"
          class="w-full bg-[#137fec] hover:bg-[#1068c4] disabled:bg-gray-600 disabled:cursor-not-allowed text-white font-bold py-4 rounded-xl shadow-[0px_10px_15px_-3px_rgba(19,127,236,0.3),0px_4px_6px_-4px_rgba(19,127,236,0.3)] transition-all duration-200 flex items-center justify-center gap-2"
        >
          <span>Place Buy Order</span>
          <svg
            class="w-6 h-6"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M13 7l5 5m0 0l-5 5m5-5H6"
            />
          </svg>
        </button>

        <!-- Warning if balance insufficient -->
        <div
          v-if="!isBalanceSufficient && estimatedTotal > 0"
          class="mt-3 text-[#fa6238] text-xs text-center"
        >
          Insufficient balance for this order
        </div>
      </div>

      <!-- Terms -->
      <div class="text-slate-500 text-[10px] text-center leading-relaxed">
        Market prices are volatile. By placing this order, you agree to our
        <span class="underline">Terms of Service</span>.
      </div>
    </div>
  </div>
</template>
