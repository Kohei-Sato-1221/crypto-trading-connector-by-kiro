<script setup lang="ts">
interface Props {
  selectedFilter: string
}

interface Emits {
  (e: 'update:selectedFilter', filter: string): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const filters = ['7D', '30D', '1Y'] as const

const selectFilter = (filter: string) => {
  emit('update:selectedFilter', filter)
}
</script>

<template>
  <div class="px-4 py-4">
    <div class="flex items-center justify-center gap-2 flex-wrap">
      <button
        v-for="filter in filters"
        :key="filter"
        @click="selectFilter(filter)"
        class="h-10 sm:h-8 px-4 sm:px-3 rounded-full text-xs font-bold transition-all duration-200 touch-manipulation"
        :class="
          selectedFilter === filter
            ? 'bg-[#137fec] text-white shadow-[0px_10px_15px_-3px_rgba(19,127,236,0.25),0px_4px_6px_-4px_rgba(19,127,236,0.25)]'
            : 'text-[#92adc9] hover:bg-[#233648] active:bg-[#2a3f54]'
        "
      >
        {{ filter }}
      </button>
    </div>
  </div>
</template>
