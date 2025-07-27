# 🛠️ Smart Goal Calendar - Технические стандарты

## 🏗️ Архитектура проекта

### Backend (Go)
**Архитектурный стиль**: Clean Architecture
**Принципы**: SOLID, DRY, KISS

```
internal/
├── domain/              # Бизнес-логика (независимый слой)
│   ├── entities/        # Основные сущности
│   ├── valueobjects/    # Неизменяемые объекты-значения
│   ├── services/        # Доменные сервисы
│   └── repositories/    # Интерфейсы репозиториев
├── application/         # Use Cases (CQRS)
│   ├── commands/        # Команды изменения состояния
│   ├── queries/         # Запросы чтения данных
│   └── handlers/        # Обработчики команд/запросов
├── adapters/            # Внешние адаптеры
│   ├── postgres/        # PostgreSQL реализация
│   ├── google/          # Google APIs интеграция
│   └── auth/            # JWT аутентификация
└── ports/               # Внешние интерфейсы
    └── http/            # REST API endpoints
```

### Frontend (React + TypeScript)
**Архитектурный стиль**: Feature-based + Layered
**Принципы**: Component composition, Hooks pattern, PWA-first

```
src/
├── components/          # Переиспользуемые компоненты
│   ├── Calendar/        # Календарные компоненты
│   ├── Goals/           # Компоненты целей
│   ├── Common/          # Общие UI компоненты
│   └── Layout/          # Макет и навигация
├── pages/               # Страницы приложения
├── store/               # Redux Toolkit слайсы
├── services/            # API клиенты
├── hooks/               # Кастомные React хуки
├── utils/               # Вспомогательные функции
└── types/               # TypeScript определения
```

---

## 🔧 Технологический стек

### Backend Dependencies
```go
// Core
"github.com/gin-gonic/gin"           // HTTP router
"gorm.io/gorm"                       // ORM для PostgreSQL
"github.com/golang-jwt/jwt/v5"       // JWT аутентификация

// Google APIs
"google.golang.org/api/calendar/v3"  // Google Calendar API
"golang.org/x/oauth2"                // OAuth2 flow

// Database
"gorm.io/driver/postgres"            // PostgreSQL driver
"github.com/golang-migrate/migrate"  // Database migrations

// Testing
"github.com/stretchr/testify"        // Test framework
"github.com/DATA-DOG/go-sqlmock"     // SQL mocking
```

### Frontend Dependencies  
```json
{
  "react": "^18.2.0",
  "typescript": "^5.0.0",
  "@reduxjs/toolkit": "^1.9.0",
  "react-redux": "^8.0.0",
  "antd": "^5.0.0",
  "@mui/material": "^5.0.0",
  "@fullcalendar/react": "^6.0.0",
  "axios": "^1.0.0",
  "dayjs": "^1.11.0"
}
```

---

## 📏 Стандарты кодирования

### Go Code Style
```go
// ✅ Правильно: Четкие интерфейсы
type GoalRepository interface {
    Create(ctx context.Context, goal *Goal) error
    GetByID(ctx context.Context, id string) (*Goal, error)
    Update(ctx context.Context, goal *Goal) error
    Delete(ctx context.Context, id string) error
}

// ✅ Правильно: Структурированные ошибки
var (
    ErrGoalNotFound = errors.New("goal not found")
    ErrInvalidGoal  = errors.New("invalid goal data")
)

// ✅ Правильно: Контекст для всех операций
func (r *goalRepository) Create(ctx context.Context, goal *Goal) error {
    // Implementation
}
```

### TypeScript Code Style
```typescript
// ✅ Правильно: Строгая типизация
interface Goal {
  id: string;
  title: string;
  description: string;
  category: GoalCategory;
  priority: Priority;
  progress: number;
  deadline: string;
  createdAt: string;
  updatedAt: string;
}

// ✅ Правильно: Redux Toolkit slices
export const goalsSlice = createSlice({
  name: 'goals',
  initialState,
  reducers: {
    setGoals: (state, action: PayloadAction<Goal[]>) => {
      state.goals = action.payload;
    },
  },
});

// ❌ Неправильно: Использование any
const data: any = response.data; // Избегать!

// ✅ Правильно: Typed API responses
interface ApiResponse<T> {
  data: T;
  message: string;
  success: boolean;
}
```

