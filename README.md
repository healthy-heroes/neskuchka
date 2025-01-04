# Neskuchka

Проект для тренировок


## Разработка
- ruff

Выполняем: `mise run info`, если не работает - надо настроить окружение

### Hurl
[Hurl](https://github.com/Orange-OpenSource/hurl) позволяет выполнять запросы и делать всякие ассерты

В папке `examples` лежат .hurl файлы, в которых написаны легковесные тесты, запустить их можно так:

```
mise run test-api
```


### Запуск приложения
```
mise run app
```

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