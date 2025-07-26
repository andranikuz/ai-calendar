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
- [ ] OAuth2 интеграция с Google
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
- [ ] OAuth2 настройка для Google API
- [ ] Двустороння синхронизация событий
- [ ] Webhook для real-time обновлений
- [ ] Обработка конфликтов данных

## Этап 4: Frontend основы (3-4 недели)

### 4.1 React приложение
- [ ] Create React App с TypeScript
- [ ] Настройка Redux Toolkit для состояния
- [ ] Роутинг с React Router
- [ ] Базовые компоненты и стили

### 4.2 Календарный интерфейс
- [ ] Интеграция FullCalendar или custom реализация
- [ ] Просмотры: день, неделя, месяц
- [ ] Drag-and-drop для событий
- [ ] Quick создание событий

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

### 📋 Todo: Next APIs
- Google Calendar integration endpoints