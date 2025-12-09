<script setup lang="ts">
interface Props {
  modelValue: number
  minAmount: number
  step: number
  symbol: string
}

interface Emits {
  (e: 'update:modelValue', value: number): void
  (e: 'increment'): void
  (e: 'decrement'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const formattedValue = computed({
  get: () => props.modelValue.toString(),
  set: (value: string) => {
    const numValue = parseFloat(value)
    if (!isNaN(numValue) && numValue >= 0) {
      emit('update:modelValue', numValue)
    }
  }
})

const handleIncrement = () => {
  emit('increment')
}

const handleDecrement = () => {
  emit('decrement')
}
</script>

<template>
  <div class="space-y-3">
    <!-- Label -->
    <div class="text-[#92adc9] text-xs font-semibold">
      Amount ({{ symbol }})
    </div>

    <!-- Input Field with Increment/Decrement Buttons -->
    <div class="bg-[#233648] rounded-xl overflow-hidden">
      <div class="flex items-center px-4 py-4 sm:py-3">
        <!-- Decrement Button -->
        <button
          @click="handleDecrement"
          class="w-10 h-10 sm:w-8 sm:h-8 flex items-center justify-center rounded-lg hover:bg-[#2a3f54] active:bg-[#324a62] transition-colors touch-manipulation"
          :disabled="modelValue <= minAmount"
          :class="modelValue <= minAmount ? 'opacity-50 cursor-not-allowed' : ''"
          aria-label="Decrease amount"
        >
          <svg
            class="w-5 h-5 sm:w-4 sm:h-4 text-white transform scale-y-[-1]"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M5 10l7-7m0 0l7 7m-7-7v18"
            />
          </svg>
        </button>

        <!-- Input -->
        <input
          :value="formattedValue"
          @input="formattedValue = ($event.target as HTMLInputElement).value"
          type="text"
          inputmode="decimal"
          class="flex-1 bg-transparent text-white text-xl sm:text-lg font-bold text-center outline-none mx-2 min-h-[24px]"
          placeholder="0"
          aria-label="Amount"
        />

        <!-- Increment Button -->
        <button
          @click="handleIncrement"
          class="w-10 h-10 sm:w-8 sm:h-8 flex items-center justify-center rounded-lg hover:bg-[#2a3f54] active:bg-[#324a62] transition-colors touch-manipulation"
          aria-label="Increase amount"
        >
          <svg
            class="w-5 h-5 sm:w-4 sm:h-4 text-white"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M5 10l7-7m0 0l7 7m-7-7v18"
            />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>
