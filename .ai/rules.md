# ⚠️ Smart Goal Calendar - Критические правила

## 🔥 КРИТИЧЕСКИЕ ПРАВИЛА (НАРУШЕНИЕ НЕДОПУСТИМО)

### 1. Документация превыше всего
- **ВСЕГДА** обновлять документацию при изменении кода
- **НИКОГДА** не коммитить без обновления соответствующих docs
- **ОБЯЗАТЕЛЬНО** поддерживать `docs/current-status.md` актуальным
- **ВСЕГДА** записывать в CHANGELOG.md выполнение команд

### 2. Безопасность и типизация
- **ЗАПРЕЩЕНО** использовать `any` типы в TypeScript
- **ОБЯЗАТЕЛЬНО** валидировать все пользовательские вводы  
- **НИКОГДА** не коммитить секреты или API ключи
- **ВСЕГДА** использовать structured errors в Go

### 3. Тестирование обязательно
- **НИКОГДА** не коммитить код без тестов для критической логики
- **ВСЕГДА** запускать `go test ./...` перед коммитом backend изменений
- **ОБЯЗАТЕЛЬНО** проверять `npm run build` перед коммитом frontend
- **ВСЕГДА** проверять TypeScript компиляцию

### 4. Clean Architecture соблюдение
- **ЗАПРЕЩЕНО** нарушать границы слоев (domain не зависит от внешних)
- **ОБЯЗАТЕЛЬНО** использовать интерфейсы для всех репозиториев
- **ВСЕГДА** следовать SOLID принципам
- **НИКОГДА** не смешивать бизнес-логику с HTTP handlers

---

## 📋 ОБЯЗАТЕЛЬНЫЕ ДЕЙСТВИЯ

### При каждом коммите
```bash
# Backend проверки
go test ./...
go vet ./...
go fmt ./...

# Frontend проверки  
cd web
npm run lint
npm run typecheck
npm run build
```

### При изменении API
- Обновить OpenAPI документацию
- Добавить/обновить тесты endpoints
- Проверить обратную совместимость
- Обновить frontend типы

### При изменении архитектуры
- Обновить `docs/architecture.md`
- Проверить все зависимости слоев
- Документировать breaking changes
- Обновить диаграммы

---

## 🚫 ЗАПРЕЩЕННЫЕ ДЕЙСТВИЯ

### Код качество
- ❌ Коммитить код с ошибками компиляции
- ❌ Использовать `console.log` в production React коде
- ❌ Игнорировать ESLint/golangci-lint ошибки
- ❌ Оставлять TODO комментарии без GitHub issues

### Безопасность  
- ❌ Хранить пароли/ключи в plaintext
- ❌ Пропускать input validation
- ❌ Использовать небезопасные SQL queries
- ❌ Логировать sensitive данные

### Архитектура
- ❌ Нарушать слои Clean Architecture
- ❌ Создавать циклические зависимости
- ❌ Смешивать domain логику с UI
- ❌ Делать прямые database вызовы из handlers

---

## ✅ ОБЯЗАТЕЛЬНЫЕ СТАНДАРТЫ

### Go Backend
```go
// ✅ Правильно: Интерфейсы в domain слое
type GoalRepository interface {
    Create(ctx context.Context, goal *Goal) error
    GetByID(ctx context.Context, id string) (*Goal, error)
}

// ✅ Правильно: Structured errors
var ErrGoalNotFound = errors.New("goal not found")

// ❌ Неправильно: Логика в HTTP handler
func (h *Handler) CreateGoal(c *gin.Context) {
    // Много бизнес-логики здесь - неправильно!
}
```

### React Frontend
```typescript
// ✅ Правильно: Строгая типизация
interface Goal {
  id: string;
  title: string;
  priority: Priority;
}

// ❌ Неправильно: any типы
const goal: any = response.data; // ЗАПРЕЩЕНО!

// ✅ Правильно: Error handling
try {
  await api.createGoal(goal);
} catch (error) {
  showNotification('error', 'Ошибка создания цели');
}
```

