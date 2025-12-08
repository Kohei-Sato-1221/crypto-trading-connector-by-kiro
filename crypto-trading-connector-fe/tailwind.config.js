/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./components/**/*.{js,vue,ts}",
    "./layouts/**/*.vue",
    "./pages/**/*.vue",
    "./plugins/**/*.{js,ts}",
    "./app/**/*.vue",
  ],
  theme: {
    extend: {
      colors: {
        'bg': {
          'primary': '#101922',
          'secondary': '#1c2630',
          'tertiary': '#2e3e50',
        },
        'accent': {
          'primary': '#137fec',
          'success': '#0bda5b',
          'danger': '#fa6238',
        },
        'bitcoin': {
          'orange': '#f7931a',
        },
        'ethereum': {
          'purple': '#627eea',
        },
      },
      fontFamily: {
        sans: ['Manrope', 'sans-serif'],
      },
    },
  },
  plugins: [],
}
