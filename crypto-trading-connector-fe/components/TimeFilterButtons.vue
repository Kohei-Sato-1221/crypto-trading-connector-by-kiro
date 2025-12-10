<script setup lang="ts">
import { useTimeFilter } from '~/composables/useTimeFilter'
import type { TimeFilter } from '~/types/crypto'

const { selectedFilter, setFilter, isSelected } = useTimeFilter()

// Map display labels to TimeFilter values
const filterMap: { display: string; value: TimeFilter }[] = [
  { display: '7D', value: '7d' },
  { display: '30D', value: '30d' },
  { display: '1Y', value: '1y' }
]

const selectFilter = (filter: TimeFilter) => {
  setFilter(filter)
}
</script>

<template>
  <div class="px-4 py-4">
    <div class="flex items-center justify-center gap-2 flex-wrap">
      <button
        v-for="filter in filterMap"
        :key="filter.value"
        @click="selectFilter(filter.value)"
        class="h-10 sm:h-8 px-4 sm:px-3 rounded-full text-xs font-bold transition-all duration-200 touch-manipulation"
        :class="
          isSelected(filter.value)
            ? 'bg-[#137fec] text-white shadow-[0px_10px_15px_-3px_rgba(19,127,236,0.25),0px_4px_6px_-4px_rgba(19,127,236,0.25)]'
            : 'text-[#92adc9] hover:bg-[#233648] active:bg-[#2a3f54]'
        "
      >
        {{ filter.display }}
      </button>
    </div>
  </div>
</template>
