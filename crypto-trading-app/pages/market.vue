<script setup lang="ts">
import { ref } from 'vue'
import type { CryptoData } from '~/types/crypto'
import { useAutoRefresh } from '~/composables/useAutoRefresh'
import { updateMockPrice, getMockCryptoData } from '~/utils/mockData'
import MarketHeader from '~/components/MarketHeader.vue'
import CryptoCard from '~/components/CryptoCard.vue'
import NavigationBar from '~/components/NavigationBar.vue'

// Initialize with mock data immediately (works for both SSR and client)
const cryptoList = ref<CryptoData[]>(getMockCryptoData())

// Auto-refresh every 5 seconds (only runs on client)
const updatePrices = () => {
  cryptoList.value = cryptoList.value.map(crypto => updateMockPrice(crypto))
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
        <CryptoCard
          v-for="crypto in cryptoList"
          :key="crypto.id"
          :crypto="crypto"
        />
      </div>
    </main>

    <!-- Navigation Bar (fixed) -->
    <NavigationBar />
  </div>
</template>
