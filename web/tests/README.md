# E2E Testing Infrastructure

## Overview

Проект использует Playwright для end-to-end тестирования, обеспечивающего проверку ключевых пользовательских сценариев в реальном браузерном окружении.

## Требования

- **Node.js**: >= 18.19.0 (для поддержки ESM modules в Playwright)
- **npm**: >= 9.6.7
- **Браузеры**: Автоматически устанавливаются через Playwright

## Установка

```bash
# Установка зависимостей
npm install

# Установка браузеров Playwright
npx playwright install
```

## Запуск тестов

### Основные команды:

```bash
# Запуск всех E2E тестов
npm run test:e2e

# Запуск с UI интерфейсом (интерактивный режим)
npm run test:e2e:ui

# Запуск в headed режиме (видимые браузеры)
npm run test:e2e:headed

# Debug режим (пошаговое выполнение)
npm run test:e2e:debug

# Просмотр отчета о последнем запуске
npm run test:e2e:report
```

### Запуск конкретных тестов:

```bash
# Запуск одного файла
npx playwright test tests/e2e/auth.spec.ts

# Запуск тестов с определенным названием
npx playwright test --grep "should redirect to login"

# Запуск только на Chrome
npx playwright test --project=chromium
```

## Структура тестов

```
tests/e2e/
├── auth.spec.ts        # Тесты аутентификации
├── dashboard.spec.ts   # Тесты главной страницы
├── calendar.spec.ts    # Тесты календаря
├── goals.spec.ts       # Тесты управления целями
├── moods.spec.ts       # Тесты отслеживания настроения
└── navigation.spec.ts  # Тесты навигации
```

## Описание тестовых сценариев

### 🔐 Authentication (auth.spec.ts)
- Редирект на страницу входа для неавторизованных пользователей
- Отображение формы регистрации
- Валидация полей формы входа
- Обработка неверных учетных данных

### 📊 Dashboard (dashboard.spec.ts)
- Отображение основных компонентов дашборда
- Responsive дизайн на мобильных устройствах
- Проверка редиректа на авторизацию

### 📅 Calendar (calendar.spec.ts)
- Отображение интерфейса календаря
- Навигация по календарю (prev/next)
- Responsive поведение
- Проверка контролов просмотра (Month/Week/Day)

### 🎯 Goals (goals.spec.ts)
- Отображение страницы целей
- Модальное окно создания цели
- Список целей и пустое состояние
- SMART цели функциональность (если доступна)

### 😊 Moods (moods.spec.ts)
- Интерфейс отслеживания настроения
- Emoji селектор для выбора настроения
- Календарь/история настроений
- Статистика и аналитика

### 🧭 Navigation (navigation.spec.ts)
- Навигация между страницами
- Мобильная навигация
- Browser back/forward
- Активное состояние навигации

## Конфигурация

### Браузеры
Тесты запускаются на:
- **Desktop**: Chrome, Firefox, Safari
- **Mobile**: Chrome (Pixel 5), Safari (iPhone 12)

### Настройки
- **Base URL**: http://localhost:5173
- **Timeout**: 60 минут для полного прогона
- **Retries**: 2 попытки на CI, 0 локально
- **Screenshots**: При падении тестов
- **Video**: При падении тестов
- **Trace**: При повторном запуске после падения

## CI/CD Integration

Тесты автоматически запускаются через GitHub Actions:
- При push в ветки: main, develop, claude
- При создании Pull Request в main, develop
- Артефакты (отчеты, скриншоты) сохраняются на 30 дней

## Troubleshooting

### Node.js версия
```bash
# Проверка версии Node.js
node --version

# Если версия < 18.19.0, обновите Node.js
# Рекомендуется использовать nvm для управления версиями
```

### Браузеры не установлены
```bash
npx playwright install --with-deps
```

### Тесты падают локально
```bash
# Убедитесь что dev сервер запущен
npm run dev

# Или используйте встроенный webServer (автоматически)
npm run test:e2e
```

### Debug тестов
```bash
# Пошаговое выполнение
npm run test:e2e:debug

# Или с UI интерфейсом
npm run test:e2e:ui
```

## Архитектурные принципы

### 🎯 Robust Selectors
Тесты используют множественные стратегии поиска элементов:
- Семантические селекторы (`role`, `aria-*`)
- Data attributes (`data-testid`)
- Text содержимое (с поддержкой i18n)
- CSS классы как fallback

### 🔄 Graceful Handling
- Тесты адаптируются к состоянию аутентификации
- Обработка редиректов на страницу входа
- Поддержка различных состояний приложения

### 🌐 Cross-Platform
- Тестирование на desktop и mobile viewports
- Multi-browser support
- Responsive design validation

### 📱 Mobile-First
- Специальные проверки для мобильных интерфейсов
- Touch navigation тестирование
- Viewport адаптация

## Best Practices

1. **Селекторы**: Используйте стабильные селекторы (data-testid, role)
2. **Ожидания**: Всегда используйте explicit waits
3. **Изоляция**: Каждый тест должен быть независимым
4. **Данные**: Используйте детерминированные тестовые данные
5. **Очистка**: Очищайте состояние после тестов
6. **Скорость**: Минимизируйте unnecessary actions

## Планы развития

- [ ] **Authentication Mock**: Bypass реальной аутентификации для ускорения тестов
- [ ] **Test Data Factory**: Генерация консистентных тестовых данных
- [ ] **Visual Regression**: Скриншот сравнения для UI стабильности
- [ ] **API Mocking**: Изоляция frontend тестов от backend
- [ ] **Performance Testing**: Web Vitals и timing проверки
- [ ] **Accessibility Testing**: Автоматизированные a11y проверки