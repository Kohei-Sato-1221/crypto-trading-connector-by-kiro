<script setup lang="ts">
interface Props {
  modelValue: number
  currentPrice: number
  currency: string
}

interface Emits {
  (e: 'update:modelValue', value: number): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const discountPercentages = [99, 97, 95] as const

const formattedValue = computed({
  get: () => props.modelValue.toLocaleString(),
  set: (value: string) => {
    const numValue = parseFloat(value.replace(/,/g, ''))
    if (!isNaN(numValue)) {
      emit('update:modelValue', numValue)
    }
  }
})

const setDiscountPrice = (percentage: number) => {
  const discounted = Math.floor(props.currentPrice * (percentage / 100))
  emit('update:modelValue', discounted)
}
</script>

<template>
  <div class="space-y-3">
    <!-- Label -->
    <div class="text-[#92adc9] text-xs font-semibold">
      Price (JPY)
    </div>

    <!-- Input Field -->
    <div class="bg-[#233648] rounded-xl overflow-hidden">
      <div class="flex items-center px-4 py-4 sm:py-3">
        <input
          :value="formattedValue"
          @input="formattedValue = ($event.target as HTMLInputElement).value"
          type="text"
          inputmode="numeric"
          class="flex-1 bg-transparent text-white text-xl sm:text-lg font-bold text-right outline-none min-h-[24px]"
          placeholder="0"
        />
        <div class="ml-3 text-slate-400 text-sm font-semibold">
          {{ currency }}
        </div>
      </div>
    </div>

    <!-- Discount Buttons -->
    <div class="flex gap-2">
      <button
        v-for="percentage in discountPercentages"
        :key="percentage"
        @click="setDiscountPrice(percentage)"
        class="flex-1 bg-[#233648] hover:bg-[#2a3f54] active:bg-[#324a62] text-slate-400 text-xs font-bold py-3 sm:py-2 rounded-lg transition-colors touch-manipulation min-h-[44px] sm:min-h-0"
      >
        {{ percentage }}%
      </button>
    </div>
  </div>
</template>