---

## 🧪 Тестирование

### Backend Testing Strategy
```go
// Unit tests - domain layer
func TestGoal_Validate(t *testing.T) {
    tests := []struct {
        name    string
        goal    Goal
        wantErr bool
    }{
        {
            name: "valid goal",
            goal: Goal{Title: "Learn Go", Category: "education"},
            wantErr: false,
        },
    }
    // Implementation
}

// Integration tests - repository layer
func TestGoalRepository_Create(t *testing.T) {
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    repo := NewGoalRepository(db)
    goal := &Goal{Title: "Test Goal"}
    
    err := repo.Create(context.Background(), goal)
    assert.NoError(t, err)
    assert.NotEmpty(t, goal.ID)
}
```

### Frontend Testing Strategy
```typescript
// Component testing
import { render, screen } from '@testing-library/react';
import { GoalModal } from './GoalModal';

test('renders goal creation form', () => {
  render(<GoalModal visible={true} onClose={() => {}} />);
  expect(screen.getByText('Создать цель')).toBeInTheDocument();
});

// Redux testing
import { store } from '../store';
import { setGoals } from '../store/slices/goalsSlice';

test('should update goals state', () => {
  const goals = [{ id: '1', title: 'Test Goal' }];
  store.dispatch(setGoals(goals));
  expect(store.getState().goals.goals).toEqual(goals);
});
```

---

## 🔐 Безопасность

### Authentication & Authorization
```go
// JWT middleware
func JWTMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := extractToken(c.GetHeader("Authorization"))
        claims, err := validateJWT(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "unauthorized"})
            c.Abort()
            return
        }
        c.Set("userID", claims.UserID)
        c.Next()
    }
}

// Input validation
type CreateGoalRequest struct {
    Title       string `json:"title" binding:"required,min=1,max=200"`
    Description string `json:"description" binding:"max=1000"`
    Category    string `json:"category" binding:"required,oneof=health career education"`
}
```

### Frontend Security
```typescript
// API interceptors для автоматического refresh токенов
axios.interceptors.response.use(
  response => response,
  async error => {
    if (error.response?.status === 401) {
      await refreshAuthToken();
      return axios.request(error.config);
    }
    return Promise.reject(error);
  }
);

// XSS защита - санитизация входных данных
const sanitizeInput = (input: string): string => {
  return DOMPurify.sanitize(input);
};
```

---

## 📊 Производительность

### Backend Optimization
```go
// Database indexing
CREATE INDEX idx_goals_user_id ON goals(user_id);
CREATE INDEX idx_events_user_id_date ON events(user_id, start_time);

// Connection pooling
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    ConnPool: &sql.DB{
        MaxOpenConns: 25,
        MaxIdleConns: 5,
        ConnMaxLifetime: time.Hour,
    },
})

// Caching strategy
func (s *goalService) GetGoals(ctx context.Context, userID string) ([]Goal, error) {
    cacheKey := fmt.Sprintf("goals:%s", userID)
    
    // Try cache first
    if cached := s.cache.Get(cacheKey); cached != nil {
        return cached.([]Goal), nil
    }
    
    // Fallback to database
    goals, err := s.repo.GetByUserID(ctx, userID)
    if err == nil {
        s.cache.Set(cacheKey, goals, 5*time.Minute)
    }
    return goals, err
}
```

