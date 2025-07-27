import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import vitePluginImp from 'vite-plugin-imp'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    react(),
    vitePluginImp({
      libList: [
        {
          libName: 'antd',
          style: () => false, // We'll import styles manually where needed
        },
      ],
    }),
  ],
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
          
          // UI libraries
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
