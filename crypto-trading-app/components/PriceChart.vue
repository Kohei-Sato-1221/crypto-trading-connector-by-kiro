<script setup lang="ts">
import { computed } from 'vue'
import { Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Filler,
  type ChartOptions
} from 'chart.js'
import type { ChartDataPoint } from '~/types/crypto'

// Register Chart.js components
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Filler
)

interface Props {
  data: ChartDataPoint[]
  isPositive: boolean
  currency: string
}

const props = defineProps<Props>()

// Chart data
const chartData = computed(() => ({
  labels: props.data.map(d => d.day),
  datasets: [
    {
      data: props.data.map(d => d.price),
      borderColor: props.isPositive ? '#0bda5b' : '#fa6238',
      backgroundColor: props.isPositive
        ? 'rgba(11, 218, 91, 0.1)'
        : 'rgba(250, 98, 56, 0.1)',
      borderWidth: 2,
      fill: true,
      tension: 0.4,
      pointRadius: 0,
      pointHoverRadius: 4,
      pointHoverBackgroundColor: props.isPositive ? '#0bda5b' : '#fa6238',
      pointHoverBorderColor: '#fff',
      pointHoverBorderWidth: 2
    }
  ]
}))

// Chart options
const chartOptions = computed<ChartOptions<'line'>>(() => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      display: false
    },
    tooltip: {
      enabled: true,
      mode: 'index',
      intersect: false,
      backgroundColor: 'rgba(28, 38, 48, 0.95)',
      titleColor: '#fff',
      bodyColor: '#94a3b8',
      borderColor: '#2e3e50',
      borderWidth: 1,
      padding: 12,
      displayColors: false,
      callbacks: {
        label: (context) => {
          const value = context.parsed.y
          return `${props.currency}: ¥${value.toLocaleString()}`
        }
      }
    }
  },
  scales: {
    x: {
      grid: {
        display: false
      },
      ticks: {
        color: '#94a3b8',
        font: {
          size: 10,
          family: 'Manrope'
        }
      },
      border: {
        display: false
      }
    },
    y: {
      position: 'right',
      grid: {
        color: 'rgba(46, 62, 80, 0.5)',
        drawBorder: false
      },
      ticks: {
        color: '#94a3b8',
        font: {
          size: 10,
          family: 'Manrope'
        },
        callback: (value) => {
          const num = Number(value)
          if (num >= 1000000) {
            return `${(num / 1000000).toFixed(1)}M`
          } else if (num >= 1000) {
            return `${(num / 1000).toFixed(0)}K`
          }
          return num.toString()
        }
      },
      border: {
        display: false,
        dash: [5, 5]
      }
    }
  },
  interaction: {
    mode: 'index',
    intersect: false
  }
}))
</script>

<template>
  <div class="w-full h-[200px]">
    <Line :data="chartData" :options="chartOptions" />
  </div>
</template>
