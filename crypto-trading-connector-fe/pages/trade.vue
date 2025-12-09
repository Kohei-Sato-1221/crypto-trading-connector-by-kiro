<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useOrderData } from '~/composables/useOrderData'
import OrderHeader from '~/components/OrderHeader.vue'
import PriceDisplay from '~/components/PriceDisplay.vue'
import TimeFilterButtons from '~/components/TimeFilterButtons.vue'
import PriceChart from '~/components/PriceChart.vue'
import OrderForm from '~/components/OrderForm.vue'
import NavigationBar from '~/components/NavigationBar.vue'

// State
const selectedPair = ref<'BTC/JPY' | 'ETH/JPY'>('BTC/JPY')
const selectedTimeFilter = ref('7D')

// Use order data composable
const {
  currentPrice,
  priceChange,
  chartData,
  availableBalance,
  loading,
  error,
  fetchAllData
} = useOrderData(selectedPair)

// Computed values
const symbol = computed(() => {
  return selectedPair.value === 'BTC/JPY' ? 'BTC' : 'ETH'
})

const isPositive = computed(() => priceChange.value > 0)

// Transform chart data for PriceChart component
const transformedChartData = computed(() => {
  return chartData.value.map(point => {
    const date = new Date(point.timestamp)
    const day = `${date.getMonth() + 1}/${date.getDate()}`
    return {
      day,
      price: point.price
    }
  })
})

// Fetch data on mount
onMounted(async () => {
  try {
    await fetchAllData(selectedTimeFilter.value)
  } catch (e) {
    console.error('Failed to fetch order data:', e)
  }
})

// Watch for pair changes
watch(selectedPair, async () => {
  try {
    await fetchAllData(selectedTimeFilter.value)
  } catch (e) {
    console.error('Failed to fetch order data:', e)
  }
})

// Watch for time filter changes
watch(selectedTimeFilter, async (newFilter) => {
  try {
    await fetchAllData(newFilter)
  } catch (e) {
    console.error('Failed to fetch chart data:', e)
  }
})

// Handle order submission
const handleSubmitOrder = (order: any) => {
  console.log('Order submitted:', order)
  // TODO: Implement API call in integration phase
  alert(`Order placed successfully!\n\nPair: ${order.pair}\nPrice: ¥${order.price.toLocaleString()}\nAmount: ${order.amount}\nTotal: ¥${order.estimatedTotal.toLocaleString()}`)
}
</script>

<template>
  <div class="min-h-screen bg-[#101922] pb-20">
    <!-- Main Content Container with max width for larger screens -->
    <div class="max-w-2xl mx-auto">
      <!-- Currency Pair Header -->
      <OrderHeader
        :selected-pair="selectedPair"
        @update:selected-pair="(pair) => selectedPair = pair"
      />

      <!-- Price Display -->
      <PriceDisplay
        :current-price="currentPrice"
        :price-change="priceChange"
        :currency="symbol"
      />

      <!-- Time Filter Buttons -->
      <TimeFilterButtons
        :selected-filter="selectedTimeFilter"
        @update:selected-filter="(filter) => selectedTimeFilter = filter"
      />

      <!-- Price Chart -->
      <div class="px-4 py-4">
        <div class="bg-[#1c2936] border border-gray-800 rounded-2xl shadow-sm p-4">
          <div class="text-slate-500 text-xs font-semibold mb-4">
            Price History (Last {{ selectedTimeFilter }})
          </div>
          
          <!-- Loading State -->
          <div v-if="loading" class="h-[200px] flex items-center justify-center">
            <div class="text-slate-400 text-sm">Loading chart...</div>
          </div>

          <!-- Error State -->
          <div v-else-if="error" class="h-[200px] flex items-center justify-center">
            <div class="text-[#fa6238] text-sm">Failed to load chart data</div>
          </div>

          <!-- Chart -->
          <PriceChart
            v-else
            :data="transformedChartData"
            :is-positive="isPositive"
            :currency="symbol"
          />
        </div>
      </div>

      <!-- Order Form -->
      <OrderForm
        :selected-pair="selectedPair"
        :current-price="currentPrice"
        :available-balance="availableBalance"
        @submit-order="handleSubmitOrder"
      />
    </div>

    <!-- Navigation Bar (full width) -->
    <NavigationBar />
  </div>
</template>
