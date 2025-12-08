<script setup lang="ts">
import { useTimeFilter } from '~/composables/useTimeFilter'
import type { TimeFilter } from '~/types/crypto'

const { selectedFilter, setFilter, isSelected } = useTimeFilter()

const timeFilters: { value: TimeFilter; label: string }[] = [
  { value: '24h', label: '24h' },
  { value: '7d', label: '7d' },
  { value: '30d', label: '30d' },
  { value: '1y', label: '1y' },
  { value: 'all', label: 'All' }
]

const handleFilterChange = (filter: TimeFilter) => {
  setFilter(filter)
}
</script>

<template>
  <header class="sticky top-0 z-50 backdrop-blur-md bg-[rgba(16,25,34,0.95)] border-b border-[#2e3e50]">
    <div class="px-4 py-4">
      <!-- App Title -->
      <div class="text-center mb-4">
        <h1 class="text-lg font-bold text-white tracking-tight">
          Crypto Trading Connector
        </h1>
      </div>

      <!-- Market Status -->
      <div class="flex items-center justify-center mb-4">
        <div class="flex items-center gap-2">
          <div class="relative">
            <div class="w-2 h-2 bg-[#0bda5b] rounded-full"></div>
            <div class="absolute inset-0 w-2 h-2 bg-[#0bda5b] rounded-full opacity-75 animate-ping"></div>
          </div>
          <p class="text-xs text-slate-400">
            Market is Open • Real-time Data
          </p>
        </div>
      </div>

      <!-- Time Filter Buttons -->
      <div class="flex items-center gap-2 overflow-x-auto pb-2">
        <button
          v-for="filter in timeFilters"
          :key="filter.value"
          @click="handleFilterChange(filter.value)"
          :class="[
            'px-4 py-2 rounded-full text-sm font-semibold transition-colors whitespace-nowrap min-w-[44px]',
            isSelected(filter.value)
              ? 'bg-[#137fec] text-white shadow-[0px_1px_2px_0px_rgba(19,127,236,0.2)]'
              : 'bg-[#2e3e50] text-slate-400'
          ]"
        >
          {{ filter.label }}
        </button>
      </div>
    </div>
  </header>
</template>

<style scoped>
@keyframes ping {
  75%, 100% {
    transform: scale(2);
    opacity: 0;
  }
}

.animate-ping {
  animation: ping 1s cubic-bezier(0, 0, 0.2, 1) infinite;
}
</style>
