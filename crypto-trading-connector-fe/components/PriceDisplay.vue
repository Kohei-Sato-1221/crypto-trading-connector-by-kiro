<script setup lang="ts">
interface Props {
  currentPrice: number
  priceChange: number
  currency: string
}

const props = defineProps<Props>()

const isPositive = computed(() => props.priceChange > 0)

const formattedPrice = computed(() => {
  return `¥${props.currentPrice.toLocaleString()}`
})

const formattedChange = computed(() => {
  const sign = isPositive.value ? '+' : ''
  return `${sign}${props.priceChange.toFixed(1)}%`
})
</script>

<template>
  <div class="px-4 py-6 sm:py-8 text-center">
    <!-- Label -->
    <div class="text-[#92adc9] text-xs sm:text-sm font-medium tracking-wider uppercase mb-3">
      Current Price
    </div>

    <!-- Price -->
    <div class="text-white text-3xl sm:text-4xl md:text-5xl font-extrabold tracking-tight mb-3">
      {{ formattedPrice }}
    </div>

    <!-- Change Badge -->
    <div class="flex items-center justify-center">
      <div
        class="px-3 py-1.5 rounded-full flex items-center gap-1.5 transition-all duration-300"
        :class="
          isPositive
            ? 'bg-[rgba(11,218,91,0.1)]'
            : 'bg-[rgba(250,98,56,0.1)]'
        "
        role="status"
        :aria-label="`Price change: ${formattedChange} in 24 hours`"
      >
        <!-- Arrow Icon -->
        <svg
          class="w-4 h-4 transition-transform duration-300"
          :class="
            isPositive
              ? 'text-[#0bda5b]'
              : 'text-[#fa6238] rotate-180'
          "
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          aria-hidden="true"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M5 10l7-7m0 0l7 7m-7-7v18"
          />
        </svg>

        <span
          class="text-sm font-bold"
          :class="
            isPositive ? 'text-[#0bda5b]' : 'text-[#fa6238]'
          "
        >
          {{ formattedChange }} (24h)
        </span>
      </div>
    </div>
  </div>
</template>
