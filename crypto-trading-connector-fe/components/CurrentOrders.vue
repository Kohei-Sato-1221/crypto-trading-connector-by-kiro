<template>
  <div class="px-4 py-4">
    <div class="bg-[#1c2936] border border-gray-800 rounded-2xl shadow-sm p-4">
      <!-- Header -->
      <div class="text-slate-500 text-xs font-semibold mb-4">
        Current Orders
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="h-[200px] flex items-center justify-center">
        <div class="text-slate-400 text-sm">Loading orders...</div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="h-[200px] flex items-center justify-center">
        <div class="text-[#fa6238] text-sm">{{ error }}</div>
      </div>

      <!-- Orders Content -->
      <div v-else class="space-y-6 transition-opacity duration-300" :class="{ 'opacity-75': loading }">
        <!-- Buy Orders Section -->
        <div>
          <div class="text-green-400 text-sm font-semibold mb-3 transition-all duration-300">
            Buy Orders ({{ buyOrders.length }} orders)
          </div>
          
          <!-- Buy Orders Header -->
          <div class="grid grid-cols-3 gap-4 text-xs text-slate-500 font-medium mb-2 px-2">
            <div>Date</div>
            <div class="text-right">Price</div>
            <div class="text-right">Amount</div>
          </div>

          <!-- Buy Orders List -->
          <div v-if="buyOrders.length === 0" class="text-slate-400 text-sm py-4 text-center transition-all duration-300">
            No buy orders
          </div>
          <TransitionGroup v-else name="order-list" tag="div" class="space-y-2">
            <div
              v-for="order in buyOrders"
              :key="order.id"
              class="grid grid-cols-3 gap-4 text-sm py-2 px-2 rounded-lg bg-[#0f1419] border border-gray-800 transition-all duration-300 hover:bg-[#1a2332]"
            >
              <div class="text-slate-300">{{ formatDate(order.createdAt) }}</div>
              <div class="text-right text-slate-300">{{ formatPrice(order.price) }}</div>
              <div class="text-right text-slate-300">{{ formatAmount(order.amount) }}</div>
            </div>
          </TransitionGroup>
        </div>

        <!-- Sell Orders Section -->
        <div>
          <div class="text-red-400 text-sm font-semibold mb-3 transition-all duration-300">
            Sell Orders ({{ sellOrders.length }} orders)
          </div>
          
          <!-- Sell Orders Header -->
          <div class="grid grid-cols-3 gap-4 text-xs text-slate-500 font-medium mb-2 px-2">
            <div>Date</div>
            <div class="text-right">Price</div>
            <div class="text-right">Amount</div>
          </div>

          <!-- Sell Orders List -->
          <div v-if="sellOrders.length === 0" class="text-slate-400 text-sm py-4 text-center transition-all duration-300">
            No sell orders
          </div>
          <TransitionGroup v-else name="order-list" tag="div" class="space-y-2">
            <div
              v-for="order in sellOrders"
              :key="order.id"
              class="grid grid-cols-3 gap-4 text-sm py-2 px-2 rounded-lg bg-[#0f1419] border border-gray-800 transition-all duration-300 hover:bg-[#1a2332]"
            >
              <div class="text-slate-300">{{ formatDate(order.createdAt) }}</div>
              <div class="text-right text-slate-300">{{ formatPrice(order.price) }}</div>
              <div class="text-right text-slate-300">{{ formatAmount(order.amount) }}</div>
            </div>
          </TransitionGroup>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useCurrentOrders } from '~/composables/useCurrentOrders'

// Props
interface Props {
  selectedPair?: 'BTC/JPY' | 'ETH/JPY'
}

const props = withDefaults(defineProps<Props>(), {
  selectedPair: 'BTC/JPY'
})

// Convert props to ref for reactivity
const selectedPairRef = toRef(props, 'selectedPair')

// Use current orders composable
const { buyOrders, sellOrders, loading, error } = useCurrentOrders(selectedPairRef)

// Format date to YYYY/MM/DD HH:MM:SS
const formatDate = (dateString: string): string => {
  const date = new Date(dateString)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  
  return `${year}/${month}/${day} ${hours}:${minutes}:${seconds}`
}

// Format price with JPY formatting
const formatPrice = (price: number): string => {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0
  }).format(price)
}

// Format amount with appropriate decimal places
const formatAmount = (amount: number): string => {
  return new Intl.NumberFormat('ja-JP', {
    minimumFractionDigits: 3,
    maximumFractionDigits: 3
  }).format(amount)
}
</script>

<style scoped>
/* Smooth transitions for order list updates */
.order-list-enter-active,
.order-list-leave-active {
  transition: all 0.3s ease;
}

.order-list-enter-from {
  opacity: 0;
  transform: translateY(-10px);
}

.order-list-leave-to {
  opacity: 0;
  transform: translateY(10px);
}

.order-list-move {
  transition: transform 0.3s ease;
}
</style>