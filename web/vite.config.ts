import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          // React core
          react: ['react', 'react-dom'],
          
          // Router
          router: ['react-router-dom'],
          
          // State management
          redux: ['@reduxjs/toolkit', 'react-redux'],
          
          // UI libraries - separate chunks for better caching
          antd: ['antd'],
          mui: ['@mui/material', '@mui/icons-material'],
          
          // Calendar library
          calendar: ['@fullcalendar/core', '@fullcalendar/react', '@fullcalendar/daygrid', '@fullcalendar/timegrid', '@fullcalendar/interaction'],
          
          // HTTP client
          axios: ['axios'],
          
          // Utilities
          utils: ['dayjs']
        }
      }
    },
    // Reduce chunk size warning limit to 800kb
    chunkSizeWarningLimit: 800
  }
})
