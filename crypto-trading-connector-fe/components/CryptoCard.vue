<script setup lang="ts">
import type { CryptoData } from '~/types/crypto'
import PriceChart from './PriceChart.vue'

interface CryptoDataWithError extends CryptoData {
  hasError?: boolean
  errorMessage?: string
}

interface Props {
  crypto: CryptoDataWithError
}

const props = defineProps<Props>()

const isPositive = computed(() => props.crypto.changePercent > 0)

const formattedPrice = computed(() => {
  return `¥${props.crypto.currentPrice.toLocaleString()}`
})

const formattedChange = computed(() => {
  const sign = isPositive.value ? '+' : ''
  return `${sign}${props.crypto.changePercent.toFixed(1)}%`
})

const handleTrade = () => {
  // Pass the pair as a query parameter
  navigateTo(`/trade?pair=${encodeURIComponent(props.crypto.pair)}`)
}
</script>

<template>
  <div class="bg-[#1c2630] border border-[#2e3e50] rounded-2xl p-5 shadow-sm transition-all duration-300">
    <!-- Header: Icon, Name, Trade Button -->
    <div class="flex items-center justify-between mb-4">
      <div class="flex items-center gap-3">
        <!-- Icon -->
        <div 
          class="w-10 h-10 rounded-full flex items-center justify-center shadow-sm"
          :style="{ backgroundColor: crypto.iconColor }"
        >
          <span class="text-white text-xl font-bold">{{ crypto.icon }}</span>
        </div>
        
        <!-- Name and Pair -->
        <div>
          <h3 class="text-lg font-bold text-white">{{ crypto.name }}</h3>
          <p class="text-xs text-slate-400 tracking-wide">{{ crypto.pair }}</p>
        </div>
      </div>
      
      <!-- Trade Button -->
      <button
        @click="handleTrade"
        :disabled="crypto.hasError"
        class="px-4 py-2 text-xs font-bold rounded-lg transition-colors min-w-[44px] min-h-[44px] flex items-center justify-center"
        :class="crypto.hasError 
          ? 'bg-[rgba(250,98,56,0.1)] text-[#fa6238] cursor-not-allowed opacity-50' 
          : 'bg-[rgba(19,127,236,0.1)] text-[#137fec] hover:bg-[rgba(19,127,236,0.2)]'"
      >
        Trade
      </button>
    </div>

    <!-- Error State -->
    <div v-if="crypto.hasError" class="mb-4">
      <div class="bg-[rgba(250,98,56,0.1)] border border-[rgba(250,98,56,0.3)] rounded-lg p-4">
        <div class="flex items-start gap-3">
          <!-- Error Icon -->
          <svg class="w-5 h-5 text-[#fa6238] flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
          
          <div class="flex-1">
            <div class="text-[#fa6238] font-semibold text-sm mb-1">
              {{ crypto.errorMessage || 'Failed to load data' }}
            </div>
            <div class="text-white/40 text-xs">
              Unable to fetch data from backend API. Please check if the server is running.
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Normal State -->
    <template v-else>
      <!-- Price and Change -->
      <div class="mb-4">
        <div class="text-3xl font-bold text-white tracking-tight mb-2 transition-all duration-300">
          {{ formattedPrice }}
        </div>
        
        <div class="flex items-center gap-2">
          <div 
            class="px-2 py-1 rounded flex items-center gap-1.5 transition-all duration-300"
            :class="isPositive ? 'bg-[rgba(11,218,91,0.1)]' : 'bg-[rgba(250,98,56,0.1)]'"
          >
            <!-- Arrow Icon -->
            <svg 
              class="w-4 h-4 transition-transform duration-300"
              :class="isPositive ? 'text-[#0bda5b]' : 'text-[#fa6238] rotate-180'"
              fill="none" 
              viewBox="0 0 24 24" 
              stroke="currentColor"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 10l7-7m0 0l7 7m-7-7v18" />
            </svg>
            
            <span 
              class="text-sm font-semibold"
              :class="isPositive ? 'text-[#0bda5b]' : 'text-[#fa6238]'"
            >
              {{ formattedChange }}
            </span>
          </div>
          
          <span class="text-xs text-slate-400">Last 7 days</span>
        </div>
      </div>

      <!-- Chart -->
      <div class="mt-4">
        <PriceChart 
          :data="crypto.chartData" 
          :isPositive="isPositive"
          :currency="crypto.symbol"
        />
      </div>
    </template>
  </div>
</template>
