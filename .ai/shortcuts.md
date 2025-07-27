# ⚡ Smart Goal Calendar - Shortcuts & Aliases

## 🚀 Основные алиасы команд

### Односимвольные алиасы
```
s    → status      # Быстрая проверка состояния
d    → develop     # Начать разработку
r    → reflect     # Глубокий анализ
p    → plan        # Планирование задач
i    → init        # Инициализация
t    → test        # Запуск тестов
l    → lint        # Проверка качества
b    → build       # Сборка проекта
```

### Быстрые комбинации
```
qf   → quick-fix   # status → develop → test → commit
fc   → full-check  # lint + test + audit + docs-sync
```

---

## 🔧 Технические shortcuts

### Backend (Go) команды
```bash
# Алиас: go-test
go test ./... && go vet ./...

# Алиас: go-check  
go test ./... && go vet ./... && go build ./cmd/api

# Алиас: go-clean
go mod tidy && go fmt ./...

# Алиас: go-race
go test -race ./...
```

### Frontend (React) команды
```bash
# Алиас: react-dev
cd web && npm run dev

# Алиас: react-check
cd web && npm run lint && npm run typecheck && npm run build

# Алиас: react-fix
cd web && npm run lint:fix && npm run typecheck

# Алиас: react-clean
cd web && rm -rf node_modules && npm install
```

---

## ⚡ Composite shortcuts

### Быстрые workflow
```
my-flow        = status + develop + test
release-check  = lint + test + audit + docs-sync
hotfix-flow    = status + bug-fix + test + commit
feature-flow   = plan + develop + test + docs-sync
```

### Quality assurance shortcuts
```
qa-backend     = go test ./... + go vet ./... + golangci-lint run
qa-frontend    = cd web && npm run lint + npm run typecheck + npm run build
qa-full        = qa-backend + qa-frontend + docs-sync
```

### Development shortcuts
```
dev-start      = status + init (если нужно)
dev-continue   = status + develop
dev-finish     = test + docs-sync + commit
dev-cycle      = status + develop + test + commit
```

---

## 🎯 Smart Goal Calendar specific

### Google Calendar shortcuts
```
gcal-test      = calendar-sync + webhook-test
gcal-debug     = Check OAuth2 flow + webhook logs
gcal-reset     = Clear tokens + re-authorize
```

### PWA shortcuts  
```
pwa-check      = Service Worker + Manifest + Offline test
pwa-build      = Build + PWA validation + lighthouse
pwa-deploy     = Build + PWA check + deployment ready
```

### Goals system shortcuts
```
goals-test     = SMART validation + progress tracking + task management
goals-demo     = Create demo goals + test workflows
goals-export   = Export goals data + validation
```

### Mood tracking shortcuts
```
mood-test      = 5-level validation + analytics + privacy check
mood-demo      = Create demo mood data + test correlations
mood-export    = Export mood data + anonymization
```

---

## 📋 Documentation shortcuts

### Docs maintenance
```
docs-check     = Validate all docs + check links + update dates
docs-sync      = Compare docs with code + update inconsistencies  
docs-full      = docs-check + docs-sync + changelog update
```

### Architecture documentation
```
arch-update    = Update architecture.md + diagrams + dependencies
arch-validate  = Check Clean Architecture compliance + layers
arch-export    = Generate architecture diagrams + documentation
```

---

## 🐛 Bug management shortcuts

### Bug workflow
```
bug-report     = Interactive bug creation with templates
bug-triage     = Analyze + prioritize + assign bugs
bug-sprint     = List high-priority bugs + fix recommendations
```

### Debug shortcuts
```
debug-backend  = Go delve + logs + database state
debug-frontend = React DevTools + Redux DevTools + Network
debug-integration = API logs + Google Calendar sync + webhook status
```

---

## 🚀 Deployment shortcuts

### Build and validation
```
pre-deploy     = full-check + security-audit + performance-test
deploy-ready   = pre-deploy + docs-sync + changelog-update
post-deploy    = health-check + integration-test + monitoring-setup
```