### Frontend Optimization
```typescript
// Lazy loading компонентов
const GoalsPage = lazy(() => import('./pages/GoalsPage'));
const CalendarPage = lazy(() => import('./pages/CalendarPage'));

// Memoization для дорогих вычислений
const MemoizedGoalList = memo(({ goals }: { goals: Goal[] }) => {
  const sortedGoals = useMemo(() => {
    return goals.sort((a, b) => a.priority - b.priority);
  }, [goals]);
  
  return <GoalList goals={sortedGoals} />;
});

// Virtual scrolling для больших списков
import { FixedSizeList as List } from 'react-window';

const GoalVirtualList = ({ goals }: { goals: Goal[] }) => (
  <List
    height={600}
    itemCount={goals.length}
    itemSize={80}
    itemData={goals}
  >
    {GoalRow}
  </List>
);
```

---

## 🚀 Deployment & DevOps

### Docker Configuration
```dockerfile
# Multi-stage build для Go
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/main .
CMD ["./main"]
```

### Environment Configuration
```yaml
# config/config.yaml
server:
  port: 8080
  timeout: 30s

database:
  host: localhost
  port: 5432
  name: smart_calendar
  max_connections: 25

google:
  client_id: ${GOOGLE_CLIENT_ID}
  client_secret: ${GOOGLE_CLIENT_SECRET}

redis:
  url: redis://localhost:6379
  ttl: 300s
```

---

## 📋 Обязательные проверки

### Pre-commit Checks
```bash
# Backend
go fmt ./...
go vet ./...
go test ./...
golangci-lint run

# Frontend  
cd web
npm run lint
npm run typecheck
npm run test
npm run build
```

### CI/CD Pipeline
```yaml
# .github/workflows/ci.yml
name: CI/CD
on: [push, pull_request]

jobs:
  backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: 1.21
      - run: go test ./...
      - run: go vet ./...
      
  frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: 18
      - run: cd web && npm ci
      - run: cd web && npm run build
```

---

## 🎯 Качественные метрики

### Code Quality Targets
- **Test Coverage**: Backend >80%, Frontend >70%
- **TypeScript Strict**: Без any типов
- **Bundle Size**: <1MB для initial load
- **Performance**: Lighthouse score >90
- **Accessibility**: WCAG 2.1 AA compliance

### Monitoring
```go
// Prometheus metrics
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
)
```

---

## 🔗 Интеграции

### Google Calendar API
```go
// OAuth2 flow
config := &oauth2.Config{
    ClientID:     cfg.Google.ClientID,
    ClientSecret: cfg.Google.ClientSecret,
    RedirectURL:  cfg.Google.RedirectURL,
    Scopes:       []string{calendar.CalendarScope},
    Endpoint:     google.Endpoint,
}

// Webhook handling
func (h *webhookHandler) HandleGoogleCalendarWebhook(c *gin.Context) {
    payload := &GoogleWebhookPayload{}
    if err := c.ShouldBindJSON(payload); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // Process incremental sync
    go h.syncService.ProcessWebhookUpdate(payload)
    c.JSON(200, gin.H{"status": "ok"})
}
```

### PWA Configuration
```typescript
// Service Worker registration
if ('serviceWorker' in navigator && process.env.NODE_ENV === 'production') {
  navigator.serviceWorker.register('/sw.js');
}

// Web App Manifest
{
  "name": "Smart Goal Calendar",
  "short_name": "GoalCal",
  "start_url": "/",
  "display": "standalone",
  "theme_color": "#1890ff",
  "background_color": "#ffffff"
}
```

---

## 📐 Принципы разработки

1. **API First** - сначала API дизайн, потом реализация
2. **Type Safety** - строгая типизация на всех уровнях
3. **Error Handling** - structured errors с контекстом
4. **Documentation** - код должен быть self-documenting
5. **Testing** - TDD подход для критической логики
6. **Performance** - измеряй и оптимизируй
7. **Security** - secure by default
8. **Accessibility** - доступность с самого начала

## 🔗 Связанные файлы
- `.ai/commands.md` - Доступные команды
- `.ai/workflows.md` - Процессы разработки
- `docs/architecture.md` - Архитектурная документация
- `docs/development-plan.md` - План разработки