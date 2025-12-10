<template>
  <div class="px-3 max-w-4xl mx-auto">
    <!-- Header -->
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-white text-xl font-semibold">Transaction Log</h3>
      <button class="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-[#2a3441] transition-colors">
        <svg class="w-5 h-6" viewBox="0 0 20 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M10 2L18 7V17L10 22L2 17V7L10 2Z" stroke="white" stroke-width="2" fill="none"/>
        </svg>
      </button>
    </div>

    <!-- Transaction List -->
    <div class="space-y-4">
      <div
        v-for="transaction in transactions"
        :key="transaction.id"
        class="bg-[#1a2332] border border-[#2a3441] rounded-lg p-4 shadow-sm hover:bg-[#1e2738] transition-colors touch-manipulation"
      >
        <!-- Transaction Header -->
        <div class="flex items-center justify-between mb-4 flex-wrap sm:flex-nowrap gap-2">
          <div class="flex items-center">
            <div class="w-10 h-10 rounded-lg bg-[#2a3441] flex items-center justify-center mr-3">
              <div class="w-6 h-7">
                <svg viewBox="0 0 24 28" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <circle cx="12" cy="14" r="10" :fill="getCryptoColor(transaction.cryptocurrency)"/>
                  <text x="12" y="18" text-anchor="middle" fill="white" font-size="10" font-weight="bold">
                    {{ getCryptoSymbol(transaction.cryptocurrency) }}
                  </text>
                </svg>
              </div>
            </div>
            <div>
              <div class="text-white text-sm font-medium">{{ transaction.cryptocurrency }}</div>
              <div class="text-white/60 text-xs">{{ formatTimestamp(transaction.timestamp) }}</div>
            </div>
          </div>
          
          <div class="text-right">
            <div class="text-[#10b981] text-sm font-medium">
              +{{ formatCurrency(transaction.profit) }}
            </div>
            <div class="bg-[#10b981]/20 text-[#10b981] text-xs px-2 py-1 rounded mt-1">
              Sell
            </div>
          </div>
        </div>

        <!-- Divider -->
        <div class="h-px bg-[#2a3441] mb-4"></div>

        <!-- Transaction Details -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <div class="text-white/60 text-xs mb-1">Amount</div>
            <div class="text-white text-sm font-medium">
              {{ formatAmount(transaction.amount, transaction.cryptocurrency) }}
            </div>
          </div>
          <div class="sm:text-right">
            <div class="text-white/60 text-xs mb-1">Order ID</div>
            <div class="text-white text-xs">{{ transaction.orderId }}</div>
          </div>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mt-3">
          <div>
            <div class="text-white/60 text-xs mb-1">Buy Price</div>
            <div class="text-white text-sm font-medium">
              {{ formatCurrency(transaction.buyPrice) }}
            </div>
          </div>
          <div class="sm:text-right">
            <div class="text-white/60 text-xs mb-1">Sell Price</div>
            <div class="text-white text-sm font-medium">
              {{ formatCurrency(transaction.sellPrice) }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Load More Button -->
    <div v-if="canLoadMore" class="mt-6">
      <button
        @click="$emit('loadMore')"
        class="w-full py-3 text-white/80 text-sm font-medium hover:text-white active:text-white/60 transition-colors touch-manipulation"
      >
        View Older Transactions
      </button>
    </div>

    <!-- Error State -->
    <div v-if="error" class="text-center py-8">
      <div class="text-red-400 text-sm mb-4">{{ error.message }}</div>
      <button
        @click="$emit('loadMore')"
        class="px-4 py-2 bg-[#2a3441] text-white rounded-lg hover:bg-[#3a4451] active:bg-[#3a4451]/70 transition-colors touch-manipulation"
      >
        Retry
      </button>
    </div>

    <!-- Empty State -->
    <div v-else-if="transactions.length === 0" class="text-center py-8">
      <div class="text-white/60 text-sm">No transactions found</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { TransactionUI } from '~/types/tradeHistory'

interface Props {
  transactions: TransactionUI[]
  canLoadMore: boolean
  error?: Error | null
}

interface Emits {
  loadMore: []
}

defineProps<Props>()
defineEmits<Emits>()

/**
 * Get cryptocurrency color
 */
const getCryptoColor = (crypto: string): string => {
  const colors: Record<string, string> = {
    'Bitcoin': '#f7931a',
    'Ethereum': '#627eea',
    'XRP': '#23292f'
  }
  return colors[crypto] || '#6366f1'
}

/**
 * Get cryptocurrency symbol
 */
const getCryptoSymbol = (crypto: string): string => {
  const symbols: Record<string, string> = {
    'Bitcoin': '₿',
    'Ethereum': 'Ξ',
    'XRP': 'X'
  }
  return symbols[crypto] || '?'
}

/**
 * Format currency with thousand separators and 1 decimal place
 */
const formatCurrency = (value: number): string => {
  return `¥${value.toLocaleString('ja-JP', { 
    minimumFractionDigits: 0,
    maximumFractionDigits: 1 
  })}`
}

/**
 * Format amount with appropriate decimal places
 */
const formatAmount = (amount: number, crypto: string): string => {
  const symbols: Record<string, string> = {
    'Bitcoin': 'BTC',
    'Ethereum': 'ETH',
    'XRP': 'XRP'
  }
  
  const symbol = symbols[crypto] || crypto.slice(0, 3).toUpperCase()
  
  // Different decimal places for different cryptocurrencies
  const decimals = crypto === 'Bitcoin' ? 3 : crypto === 'Ethereum' ? 2 : 0
  
  return `${amount.toFixed(decimals)} ${symbol}`
}

/**
 * Format timestamp to user-friendly format
 */
const formatTimestamp = (timestamp: Date): string => {
  const now = new Date()
  // Calculate days difference by comparing dates only (ignore time)
  const nowDate = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  const timestampDate = new Date(timestamp.getFullYear(), timestamp.getMonth(), timestamp.getDate())
  const diffMs = nowDate.getTime() - timestampDate.getTime()
  const days = Math.floor(diffMs / (1000 * 60 * 60 * 24))
  
  if (days === 0) {
    return `Today, ${timestamp.toLocaleTimeString('ja-JP', { 
      hour: '2-digit', 
      minute: '2-digit',
      hour12: false 
    })}`
  } else if (days === 1) {
    return `Yesterday, ${timestamp.toLocaleTimeString('ja-JP', { 
      hour: '2-digit', 
      minute: '2-digit',
      hour12: false 
    })}`
  } else {
    return timestamp.toLocaleDateString('ja-JP', { 
      month: 'short', 
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      hour12: false
    }).replace(',', ', ')
  }
}
</script>