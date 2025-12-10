<template>
  <div class="grid grid-cols-2 gap-3 px-3 mb-6">
    <!-- Total Profit Card -->
    <div class="bg-[#1a2332] border border-[#2a3441] rounded-lg p-4 shadow-sm">
      <div class="flex items-center mb-2">
        <div class="w-5 h-6 mr-3">
          <svg viewBox="0 0 20 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M10 2L18 7V17L10 22L2 17V7L10 2Z" stroke="#10b981" stroke-width="2" fill="none"/>
          </svg>
        </div>
        <span class="text-white text-sm font-medium">Total Profit</span>
      </div>
      
      <div class="mb-3">
        <span class="text-white text-2xl font-bold">
          {{ formatCurrency(statistics.totalProfit) }}
        </span>
      </div>
      
      <div class="flex items-center">
        <div class="w-3.5 h-4 mr-2">
          <svg viewBox="0 0 14 16" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M7 1L13 6L7 11L1 6L7 1Z" fill="#10b981"/>
          </svg>
        </div>
        <span class="text-[#10b981] text-sm font-medium">
          +{{ formatPercentage(statistics.profitPercentage) }}%
        </span>
      </div>
    </div>

    <!-- Executions Card -->
    <div class="bg-[#1a2332] border border-[#2a3441] rounded-lg p-4 shadow-sm">
      <div class="flex items-center mb-2">
        <div class="w-5 h-6 mr-3">
          <svg viewBox="0 0 20 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M10 2L18 7V17L10 22L2 17V7L10 2Z" stroke="#6366f1" stroke-width="2" fill="none"/>
          </svg>
        </div>
        <span class="text-white text-sm font-medium">Executions</span>
      </div>
      
      <div class="mb-3">
        <span class="text-white text-2xl font-bold">
          {{ statistics.executionCount }}
        </span>
      </div>
      
      <div>
        <span class="text-white/60 text-sm">Sell Orders</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { TradeStatistics } from '~/types/tradeHistory'

interface Props {
  statistics: TradeStatistics
}

defineProps<Props>()

/**
 * Format currency value with thousand separators and 1 decimal place
 */
const formatCurrency = (value: number): string => {
  return `¥${value.toLocaleString('ja-JP', { 
    minimumFractionDigits: 0,
    maximumFractionDigits: 1 
  })}`
}

/**
 * Format percentage with 1 decimal place
 */
const formatPercentage = (value: number): string => {
  return value.toFixed(1)
}
</script>