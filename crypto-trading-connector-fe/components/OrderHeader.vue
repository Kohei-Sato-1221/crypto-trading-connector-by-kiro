<script setup lang="ts">
interface Props {
  selectedPair: 'BTC/JPY' | 'ETH/JPY'
}

interface Emits {
  (e: 'update:selectedPair', pair: 'BTC/JPY' | 'ETH/JPY'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const pairs = ['BTC/JPY', 'ETH/JPY'] as const

const selectPair = (pair: 'BTC/JPY' | 'ETH/JPY') => {
  emit('update:selectedPair', pair)
}
</script>

<template>
  <div class="px-4 pt-4">
    <div class="bg-[#233648] rounded-xl p-1 flex gap-1">
      <button
        v-for="pair in pairs"
        :key="pair"
        @click="selectPair(pair)"
        class="flex-1 py-2.5 rounded-lg text-sm font-bold transition-all duration-200"
        :class="
          selectedPair === pair
            ? 'bg-[#101922] text-[#92adc9] shadow-sm'
            : 'text-[#92adc9] hover:bg-[#101922]/50'
        "
      >
        {{ pair }}
      </button>
    </div>
  </div>
</template>
