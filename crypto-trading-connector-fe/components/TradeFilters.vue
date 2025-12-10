<template>
  <div class="px-3 mb-6">
    <!-- Time Filter -->
    <div class="bg-[#1a2332] rounded-lg p-3 mb-4">
      <div class="flex justify-center space-x-4">
        <button
          @click="$emit('timeFilterChange', 'all')"
          :class="[
            'px-6 py-2 text-sm font-medium rounded-md transition-colors',
            timeFilter === 'all'
              ? 'text-white bg-[#2a3441]'
              : 'text-white/60 hover:text-white hover:bg-[#2a3441]/50'
          ]"
        >
          All Time
        </button>
        <button
          @click="$emit('timeFilterChange', '7days')"
          :class="[
            'px-6 py-2 text-sm font-medium rounded-md transition-colors',
            timeFilter === '7days'
              ? 'text-white bg-[#2a3441]'
              : 'text-white/60 hover:text-white hover:bg-[#2a3441]/50'
          ]"
        >
          Last 7 Days
        </button>
      </div>
    </div>

    <!-- Asset Filter -->
    <div class="flex space-x-2 overflow-x-auto pb-2">
      <button
        @click="$emit('assetFilterChange', 'all')"
        :class="[
          'flex-shrink-0 px-4 py-2 text-sm font-medium rounded-lg transition-colors',
          assetFilter === 'all'
            ? 'text-white bg-[#2a3441]'
            : 'text-white/60 hover:text-white hover:bg-[#2a3441]/50'
        ]"
      >
        All Assets
      </button>
      
      <button
        @click="$emit('assetFilterChange', 'BTC')"
        :class="[
          'flex-shrink-0 flex items-center px-4 py-2 text-sm font-medium rounded-lg transition-colors',
          assetFilter === 'BTC'
            ? 'text-white bg-[#2a3441]'
            : 'text-white/60 hover:text-white hover:bg-[#2a3441]/50'
        ]"
      >
        <div class="w-4 h-5 mr-2">
          <svg viewBox="0 0 16 20" fill="none" xmlns="http://www.w3.org/2000/svg">
            <circle cx="8" cy="10" r="7" fill="#f7931a"/>
            <text x="8" y="14" text-anchor="middle" fill="white" font-size="8" font-weight="bold">₿</text>
          </svg>
        </div>
        BTC
      </button>
      
      <button
        @click="$emit('assetFilterChange', 'ETH')"
        :class="[
          'flex-shrink-0 flex items-center px-4 py-2 text-sm font-medium rounded-lg transition-colors',
          assetFilter === 'ETH'
            ? 'text-white bg-[#2a3441]'
            : 'text-white/60 hover:text-white hover:bg-[#2a3441]/50'
        ]"
      >
        <div class="w-4 h-5 mr-2">
          <svg viewBox="0 0 16 20" fill="none" xmlns="http://www.w3.org/2000/svg">
            <circle cx="8" cy="10" r="7" fill="#627eea"/>
            <text x="8" y="14" text-anchor="middle" fill="white" font-size="8" font-weight="bold">Ξ</text>
          </svg>
        </div>
        ETH
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  timeFilter: 'all' | '7days'
  assetFilter: 'all' | 'BTC' | 'ETH'
}

interface Emits {
  timeFilterChange: [filter: 'all' | '7days']
  assetFilterChange: [filter: 'all' | 'BTC' | 'ETH']
}

defineProps<Props>()
defineEmits<Emits>()
</script>