### Environment shortcuts
```
env-dev        = Start development environment + database + services
env-test       = Setup test environment + mock services + test data
env-prod       = Production health check + monitoring + alerts
```

---

## 🎨 Custom user shortcuts

### Создание собственных алиасов
Добавьте в `.ai-config.json`:

```json
{
  "commands": {
    "aliases": {
      "my-command": "status + develop + custom-action"
    },
    "shortcuts": {
      "weekend-deploy": "full-check + deploy-ready + notify-team",
      "morning-routine": "status + plan + priority-review",
      "code-review": "lint + test + docs-check + security-scan"
    }
  }
}
```

### Conditional shortcuts
```
# Если есть критические баги
critical-mode  = bug-list + bug-fix + immediate-test + hotfix-deploy

# Если close to deadline
crunch-mode    = status + highest-priority-only + minimal-testing

# Если новый feature
feature-mode   = plan + design-review + develop + comprehensive-test
```

---

## 📱 Mobile/PWA specific shortcuts

### Mobile testing
```
mobile-test    = Responsive test + Touch interactions + PWA install
mobile-debug   = Mobile DevTools + Performance + Accessibility
mobile-deploy  = Mobile build + App store ready + PWA manifest
```

### Performance shortcuts
```
perf-audit     = Lighthouse + Bundle analyzer + Core Web Vitals
perf-optimize  = Bundle optimization + Image optimization + Caching
perf-monitor   = Performance monitoring + Alerts + Reporting
```

---

## 🔍 Search and analysis shortcuts

### Code analysis
```
code-search    = Find patterns + Dependencies + Unused code
code-metrics   = Lines of code + Complexity + Technical debt
code-security  = Security scan + Vulnerability check + Compliance
```

### Git shortcuts
```
git-clean      = Clean branches + Unused files + Optimize repo
git-history    = Commit analysis + Contributors + Change patterns
git-release    = Tag + Release notes + Deploy preparation
```

---

## ⌨️ Keyboard shortcuts for CLI

### Quick access patterns
```
Ctrl+S    → Equivalent to 's' (status)
Ctrl+D    → Equivalent to 'd' (develop)  
Ctrl+T    → Equivalent to 't' (test)
Ctrl+B    → Equivalent to 'b' (build)
```

### Command completion
```
Tab       → Auto-complete available commands
Shift+Tab → Show command help
Ctrl+R    → Recent commands history
Ctrl+L    → Clear screen, show available shortcuts
```

---

## 📊 Reporting shortcuts

### Status reports
```
daily-report   = status + progress-summary + next-actions
weekly-report  = reflect + accomplishments + roadmap-update  
monthly-report = comprehensive-analysis + metrics + planning
```

### Metrics shortcuts
```
health-check   = Technical metrics + Performance + Quality scores
progress-track = Goal completion + Feature progress + Timeline
quality-report = Test coverage + Code quality + Security status
```

---

## 🎯 Usage examples

### Утренний старт
```bash
# Быстрая проверка и начало работы
> s                    # status
> d                    # develop (если есть готовые задачи)
# или
> qf                   # quick-fix для небольших задач
```

### Подготовка к релизу
```bash
# Полная проверка качества
> fc                   # full-check
> release-check        # comprehensive validation
> deploy-ready         # final preparation
```

### Работа с багами
```bash
# Анализ и исправление
> bug-list             # посмотреть активные баги  
> bug-fix              # выбрать и исправить
> qa-full              # полная проверка качества
```

### Еженедельная рефлексия
```bash
# Глубокий анализ
> r                    # reflect
> weekly-report        # сформировать отчет
> p                    # plan (обновить roadmap)
```

## 🔗 Configuration

Все shortcuts настраиваются через:
- `.ai-config.json` - основная конфигурация
- `.ai/shortcuts.md` - документация (этот файл)
- Локальные alias в shell (опционально)

**Принцип**: Чем чаще команда используется, тем короче должен быть алиас.