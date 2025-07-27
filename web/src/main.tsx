import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
// Import Ant Design styles once
import 'antd/dist/reset.css'
import App from './App.tsx'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
