# Smart Goal Calendar - План разработки

## Этап 1: Инфраструктура и основы (1-2 недели)

### 1.1 Инициализация проекта
- [x] ~~Создание структуры каталогов docs/~~
- [x] ~~Создание структуры каталогов согласно Clean Architecture~~
- [x] ~~Настройка go.mod, базовые зависимости (gin/echo, postgres driver)~~
- [x] ~~Создание Dockerfile и docker-compose.yml~~
- [x] ~~Настройка .gitignore и базовой документации~~

### 1.2 База данных и миграции
- [ ] Настройка PostgreSQL через Docker
- [ ] Создание миграций для основных таблиц (users, goals, events, moods)
- [ ] Настройка подключения к БД и connection pool
- [ ] Базовая система миграций

### 1.3 Domain слой
- [x] ~~Создание entities: User, Goal, Event, Mood~~
- [x] ~~Реализация value objects: Priority, Recurrence, Timezone~~
- [x] ~~Базовые domain services для бизнес-логики~~
- [x] ~~Валидация и доменные правила~~

## Этап 2: Backend API (2-3 недели)

### 2.1 Repository слой
- [x] ~~Интерфейсы репозиториев для каждой entity~~
- [x] ~~PostgreSQL имплементация с CRUD операциями~~
- [x] ~~Настройка connection pooling и транзакций~~
- [ ] Базовые тесты для репозиториев

### 2.2 HTTP API
- [x] ~~Настройка роутинга (gin/echo)~~
- [x] ~~Endpoints для пользователей (регистрация, авторизация)~~
- [ ] CRUD endpoints для целей и событий
- [x] ~~Middleware для аутентификации и логирования~~

### 2.3 Аутентификация
- [x] ~~JWT токены для аутентификации~~
- [x] ~~OAuth2 интеграция с Google~~
- [x] ~~Middleware для защиты endpoints~~
- [x] ~~Refresh token механизм~~

## Этап 3: Core функциональность (2-3 недели)

### 3.1 Календарная система
- [ ] CRUD операции для событий
- [ ] Поддержка повторяющихся событий (RRULE)
- [ ] Временные зоны и локализация
- [ ] Drag-and-drop API для перемещения событий

### 3.2 Система целей
- [ ] SMART цели с валидацией
- [ ] Разбивка целей на подзадачи
- [ ] Автоматическое планирование времени для целей
- [ ] Прогресс трекинг

### 3.3 Google Calendar интеграция
- [x] ~~OAuth2 настройка для Google API~~
- [x] ~~Двустороння синхронизация событий~~
- [ ] Webhook для real-time обновлений
- [ ] Обработка конфликтов данных

## Этап 4: Frontend основы (3-4 недели)

### 4.1 React приложение
- [x] ~~Vite React App с TypeScript~~
- [x] ~~Настройка Redux Toolkit для состояния~~
- [x] ~~Роутинг с React Router~~
- [x] ~~Базовые компоненты и стили~~

### 4.2 Календарный интерфейс
- [x] ~~Интеграция FullCalendar~~
- [x] ~~Просмотры: день, неделя, месяц~~
- [x] ~~Drag-and-drop для событий~~
- [x] ~~Quick создание событий~~

### 4.3 Управление целями
- [ ] Компонент создания SMART целей
- [ ] Визуализация прогресса (progress bars)
- [ ] Дерево целей и подзадач
- [ ] Интеграция с календарем

## Этап 5: Mood tracking и аналитика (2 недели)

### 5.1 Mood tracking
- [ ] Простой emoji selector (5 эмоций)
- [ ] Ежедневные напоминания
- [ ] Визуализация данных (heat map)
- [ ] Privacy-first хранение

### 5.2 Базовая аналитика
- [ ] Корреляция настроения и продуктивности
- [ ] Простые дашборды прогресса
- [ ] Экспорт базовых отчетов
- [ ] KPI отслеживание для целей

## Этап 6: Полировка и тестирование (2-3 недели)

### 6.1 Тестирование
- [ ] Unit тесты для критичной бизнес-логики
- [ ] Integration тесты для API endpoints
- [ ] E2E тесты для ключевых user flows
- [ ] Performance тестирование

