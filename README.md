# Neskuchka

Проект для тренировок


## Разработка
Выполняем: `mise run info`, если не работает - надо настроить окружение

###  Настройка окружения
В случае с виндой лучше использовать WSL

Поставить [mise](https://mise.jdx.dev/)

Клонируем репозиторий и устанавливаем зависимости
```
git clone git@github.com:healthy-heroes/neskuchka.git
cd neskuchka

# Добавляем конфиг файл `.mise.toml` в доверяемые 
mise trust
```

Настраиваем pre-комитные хуки, чтобы автоматом запускать линтинг перед каждым комитом.

[Фича экспериментальная](https://mise.jdx.dev/cli/generate/git-pre-commit.html), поэтому нужно ее предварительно включить:

```
mise settings experimental=true
mise generate git-pre-commit --write --task=pre-commit
```

### Запуск локального сервера
Нужно запустить бекенд и фронтенд в разных терминалах

После открыть http://127.0.0.1:5173/

#### Backend
```
# Переходит в папку бекенда
cd backend

# Устанавливаем зависимости
go mod tidy

# Запуск локального сервера
go run app/main.go server

```


#### Frontend

```
cd frontend/app

# Устанавливаем зависимости
mise run install

# Запускаем приложение
mise run app
```






