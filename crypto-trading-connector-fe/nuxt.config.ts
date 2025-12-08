// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  
  modules: ['@nuxtjs/tailwindcss'],
  
  css: ['~/assets/css/main.css'],
  
  app: {
    head: {
      title: 'Crypto Trading Connector',
      link: [
        {
          rel: 'stylesheet',
          href: 'https://fonts.googleapis.com/css2?family=Manrope:wght@400;500;600;700;800&display=swap'
        }
      ]
    }
  },

  runtimeConfig: {
    public: {
      // Use mock data by default in development, set to false to use API
      useMockData: process.env.USE_MOCK_DATA === 'true',
      // Backend API base URL
      apiBaseUrl: process.env.API_BASE_URL || 'http://localhost:8080'
    }
  }
})
