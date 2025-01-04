# Neskuchka

Проект для тренировок


## Разработка
- ruff

Выполняем: `mise run info`, если не работает - надо настроить окружение

### Запуск приложения
```
mise run app
```

Сервер стартует на http://localhost:8000
Документация доступна по адресу http://localhost:8000/docs или http://localhost:8000/redoc

###  Настройка окружения
В случае с виндой лучше использовать WSL

Поставить [mise](https://mise.jdx.dev/) и [uv](https://docs.astral.sh/uv/)

Клонируем репозиторий и устанавливаем зависимости
```
git clone git@github.com:healthy-heroes/neskuchka.git
cd neskuchka

# Добавляем конфиг файл `.mise.toml` в доверяемые 
mise trust

# Устанавливаем зависимости
mise run install
```