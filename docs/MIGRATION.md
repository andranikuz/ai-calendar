# 🔄 Migration Guide - v1.0 to v2.0

## 📋 Overview

Это руководство поможет вам мигрировать с старой системы инструкций (v1.0) на новую (v2.0) с английскими командами и улучшенной структурой.

## 🚀 Quick Migration (5 минут)

### Step 1: Backup Current State
```bash
# Создайте backup текущей документации
cp -r docs/ docs-backup/
cp CLAUDE.md CLAUDE.md.backup
```

### Step 2: Install New System
1. Распакуйте архив в корень проекта
2. Файлы будут размещены:
    - `instructions.md` - в корне
    - `.ai/` - папка с инструкциями
    - `.ai-config.json` - конфигурация

### Step 3: Initialize
```bash
# Запустите инициализацию
AI: run command "init"
```

### Step 4: Migrate Documentation
```bash
# Автоматическая миграция
AI: run command "migrate-docs"
```

## 📝 Manual Migration

### Mapping Old Commands to New

| Старая команда | Новая команда | Изменения |
|----------------|---------------|-----------|
| инициализация | `init` | Те же функции |
| статус | `status` | Расширенный вывод |
| продолжи разработку | `develop` | Больше опций |
| рефлексия | `reflect` | Структурированный отчет |
| идеи развития | `innovate` | Категоризация идей |
| создать issue | `bug` | Интерактивный режим |
| список задач | `bug-list` + `plan` | Разделение багов и features |

### New Features in v2.0

1. **Shortcuts & Aliases**
    - `s` вместо `status`
    - `d` вместо `develop`
    - `b` вместо `bug` (новое!)
    - Custom aliases in config

2. **Command Composition**
    - Chain commands: `status + develop + test`
    - Conditional: `test ? commit : fix`

3. **Smart Commands**
    - `auto` - AI выбирает команду
    - `explain` - объяснение любой части
    - `bug` - интерактивный репорт багов

4. **Bug Management** (новое!)
    - `bug` - создать репорт о баге
    - `bug-list` - список всех багов
    - `bug-fix` - работа над исправлением
    - Интеграция с workflow

5. **Better Structure**
    - Модульные инструкции
    - Отдельные файлы для каждой области
    - JSON конфигурация

## 🔧 Configuration Migration

### Old Style (CLAUDE.md)
```markdown
Команды жестко заданы в файле
Нет возможности кастомизации
```

### New Style (.ai-config.json)
```json
{
  "aliases": {
    "my-flow": "status + develop + test"
  },
  "settings": {
    "language": "ru",
    "commandLanguage": "en"
  }
}
```

## 📊 Validation

После миграции проверьте:

```bash
# 1. Проверка структуры
AI: run "validate --structure"

# 2. Проверка команд
AI: run "help" (должны показаться все команды)

# 3. Проверка документации
AI: run "docs-check"

# 4. Тестовый прогон
AI: run "status"
```

## 🚨 Common Issues

### Issue: "Command not found"
**Fix**: Используйте английские названия команд

### Issue: "Documentation mismatch"
**Fix**: Запустите `docs-sync`

### Issue: "Old aliases not working"
**Fix**: Добавьте их в `.ai-config.json`

## 🎯 Best Practices

1. **Start Fresh**
    - Лучше начать с чистой инициализации
    - Импортируйте старые данные постепенно

2. **Learn Shortcuts**
    - Используйте короткие команды
    - Создайте свои aliases

3. **Use English Commands**
    - Консистентность с другими инструментами
    - Легче интегрировать

4. **Keep Russian Output**
    - Документация на русском
    - Коммиты на русском
    - Отчеты на русском

## 📅 Migration Timeline

### Day 1
- Backup everything
- Install new system
- Run init
- Basic testing

### Day 2-3
- Migrate active tasks
- Update team documentation
- Configure aliases

### Week 1
- Full transition
- Team training
- Optimize workflow

## 🆘 Getting Help

Если возникли проблемы:

1. Check `troubleshooting.md`
2. Run `explain --migration`
3. Contact support

## ✅ Migration Checklist

- [ ] Backup created
- [ ] New system installed
- [ ] Init completed
- [ ] Documentation migrated
- [ ] Commands working
- [ ] Team notified
- [ ] Aliases configured
- [ ] Validation passed

## 🎉 Welcome to v2.0!

Новая система предоставляет:
- ⚡ Более быструю работу
- 🎯 Лучшую организацию
- 🔧 Больше возможностей
- 📊 Улучшенную аналитику

Успешной миграции!