### 6.2 UI/UX улучшения
- [ ] Responsive дизайн для мобильных устройств
- [ ] Accessibility (WCAG 2.1)
- [ ] Анимации и микроинтеракции
- [ ] Error handling и loading states

### 6.3 Deployment
- [ ] CI/CD pipeline с GitHub Actions
- [ ] Kubernetes конфигурация или простой VPS
- [ ] Мониторинг и логирование
- [ ] Backup и disaster recovery

## Приоритизация задач

**Высокий приоритет (MVP):**
- Базовый календарь с CRUD операциями
- Простая система целей
- Google Calendar синхронизация
- Базовая аутентификация

**Средний приоритет:**
- Mood tracking
- Расширенная аналитика
- UI полировка
- Mobile responsive

**Низкий приоритет (post-MVP):**
- Дополнительные интеграции
- AI-планировщик
- Командные функции
- Advanced аналитика

## Текущий статус

**Выполнено:**
- [x] Создание структуры документации
- [x] План разработки
- [x] Инициализация Go проекта
- [x] Базовая архитектура
- [x] Docker окружение

**В работе:**
- [ ] Goal/Event CRUD endpoints
- [ ] Google Calendar интеграция

**Следующие шаги:**
1. Goal management API
2. Event calendar API  
3. Google Calendar OAuth2
4. Frontend React приложение

## Протестированные API endpoints

### ✅ User Management API
- `POST /api/v1/auth/register` - регистрация пользователя
- `POST /api/v1/auth/login` - аутентификация
- `POST /api/v1/auth/refresh` - обновление токенов
- `GET /api/v1/users/me` - профиль пользователя
- `PUT /api/v1/users/me` - обновление профиля
- `DELETE /api/v1/users/me` - удаление аккаунта

### ✅ Goal Management API (COMPLETED)
- `POST /api/v1/goals` - создание цели
- `GET /api/v1/goals` - получение целей пользователя (с пагинацией)
- `GET /api/v1/goals/:id` - получение конкретной цели
- `PUT /api/v1/goals/:id` - обновление цели
- `DELETE /api/v1/goals/:id` - удаление цели
- `POST /api/v1/goals/:id/tasks` - создание задачи для цели
- `GET /api/v1/goals/:id/tasks` - получение задач цели
- `POST /api/v1/goals/tasks/:taskId/complete` - завершение задачи
- `POST /api/v1/goals/:id/milestones` - создание milestone для цели
- `GET /api/v1/goals/:id/milestones` - получение milestones цели
- `POST /api/v1/goals/milestones/:milestoneId/complete` - завершение milestone

### ✅ Event Calendar API (COMPLETED)
- `POST /api/v1/events` - создание события
- `GET /api/v1/events` - получение событий пользователя (с пагинацией)
- `GET /api/v1/events/:id` - получение конкретного события
- `PUT /api/v1/events/:id` - обновление события
- `DELETE /api/v1/events/:id` - удаление события
- `GET /api/v1/events/search` - поиск событий
- `GET /api/v1/events/upcoming` - получение предстоящих событий
- `GET /api/v1/events/today` - получение событий на сегодня
- `GET /api/v1/events/time-range` - получение событий по временному диапазону
- `GET /api/v1/events/conflict-check` - проверка конфликтов времени
- `POST /api/v1/events/:id/move` - перемещение события
- `POST /api/v1/events/:id/duplicate` - дублирование события
- `POST /api/v1/events/:id/status` - изменение статуса события
- `POST /api/v1/events/:id/link-goal` - связывание события с целью
- `POST /api/v1/events/:id/unlink-goal` - отвязывание события от цели

### ✅ Mood Tracking API (COMPLETED)
- `POST /api/v1/moods` - создание записи настроения
- `GET /api/v1/moods` - получение записей настроения пользователя (с пагинацией)
- `GET /api/v1/moods/:id` - получение конкретной записи настроения
- `PUT /api/v1/moods/:id` - обновление записи настроения
- `DELETE /api/v1/moods/:id` - удаление записи настроения
- `GET /api/v1/moods/by-date` - получение настроения по дате
- `GET /api/v1/moods/date-range` - получение настроений по диапазону дат
- `GET /api/v1/moods/latest` - получение последней записи настроения
- `POST /api/v1/moods/upsert-by-date` - создание или обновление настроения на дату
- `GET /api/v1/moods/stats` - статистика настроений
- `GET /api/v1/moods/trends` - анализ трендов настроения

