# 🔧 Troubleshooting Guide

[← Shortcuts](./shortcuts.md) | [Back to Instructions →](../instructions.md)

## 🚨 Common Issues

### "Documentation out of sync"
**Problem**: Документация не соответствует коду

**Solution**:
```bash
init              # Проверить структуру
docs-sync         # Синхронизировать
status            # Проверить результат
```

**Prevention**:
- Always update docs with code
- Use pre-commit hooks
- Regular sync checks

---

### "Command not recognized"
**Problem**: AI не понимает команду

**Possible Causes**:
1. Опечатка в команде
2. Используется старый alias
3. Команда не существует

**Solution**:
```bash
help              # Показать все команды
explain <command> # Объяснить команду
```

---

### "Tests failing"
**Problem**: Тесты не проходят

**Debug Steps**:
1. `test --verbose` - детальный вывод
2. `test --single <test>` - отдельный тест
3. `test --debug` - режим отладки

**Common Fixes**:
- Clear cache: `clean-cache`
- Update snapshots: `test -u`
- Check environment: `check-env`

---

### "Build errors"
**Problem**: Проект не собирается

**Quick Fixes**:
```bash
clean-all         # Очистить все кеши
install --force   # Переустановить зависимости
typecheck         # Проверить типы
build --verbose   # Детальная сборка
```

---

### "Performance issues"
**Problem**: Приложение работает медленно

**Diagnostic**:
```bash
audit --performance  # Анализ производительности
analyze --bundle     # Анализ размера
profile             # Профилирование
```

## 🔍 Debugging Commands

### Information Gathering
```bash
context           # Полный контекст проекта
status --detailed # Детальный статус
explain --system  # Объяснить систему
check-env         # Проверить окружение
```

### Problem Identification
```bash
audit --errors    # Найти ошибки
audit --warnings  # Найти предупреждения
check-deps        # Проверить зависимости
validate          # Валидация проекта
```

### Quick Fixes
```bash
fix --auto        # Автоисправление
fix --lint        # Исправить стиль
fix --types       # Исправить типы
fix --format      # Форматирование
```

## 📋 Checklists

### Before Asking for Help
- [ ] Read error message carefully
- [ ] Check documentation
- [ ] Try common fixes
- [ ] Isolate the problem
- [ ] Prepare minimal reproduction

### After Fixing Issue
- [ ] Document the solution
- [ ] Add test to prevent regression
- [ ] Update troubleshooting guide
- [ ] Share knowledge with team

## 🆘 Getting Help

### Self-Help Resources
1. This troubleshooting guide
2. `explain --error <e>`
3. `help --topic <topic>`
4. Project documentation

### When to Escalate
- Security issues → immediate
- Data loss risk → immediate
- Blocking entire team → high priority
- Performance regression → medium priority

### How to Report Issues
```markdown
## Issue Template
**Command**: {what command failed}
**Expected**: {what should happen}
**Actual**: {what happened}
**Steps**: {how to reproduce}
**Environment**: {output of check-env}
**Logs**: {relevant logs}
```

### Using Bug Command
```bash
# Простой репорт
bug "Описание проблемы"

# Детальный репорт
bug --critical --reproduce
# Следуйте инструкциям для заполнения всех полей

# С контекстом
bug --context "После деплоя v2.1" --component "api" "500 ошибка на /users"
```

## 🛠️ Recovery Procedures

### Corrupted State
```bash
backup --emergency   # Сохранить текущее
clean-all           # Очистить все
init --fresh        # Свежая инициализация
restore --latest    # Восстановить данные
```

### Failed Deployment
```bash
rollback --immediate # Откатить немедленно
status --production  # Проверить прод
hotfix              # Применить исправление
monitor             # Мониторить результат
```

### Lost Work
```bash
recover --git       # Из git reflog
recover --backup    # Из backup
recover --cache     # Из кеша
```

## 💡 Prevention Tips

### Daily Practices
- Run `status` at start
- Commit frequently
- Test before commit
- Update docs immediately

### Weekly Practices
- Run `reflect`
- Clean old branches
- Update dependencies
- Review alerts

### Project Health
- Maintain >80% test coverage
- Keep bundle size <1MB
- Monitor error rates
- Regular security audits

## 🔄 Reset Commands

### Soft Reset
```bash
reset --soft     # Сохранить изменения
clean-cache      # Очистить кеши
restart          # Перезапустить
```

### Hard Reset
```bash
reset --hard     # Удалить изменения
clean-all        # Очистить все
init --force     # Переинициализация
```

### Emergency Reset
```bash
emergency-stop   # Остановить все
backup --force   # Принудительный backup
reset --emergency # Экстренный сброс
```

## 📞 Support Contacts

### Internal
- Tech Lead: {contact}
- DevOps: {contact}
- Security: {contact}

### External
- Hosting Support: {contact}
- API Support: {contact}