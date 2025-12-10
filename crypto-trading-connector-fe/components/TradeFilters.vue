<template>
  <div class="px-3 mb-6 max-w-4xl mx-auto">
    <!-- Time Filter -->
    <div class="bg-[#1a2332] rounded-lg p-3 mb-4">
      <div class="flex justify-center space-x-2 sm:space-x-4">
        <button
          @click="$emit('timeFilterChange', 'all')"
          :class="[
            'px-4 sm:px-6 py-2 text-sm font-medium rounded-md transition-colors touch-manipulation',
            timeFilter === 'all'
              ? 'text-white bg-[#2a3441]'
              : 'text-white/60 hover:text-white hover:bg-[#2a3441]/50 active:bg-[#2a3441]/70'
          ]"
        >
          All Time
        </button>
        <button
          @click="$emit('timeFilterChange', '7days')"
          :class="[
            'px-4 sm:px-6 py-2 text-sm font-medium rounded-md transition-colors touch-manipulation',
            timeFilter === '7days'
              ? 'text-white bg-[#2a3441]'
              : 'text-white/60 hover:text-white hover:bg-[#2a3441]/50 active:bg-[#2a3441]/70'
          ]"
        >
          Last 7 Days
        </button>
      </div>
    </div>

    <!-- Asset Filter -->
    <div class="flex space-x-2 overflow-x-auto pb-2 scrollbar-hide">
      <button
        @click="$emit('assetFilterChange', 'all')"
        :class="[
          'flex-shrink-0 px-4 py-2 text-sm font-medium rounded-lg transition-colors touch-manipulation',
          assetFilter === 'all'
            ? 'text-white bg-[#2a3441]'
            : 'text-white/60 hover:text-white hover:bg-[#2a3441]/50 active:bg-[#2a3441]/70'
        ]"
      >
        All Assets
      </button>
      
      <button
        @click="$emit('assetFilterChange', 'BTC')"
        :class="[
          'flex-shrink-0 flex items-center px-4 py-2 text-sm font-medium rounded-lg transition-colors touch-manipulation',
          assetFilter === 'BTC'
            ? 'text-white bg-[#2a3441]'
            : 'text-white/60 hover:text-white hover:bg-[#2a3441]/50 active:bg-[#2a3441]/70'
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
          'flex-shrink-0 flex items-center px-4 py-2 text-sm font-medium rounded-lg transition-colors touch-manipulation',
          assetFilter === 'ETH'
            ? 'text-white bg-[#2a3441]'
            : 'text-white/60 hover:text-white hover:bg-[#2a3441]/50 active:bg-[#2a3441]/70'
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