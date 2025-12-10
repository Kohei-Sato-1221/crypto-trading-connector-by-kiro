<template>
  <div class="min-h-screen bg-[#101922] pb-20">
    <!-- Header -->
    <div class="bg-[#1a2332] border-b border-[#2a3441]">
      <div class="flex items-center justify-between px-4 py-3">
        <button class="w-10 h-10 flex items-center justify-center rounded-lg hover:bg-[#2a3441] transition-colors">
          <svg class="w-6 h-7" viewBox="0 0 24 28" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M15 18L9 12L15 6" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </button>
        <h1 class="text-white text-lg font-semibold">Trade History</h1>
        <div class="w-10"></div> <!-- Spacer for centering -->
      </div>
    </div>

    <!-- Main Content -->
    <main class="py-6">
      <!-- Loading State -->
      <div v-if="loading" class="text-center py-8">
        <div class="text-white/60">Loading...</div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center py-8 px-4">
        <div class="text-red-400 mb-4">{{ error.message }}</div>
        <button
          @click="refresh"
          class="px-4 py-2 bg-[#2a3441] text-white rounded-lg hover:bg-[#3a4451] transition-colors"
        >
          Retry
        </button>
      </div>

      <!-- Content -->
      <div v-else>
        <!-- Statistics Cards -->
        <TradeStatistics :statistics="statistics" />

        <!-- Filters -->
        <TradeFilters
          :time-filter="filters.timeFilter"
          :asset-filter="filters.assetFilter"
          @time-filter-change="setTimeFilter"
          @asset-filter-change="setAssetFilter"
        />

        <!-- Transaction Log -->
        <TransactionLog
          :transactions="transactions"
          :can-load-more="canLoadMore"
          @load-more="loadMoreTransactions"
        />
      </div>
    </main>

    <!-- Navigation Bar -->
    <NavigationBar />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useTradeHistory } from '~/composables/useTradeHistory'
import TradeStatistics from '~/components/TradeStatistics.vue'
import TradeFilters from '~/components/TradeFilters.vue'
import TransactionLog from '~/components/TransactionLog.vue'
import NavigationBar from '~/components/NavigationBar.vue'

// Use trade history composable
const {
  statistics,
  transactions,
  loading,
  error,
  filters,
  canLoadMore,
  setTimeFilter,
  setAssetFilter,
  loadMoreTransactions,
  refresh,
  initialize
} = useTradeHistory()

// Initialize data on mount
onMounted(async () => {
  await initialize()
})
</script>
