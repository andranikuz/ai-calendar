# Smart Goal Calendar - Архитектура системы

## Текущее состояние

**Статус:** Инициализация проекта  
**Дата:** 2025-07-27  
**Версия:** 0.1.1 (MVP + Webhooks)

## Общая архитектура

Smart Goal Calendar использует **Clean Architecture** с четким разделением слоев и зависимостей. Система построена как монолит с возможностью последующего разделения на микросервисы.

### Диаграмма архитектуры

```
┌─────────────────────────────────────────────────────────┐
│                    Presentation Layer                    │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────┐ │
│  │   Web Client    │  │  Mobile Client  │  │    API   │ │
│  │  (React + TS)   │  │ (React Native)  │  │   Docs   │ │
│  └─────────────────┘  └─────────────────┘  └──────────┘ │
└─────────────────────────────────────────────────────────┘
                               │
┌─────────────────────────────────────────────────────────┐
│                     Interface Layer                      │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────┐ │
│  │   HTTP/REST     │  │    Webhooks    │  │   gRPC   │ │
│  │   Endpoints     │  │  (Google Cal)   │  │ Internal │ │
│  └─────────────────┘  └─────────────────┘  └──────────┘ │
└─────────────────────────────────────────────────────────┘
                               │
┌─────────────────────────────────────────────────────────┐
│                   Application Layer                      │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────┐ │
│  │    Commands     │  │     Queries     │  │ Handlers │ │
│  │   (CQRS Write)  │  │   (CQRS Read)   │  │ (UseCases)│ │
│  └─────────────────┘  └─────────────────┘  └──────────┘ │
└─────────────────────────────────────────────────────────┘
                               │
┌─────────────────────────────────────────────────────────┐
│                     Domain Layer                         │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────┐ │
│  │    Entities     │  │ Value Objects   │  │ Services │ │
│  │ (Core Models)   │  │   (Immutable)   │  │(Business)│ │
│  └─────────────────┘  └─────────────────┘  └──────────┘ │
└─────────────────────────────────────────────────────────┘
                               │
┌─────────────────────────────────────────────────────────┐
│                  Infrastructure Layer                    │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────┐ │
│  │  Repositories   │  │   Integrations  │  │   Queue  │ │
│  │  (PostgreSQL)   │  │  (Google, etc)  │  │(Temporal)│ │
│  └─────────────────┘  └─────────────────┘  └──────────┘ │
└─────────────────────────────────────────────────────────┘
```

## Структура проекта

```
smart-goal-calendar/
├── docs/                    # Документация
│   ├── PROJECT.MD           # Техническое задание
│   ├── development-plan.md  # План разработки
│   └── architecture.md      # Текущий документ
├── cmd/                     # Точки входа приложения
│   ├── api/                 # HTTP API сервер
│   ├── worker/              # Background workers
│   └── migrate/             # Миграции БД
├── internal/                # Приватный код приложения
│   ├── domain/              # Доменный слой
│   │   ├── entities/        # Основные сущности
│   │   ├── valueobjects/    # Объекты-значения
│   │   ├── services/        # Доменные сервисы
│   │   └── repositories/    # Интерфейсы репозиториев
│   ├── application/         # Слой приложения
│   │   ├── commands/        # Команды (CQRS)
│   │   ├── queries/         # Запросы (CQRS)
│   │   └── handlers/        # Обработчики
│   ├── adapters/            # Адаптеры
│   │   ├── postgres/        # PostgreSQL реализация
│   │   ├── redis/           # Redis кеширование
│   │   ├── google/          # Google APIs + Webhooks
│   │   └── temporal/        # Temporal workflows
│   └── ports/               # Внешние интерфейсы
│       ├── http/            # HTTP handlers
│       │   ├── handlers/    # Including webhook handlers
│       │   └── routes/      # Including webhook routes
│       ├── grpc/            # gRPC сервисы
│       └── websocket/       # WebSocket handlers
├── web/                     # Frontend приложение
│   ├── src/                 # React TypeScript код
│   ├── public/              # Статические файлы
│   └── package.json         # Frontend зависимости
├── migrations/              # SQL миграции
├── deployments/             # Docker, K8s конфигурации
├── scripts/                 # Утилиты и скрипты
├── .gitignore              # Git ignore правила
├── go.mod                  # Go модуль
├── go.sum                  # Go зависимости
├── Dockerfile              # Docker образ
└── docker-compose.yml      # Локальная разработка
```

