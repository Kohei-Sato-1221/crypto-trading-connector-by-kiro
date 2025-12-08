<script setup lang="ts">
interface NavigationItem {
  name: string
  path: string
  icon: string
  label: string
}

const navItems: NavigationItem[] = [
  { name: 'market', path: '/market', icon: '📊', label: 'Market' },
  { name: 'trade', path: '/trade', icon: '💱', label: 'Trade' },
  { name: 'history', path: '/history', icon: '📜', label: 'History' },
  { name: 'portfolio', path: '/portfolio', icon: '💼', label: 'Portfolio' }
]

// Get current route - use composable in production, allow override for testing
let currentPath: string
try {
  const route = useRoute()
  currentPath = route.path
} catch {
  // Fallback for testing environment
  currentPath = '/market'
}

const isActive = (path: string) => {
  return currentPath === path
}
</script>

<template>
  <nav class="fixed bottom-0 left-0 right-0 bg-[#192633] border-t border-[#2e3e50] h-16 z-50">
    <div class="flex items-center justify-around h-full max-w-2xl mx-auto px-2">
      <NuxtLink
        v-for="item in navItems"
        :key="item.name"
        :to="item.path"
        class="flex flex-col items-center justify-center flex-1 h-full min-w-[44px] transition-colors"
        :class="isActive(item.path) ? 'text-[#137fec]' : 'text-slate-400'"
      >
        <span class="text-2xl mb-1">{{ item.icon }}</span>
        <span class="text-[10px] font-medium">{{ item.label }}</span>
      </NuxtLink>
    </div>
  </nav>
</template>
