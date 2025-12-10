<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useAutoRefresh } from '~/composables/useAutoRefresh'
import { useCryptoData } from '~/composables/useCryptoData'
import { useTimeFilter } from '~/composables/useTimeFilter'
import MarketHeader from '~/components/MarketHeader.vue'
import CryptoCard from '~/components/CryptoCard.vue'
import NavigationBar from '~/components/NavigationBar.vue'

// Use crypto data composable
const { cryptoData, loading, error, fetchCryptoData, useMockData } = useCryptoData()

// Use time filter composable
const { selectedFilter } = useTimeFilter()

// Fetch data on mount
onMounted(async () => {
  try {
    await fetchCryptoData(selectedFilter.value)
  } catch (e) {
    console.error('Failed to fetch crypto data:', e)
  }
})

// Watch for time filter changes and refetch data
watch(selectedFilter, async (newFilter) => {
  try {
    await fetchCryptoData(newFilter)
  } catch (e) {
    console.error('Failed to fetch crypto data for filter:', newFilter, e)
  }
})

// Auto-refresh every 5 seconds (only runs on client)
const updatePrices = async () => {
  try {
    await fetchCryptoData(selectedFilter.value)
  } catch (e) {
    console.error('Failed to refresh crypto data:', e)
  }
}

const autoRefresh = useAutoRefresh(updatePrices, 5000)

// Start auto-refresh on client side
if (process.client) {
  autoRefresh.start()
}
</script>

<template>
  <div class="min-h-screen bg-[#101922] pb-20">
    <!-- Header (sticky) -->
    <MarketHeader />

    <!-- Scrollable Content Area -->
    <main class="px-4 py-4">
      <div class="space-y-4 max-w-2xl mx-auto">
        <!-- Loading State (initial load only) -->
        <div v-if="loading && cryptoData.length === 0" class="text-center py-8">
          <div class="text-white/60">Loading...</div>
        </div>

        <!-- Crypto Cards (show individual errors per card) -->
        <CryptoCard
          v-for="crypto in cryptoData"
          :key="crypto.id"
          :crypto="crypto"
        />
      </div>
    </main>

    <!-- Navigation Bar (fixed) -->
    <NavigationBar />
  </div>
</template>