## Доменные сущности

### Core Entities

#### User
```go
type User struct {
    ID        UserID
    Email     Email
    Profile   UserProfile
    Settings  UserSettings
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

#### Goal
```go
type Goal struct {
    ID          GoalID
    UserID      UserID
    Title       string
    Description string
    Category    GoalCategory
    Priority    Priority
    Status      GoalStatus
    Progress    Progress
    Deadline    time.Time
    Milestones  []Milestone
    Tasks       []Task
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

#### Event
```go
type Event struct {
    ID            EventID
    UserID        UserID
    GoalID        *GoalID // Optional связь с целью
    Title         string
    Description   string
    StartTime     time.Time
    EndTime       time.Time
    Timezone      Timezone
    Recurrence    *RecurrenceRule
    Location      *Location
    Attendees     []Attendee
    Status        EventStatus
    GoogleEventID *string // Google Calendar Event ID для webhook sync
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

#### GoogleCalendarSync
```go
type GoogleCalendarSync struct {
    ID                  string
    UserID              UserID
    GoogleIntegrationID GoogleIntegrationID
    CalendarID          string
    CalendarName        string
    SyncDirection       CalendarSyncDirection
    SyncStatus          CalendarSyncStatus
    LastSyncAt          time.Time
    WebhookChannelID    string     // Webhook канал для real-time
    WebhookURL          string     // URL для webhook endpoint
    WebhookResourceID   string     // Resource ID от Google
    WebhookExpiresAt    *time.Time // Время истечения webhook
    Settings            CalendarSyncSettings
    CreatedAt           time.Time
    UpdatedAt           time.Time
}
```

#### Mood
```go
type Mood struct {
    ID          MoodID
    UserID      UserID
    Date        Date
    Level       MoodLevel // 1-5 scale
    Notes       string
    Tags        []MoodTag
    RecordedAt  time.Time
}
```

### Value Objects

#### Priority
```go
type Priority int

const (
    PriorityLow Priority = iota + 1
    PriorityMedium
    PriorityHigh
    PriorityCritical
)
```

#### RecurrenceRule
```go
type RecurrenceRule struct {
    Frequency Frequency // DAILY, WEEKLY, MONTHLY, YEARLY
    Interval  int       // Every N frequency
    Until     *time.Time
    Count     *int
    ByDay     []Weekday
    ByMonth   []Month
}
```

## Технологический стек

### Backend
- **Язык:** Go 1.21+
- **Web Framework:** Gin (HTTP routing и middleware)
- **База данных:** PostgreSQL 15+ с JSON полями
- **Кеширование:** Redis 7+ для сессий и кеширования
- **Очереди:** Temporal.io для background jobs
- **Миграции:** golang-migrate/migrate
- **Логирование:** zerolog
- **Конфигурация:** viper
- **Тестирование:** testify, gomock

### Frontend
- **Язык:** TypeScript 5+
- **Framework:** React 18 с хуками
- **State Management:** Redux Toolkit
- **UI Library:** Ant Design
- **Calendar:** FullCalendar
- **HTTP Client:** Axios
- **Build Tool:** Vite
- **Testing:** Jest, React Testing Library

### Infrastructure
- **Контейнеризация:** Docker с multi-stage builds
- **Оркестрация:** docker-compose (dev), Kubernetes (prod)
- **CI/CD:** GitHub Actions
- **Мониторинг:** Prometheus + Grafana
- **Tracing:** Jaeger
- **Cloud:** AWS или Google Cloud Platform

## Интеграции

### Текущие (MVP)
- **Google Calendar API:** OAuth2, двусторонняя синхронизация
- **Google Calendar Webhooks:** Real-time push уведомления об изменениях
- **Google OAuth:** Аутентификация пользователей

### Планируемые (Post-MVP)
- **Microsoft Graph:** Outlook Calendar
- **Notion API:** Database sync
- **CalDAV:** iCloud, Nextcloud
- **Task Management:** Todoist, Asana, Trello

## Webhook Architecture

### Google Calendar Webhooks

Система использует Google Calendar Push Notifications для получения real-time обновлений:

```
┌─────────────────┐      ┌─────────────────┐      ┌─────────────────┐
│  Google Calendar│──────>│ Webhook Handler │──────>│   Event Sync    │
│   Push Service  │ POST  │ /api/v1/google/ │      │   Processor     │
│                 │       │    /webhook     │      │                 │
└─────────────────┘      └─────────────────┘      └─────────────────┘
                                  │                         │
                                  ▼                         ▼
                         ┌─────────────────┐      ┌─────────────────┐
                         │ Channel Manager │      │  Event Storage  │
                         │  (PostgreSQL)   │      │  (PostgreSQL)   │
                         └─────────────────┘      └─────────────────┘
```

### Webhook Flow

1. **Setup Phase:**
   - User authorizes calendar access
   - Application registers webhook with Google
   - Channel ID and resource ID stored in database

2. **Notification Phase:**
   - Google sends POST request on calendar changes
   - Handler validates webhook headers
   - Asynchronous processing of notification

3. **Sync Phase:**
   - Incremental sync fetches only changed events
   - Updates local database with changes
   - Handles creates, updates, and deletes

4. **Management:**
   - Webhooks expire after ~1 week
   - Automatic renewal before expiration (planned)
   - Channel cleanup on user disconnect

## Безопасность

### Аутентификация
- OAuth2 с PKCE для внешних провайдеров
- JWT токены с refresh rotation
- Rate limiting по IP и пользователю

### Авторизация
- RBAC с ролями (User, Admin)
- Resource-based permissions
- Multi-tenant изоляция данных

### Защита данных
- Шифрование sensitive данных at rest (AES-256)
- TLS 1.3 для всех соединений
- Аудит логи для критических операций
- GDPR compliance для mood данных

## Производительность

### Оптимизации
- Database индексы для частых запросов
- Redis кеширование для session data
- Connection pooling для PostgreSQL
- Lazy loading для больших datasets

### Масштабируемость
- Горизонтальное масштабирование API серверов
- Read replicas для PostgreSQL
- CDN для статических файлов
- Партиционирование таблиц по времени

## Мониторинг

### Метрики
- Application metrics (Prometheus)
- Database performance (pg_stat_statements)
- HTTP request metrics (response time, errors)
- Business metrics (MAU, DAU, conversion)

### Логирование
- Structured logging (JSON format)
- Centralized log aggregation
- Error tracking (Sentry интеграция)
- Audit trails для важных операций

## Статус реализации

### ✅ Выполнено
- [x] Документация проекта
- [x] План архитектуры
- [x] Техническое задание
- [x] Инициализация Go проекта
- [x] Структура каталогов
- [x] Docker setup
- [x] Google Calendar OAuth2 интеграция
- [x] Google Calendar Webhook интеграция

### 🔄 В процессе
- [ ] Базовые тесты репозиториев
- [ ] Автоматическое продление webhook подписок

### 📋 Запланировано
- [ ] Базовые тесты репозиториев
- [ ] Дополнительные PWA возможности
- [ ] Улучшение аналитики и дашбордов

## Следующие шаги
1. **Базовые unit-тесты репозиториев**
2. **Расширение PWA функциональности и offline sync**
3. **Продолжение работы над аналитикой и отчетами**