# Neskuchka [![Coverage Status](https://coveralls.io/repos/github/healthy-heroes/neskuchka/badge.svg?branch=main)](https://coveralls.io/github/healthy-heroes/neskuchka?branch=main)

Проект для тренировок


## Разработка
Выполняем: `mise run info`, если не работает - надо настроить окружение

###  Настройка окружения
Поставить [mise](https://mise.jdx.dev/)

Клонируем репозиторий и устанавливаем зависимости
```
git clone git@github.com:healthy-heroes/neskuchka.git
cd neskuchka

# Добавляем конфиг файл `.mise.toml` в доверяемые 
mise trust --all
```

Используем экспериментальные фичи mise, поэтому включаем настройку:
```
mise settings experimental=true
```

Настраиваем git-хуки:
```
## pre commit hook
mise generate git-pre-commit --write --task=pre-commit

## pre push hook
mise generate git-pre-commit --write --task=pre-push --hook=pre-push
```

**Для Windows**
Для разработки рекомендуется использовать WSL 2, а для запуска docker-compose должна быть включена [интеграция с WSL](https://docs.docker.com/desktop/features/wsl/)

### Запуск линтеров и тестов
Для запуска линтеров и тестов используется mise. Запускать можно в проекте или в корневом каталоге.

- `checks` - запуск линтеров, поддерживает флаг `--fix` для автоматического исправления ошибок
- `tests` - запуск тестов

При запуске из корня испольузется фича `experimental_monorepo_root`, позволяет удобно запускать задачи из любого проекта:
```
## запуск общей задачи проверки кода
mise run checks

## запуск задачи проверки кода для бэкенда
mise run //backend:checks

## что идентично
(cd backend && mise run checks)

## запуск задачи проверки кода для фронтенда с автоматическим исправлением ошибок
mise run //frontend:checks --fix

## запуск тестов фронтенда
mise run //frontend:tests
```

### Запуск локального сервера
Бэкенд можно поднимать локально или в Docker-контейнере. Контейнер хранит собранный фронтенд, поэтому не будет подхватывать изменения на фронтенде.

При разработке бэкенд можно поднимать в контейнере или локально, зависит от предпочтений. В контейнере удобнее, когда не нужно постоянно обновлять фронтенд и достаточно собранной версии посмотреть на работоспособность. В остальном разницы практически нет.

#### Запуск бэкенда в контейнере
Поднимает бэкенд с собранным фронтендом и Mailpit (интерфейс для перехвата писем):
```
docker compose -f docker-dev-compose.yml up -d --build
```

- Приложение: http://localhost:8080/
- Mailpit (почта): http://localhost:8025/

База данных хранится в Docker volume. При первом запуске автоматически заполняется начальными данными (seed). При перезапуске контейнера база сохраняется.

Для полного сброса базы:
```
docker compose -f docker-dev-compose.yml down -v
docker compose -f docker-dev-compose.yml up -d --build
```

#### Запуск бэкенда локально
Запуск производится с помощью `air`, которая прилетает через `mise`. `air` — это вотчер, конфигурация лежит в `app/.air.toml`, там же параметры запуска.

Для работы авторизации по email нужен запущенный Mailpit:
```
docker compose -f docker-dev-compose.yml up -d mailpit
```

Mailpit Web UI: http://localhost:8025/

```
cd backend

go mod tidy

# Заполнить базу начальными данными (идемпотентно, можно запускать повторно)
mise run seed

# Запуск локального сервера
cd app
air
```

### Запуск локального фронтенда
Для активной разработки фронтенда нужно его запустить локально:

```
cd frontend

# Устанавливаем зависимости
mise run install
```

Фронтенд запускается с проксей, которая перенаправляет запросы `*/api` в бекенд на localhost, в какой порт указывается переменной окружения `VITE_BACKEND_PORT`

Для удобства есть есть две команды с преднастроенными портами (при условии, что их не меняли при поднятии бекенда или докера)

```
mise run app            # фронт нацелен на докер (http://localhost:8080/)
mise run app --back     # Фронт нацелен на локальный бек (http://localhost:8000/)
```

В обоих случаях фронт поднимается на http://localhost:5173/ с вотчером изменений.

### Storybook
Для разработки компонентов можно использовать Storybook:

```
cd frontend
mise run storybook
```

Storybook поднимается на http://localhost:6006/