### ✅ Google Calendar Integration API (COMPLETED)
- `GET /api/v1/google/auth-url` - получение URL для OAuth2 авторизации
- `POST /api/v1/google/callback` - обработка OAuth2 callback
- `GET /api/v1/google/integration` - получение информации об интеграции
- `DELETE /api/v1/google/integration` - отключение интеграции
- `GET /api/v1/google/calendars` - получение списка календарей пользователя

### ✅ Google Calendar Sync API (COMPLETED)
- `POST /api/v1/google/calendar-syncs` - создание конфигурации синхронизации
- `GET /api/v1/google/calendar-syncs` - получение списка конфигураций синхронизации
- `PUT /api/v1/google/calendar-syncs/:id` - обновление конфигурации синхронизации
- `DELETE /api/v1/google/calendar-syncs/:id` - удаление конфигурации синхронизации
- `POST /api/v1/google/calendar-syncs/:id/sync` - запуск синхронизации вручную

### ✅ Frontend React Application (COMPLETED)

#### Реализованные компоненты:
- **Layout** - основной макет приложения с боковой навигацией
- **Authentication Pages** - страницы входа и регистрации
- **Dashboard** - главная страница с обзором целей, событий и настроения
- **Calendar Page** - календарь с FullCalendar интеграцией
- **Event Modal** - модальное окно для создания/редактирования событий

#### Реализованные возможности:
- **Redux Toolkit** - управление состоянием приложения
- **React Router** - навигация между страницами
- **Ant Design** - UI компоненты
- **FullCalendar** - интерактивный календарь с drag-and-drop
- **TypeScript** - типизация всего кода
- **Axios** - HTTP клиент с автоматическим обновлением токенов
- **Environment Variables** - конфигурация для разных окружений

#### Структура проекта:
```
web/
├── src/
│   ├── components/        # Переиспользуемые компоненты
│   │   ├── Calendar/      # Календарные компоненты
│   │   ├── Common/        # Общие компоненты
│   │   └── Layout/        # Компоненты макета
│   ├── pages/             # Страницы приложения
│   ├── services/          # API сервисы
│   ├── store/             # Redux store и слайсы
│   ├── types/             # TypeScript типы
│   ├── hooks/             # Кастомные хуки
│   └── utils/             # Утилиты
├── public/                # Статические файлы
└── dist/                  # Собранное приложение
```

### ✅ Goal Management UI Components (COMPLETED)
- **GoalsPage** - полный интерфейс управления целями с:
  - Статистические карточки (общее количество, активные, завершенные, средний прогресс)
  - Поиск и фильтрация по статусу и категории
  - Список целей с индикаторами статуса и приоритета
  - Интеграция с детальной панелью целей
- **GoalModal** - модальное окно для создания/редактирования целей:
  - SMART framework для создания целей
  - Категории с описаниями и эмодзи
  - Приоритеты с цветовыми индикаторами
  - Валидация форм и обработка ошибок
- **GoalDetailPanel** - детальная панель цели:
  - Полная информация о цели
  - Обновление прогресса в реальном времени
  - Управление задачами (tasks) и этапами (milestones)
  - Модальные окна для создания/редактирования задач и этапов

### ✅ Mood Tracking Interface (COMPLETED)
- **MoodsPage** - полный интерфейс отслеживания настроения:
  - 5-уровневая система эмодзи (😢 😔 😐 😊 😄)
  - Интерактивный emoji selector с описаниями
  - Быстрые действия для записи настроения сегодня
  - Переключение между календарным и списочным видом
- **Mood Calendar** - интеграция с календарем:
  - Отображение эмодзи на датах в календаре
  - Клик по дате для создания/редактирования записи
  - Месячные агрегаты с средним настроением
- **Mood Statistics & Analytics**:
  - Статистические карточки (дни отслеживания, средний уровень, тренд)
  - Еженедельный обзор с детализацией по дням
  - Распределение настроений с progress bars
  - Цветовая индикация по уровням настроения