---

## 📐 ПРИНЦИПЫ КОММИТОВ

### Формат коммит сообщений
```
<type>(<scope>): <description>

[optional body]

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

### Типы коммитов
- `feat`: Новая функциональность
- `fix`: Исправление бага
- `docs`: Изменения документации
- `style`: Форматирование кода
- `refactor`: Рефакторинг без изменения функциональности
- `test`: Добавление тестов
- `chore`: Обновление build процесса

### Обязательные элементы
- Осмысленное описание на русском языке
- Детали изменений в body (если нужно)
- Ссылки на issues (если применимо)
- Подпись Claude Code

---

## 🎯 КАЧЕСТВЕННЫЕ ТРЕБОВАНИЯ

### Метрики минимального качества
- **Test Coverage**: Backend >70%, Frontend >60%  
- **TypeScript**: 100% typed, 0 any
- **ESLint**: 0 errors, warnings допустимы
- **Build**: Успешная сборка обязательна
- **Performance**: Lighthouse >85

### Code Review критерии
- Соответствие архитектурным принципам
- Покрытие тестами новой логики  
- Обновление документации
- Безопасность implementation
- Performance implications

---

## 🔄 WORKFLOW ПРАВИЛА

### TodoWrite обязательное использование
- **ВСЕГДА** создавать todolist для сложных задач
- **ТОЛЬКО** одна задача "in_progress" одновременно
- **НЕМЕДЛЕННО** отмечать задачи completed после завершения
- **НИКОГДА** не оставлять stale todos

### Документация workflow
- При выполнении команды `init` - создать коммит
- При выполнении команды `status` - обновить current-status.md
- При выполнении команды `develop` - обновить development-plan.md
- При выполнении команды `reflect` - создать файл в ideas/

### Git workflow
- **НИКОГДА** не пушить напрямую в main
- **ВСЕГДА** создавать feature branches для больших изменений
- **ОБЯЗАТЕЛЬНО** делать meaningful коммиты
- **НИКОГДА** не force push без крайней необходимости

---

## 🚨 КРИТИЧЕСКИЕ СЦЕНАРИИ

### При обнаружении security vulnerability
1. **НЕМЕДЛЕННО** прекратить работу
2. Создать приватный security issue
3. Не коммитить уязвимый код
4. Уведомить о проблеме

### При breaking changes
1. Обновить CHANGELOG.md с BREAKING CHANGE меткой
2. Обновить API документацию
3. Создать migration guide
4. Обновить версию в package.json/go.mod

### При production bugs
1. Создать hotfix branch
2. Minimal reproducible fix
3. Немедленное тестирование
4. Emergency deployment процедура

---

## 🎯 SMART GOAL CALENDAR СПЕЦИФИЧНЫЕ ПРАВИЛА

### Google Calendar Integration
- **ВСЕГДА** handle rate limiting
- **ОБЯЗАТЕЛЬНО** implement retry логику
- **НИКОГДА** не терять пользовательские данные при sync
- **ВСЕГДА** логировать sync операции

### PWA Requirements
- **ОБЯЗАТЕЛЬНО** поддерживать offline functionality
- **ВСЕГДА** handle network failures gracefully
- **НИКОГДА** не показывать белый экран
- **ВСЕГДА** показывать loading states

### Mood Tracking Privacy
- **НИКОГДА** не логировать mood данные в analytics
- **ВСЕГДА** шифровать sensitive данные
- **ОБЯЗАТЕЛЬНО** предоставлять data export
- **НИКОГДА** не передавать mood данные третьим сторонам

---

## 🔗 Проверка соблюдения правил

Эти правила проверяются через:
- Automated linting (ESLint, golangci-lint)
- Pre-commit hooks
- CI/CD pipeline
- Code review процесс
- Automated testing

**ПОМНИ**: Нарушение критических правил недопустимо и может привести к rollback изменений.