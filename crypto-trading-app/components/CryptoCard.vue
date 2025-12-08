<script setup lang="ts">
import type { CryptoData } from '~/types/crypto'
import PriceChart from './PriceChart.vue'

interface Props {
  crypto: CryptoData
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
  navigateTo('/trade')
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
        class="px-4 py-2 bg-[rgba(19,127,236,0.1)] text-[#137fec] text-xs font-bold rounded-lg hover:bg-[rgba(19,127,236,0.2)] transition-colors min-w-[44px] min-h-[44px] flex items-center justify-center"
      >
        Trade
      </button>
    </div>

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
  </div>
</template>