- **Mood Modal** - форма записи настроения:
  - Интуитивный выбор уровня настроения
  - Теги для описания влияющих факторов
  - Заметки для рефлексии и контекста
  - Валидация и обработка ошибок

### ✅ Settings Page with Google Integration (COMPLETED)
- **SettingsPage** - комплексная страница настроек:
  - Настройки профиля пользователя с валидацией
  - Изменение имени и email (email заблокирован)
  - Обновление данных через Redux и API
- **Google Calendar Integration Management**:
  - Подключение/отключение Google Calendar с OAuth2
  - Список доступных календарей пользователя
  - Настройка синхронизации для каждого календаря
  - Управление направлением синхронизации (bidirectional/from_google/to_google)
  - Автоматическая синхронизация с настраиваемыми интервалами
  - Ручной запуск синхронизации для каждого календаря
  - Разрешение конфликтов (Google wins/Local wins/Manual)
- **CalendarSyncModal** - детальные настройки синхронизации:
  - Выбор направления синхронизации с эмодзи-индикаторами
  - Настройки автосинхронизации и интервалов
  - Опции синхронизации прошлых и будущих событий
  - Стратегии разрешения конфликтов
- **Notification Settings**:
  - Email уведомления для событий
  - Ежедневные напоминания о mood tracking
  - Еженедельные обновления прогресса целей
  - Настройка времени напоминаний
- **Data & Privacy**:
  - Экспорт всех данных пользователя
  - Удаление аккаунта с подтверждением
  - Ссылки на Privacy Policy и Terms of Service

### ✅ Mobile & PWA Implementation (COMPLETED)
- **Responsive Design** - полная мобильная адаптация:
  - Адаптивная сетка с breakpoints для всех устройств
  - Mobile-first подход к дизайну компонентов
  - Отзывчивые карточки и списки
  - Оптимизированные размеры и отступы для мобильных устройств
- **Mobile Navigation**:
  - Bottom navigation bar для мобильных устройств
  - Floating Action Button с Speed Dial для быстрых действий
  - Drawer menu с полным функционалом
  - Badge индикаторы для уведомлений и счетчиков
- **Notification System**:
  - Контекстные уведомления с React Context
  - Типизированные toast сообщения (success, error, info, warning)
  - Системные уведомления (напоминания о mood, прогресс целей)
  - Действия в уведомлениях (кнопки, навигация)
  - Настраиваемые позиции и длительность показа
- **Progressive Web App (PWA)**:
  - Web App Manifest с полными метаданными
  - Service Worker для offline функциональности
  - App-like поведение на мобильных устройствах
  - Install prompts и shortcuts
  - Background sync и push notifications поддержка
- **Mobile Components**:
  - MobileCard - универсальный компонент карточек
  - ResponsiveLayout - адаптивный макет
  - MobileDashboardPage - оптимизированная главная страница
  - Контекстные меню и quick actions

### ✅ Performance Optimization (COMPLETED - 26.07.2025)
- **Bundle Size Optimization**: Уменьшен с 1,580.67 kB до максимального chunk 927.42 kB (41% improvement)
- **Lazy Loading**: Все страницы теперь загружаются динамически через React.lazy()
- **Code Splitting**: Настроены manual chunks для оптимального разделения библиотек:
  - react: React core (11.83 kB)
  - router: React Router (32.78 kB) 
  - redux: Redux Toolkit (26.30 kB)
  - antd: Ant Design (927.42 kB)
  - mui: Material-UI (0.75 kB)
  - calendar: FullCalendar (259.55 kB)
  - axios: HTTP client (35.41 kB)
- **Suspense Integration**: Добавлен Suspense wrapper для обработки lazy loading
- **Production Ready**: Приложение готово к deployment с оптимизированными бандлами

### 📋 Todo: Next Frontend Components
- Дополнительные PWA features (offline sync, background tasks)
- Улучшенная анимационная система и микроинтеракции
- Advanced accessibility (WCAG 2.1 compliance)

### 📋 Todo: Next APIs
- Webhook endpoints для real-time синхронизации с Google Calendar
- Batch operations для массовых операций с событиями