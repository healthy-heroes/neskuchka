## Backend
### Task: Refactoring tests
- Тесты на API
    - Написать тесты
    - Написать фикстуры
- подумать
    - in-memory база для тестов?
    - async?


### Task: Frontend
- Поднять nextjs и выбрать между spectrum и shadcn

### Task: Refactor API Schema
- Сделать удобно для фронта

### Task: Docker
- docker compose

### Deploy
- CI/CD
- Cloud


### Task: Refactoring
- Подключить mypy
- Рефакторинг dependencies
- Рефакторинг моделей (уменьшить?)
- Обработка ошибок
- Схема возвращаемых значений из API (ошибки, данные)
- [v] lifespans
- сделать типы для возвращаемых ошибок, чтобы не хардкодить текст, поправить тесты и сделать проверку
- документация для апи



## Part 2
### Task: Auth
- Users
- Auth
- JWT

### Task: Admin
- Админка для добавления тренировок и упражнений


### Task: Логирование
- поизучать чем логируют, прикрутить
- opentelemetry

### Task: Postgres
- connection pool
- Миграции
- async db and session
- тестирование через https://testcontainers.com/guides/getting-started-with-testcontainers-for-python/

## Utils
- Сделать так чтобы `ruff lint` проверял форматирование, которое приносит `ruff format`
- Подумать над тем, чтобы `ruff` и `uv` использовать без `mise`
- Посмотреть на pyproject.toml технологию
