# 🚀 Quick Start Commands

## Frontend Commands

```bash
# Navigate to frontend
cd web

# Install dependencies
npm i

# Development server
npm run dev

# Production build
npm run build

# Preview production build
npm run preview

# Lint code
npm run lint
```

## Project Structure

```
ai-calendar/
├── web/           # React Frontend
├── docs/          # Documentation
└── README.md      # Main readme
```

## Quick Commands Reference

| Command | What it does |
|---------|-------------|
| `npm i` | Install dependencies |
| `npm run dev` | Start dev server (port 5173) |
| `npm run build` | Build for production |
| `npm run preview` | Preview prod build (port 4173) |
| `npm run lint` | Check code quality |

## Environment Setup

Create `web/.env`:
```env
VITE_API_URL=http://localhost:8080
VITE_GOOGLE_CLIENT_ID=your_client_id
```

## Troubleshooting

### Build Issues
```bash
# Clean install
rm -rf node_modules package-lock.json
npm i
npm run build
```

### Dev Server Issues
```bash
# Use preview instead
npm run build
npm run preview
```

### Port Issues
```bash
# Dev server - custom port
npm run dev -- --port 3000

# Preview - custom port  
npm run preview -- --port 4000
```

## Browser Support

- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

## Mobile Testing

1. Start server: `npm run preview`
2. Find your IP: `ifconfig | grep inet`
3. Access: `http://YOUR_IP:4173`
4. Or use: `npm run preview -- --host`

## PWA Testing

1. Open in Chrome
2. Install prompt will appear
3. Or manually: Chrome menu → Install app

## Production Deployment

```bash
# Build optimized version
npm run build

# Deploy dist/ folder to your hosting
# Files are in web/dist/
```

---

**Need help?** Check [README.md](README.md) for detailed documentation.