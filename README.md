# Neskuchka

Проект для тренировок


## Разработка
Выполняем: `mise run info`, если не работает - надо настроить окружение

### Линтинг
Линтингом занимается [ruff](https://docs.astral.sh/ruff/)

Запуск через `mise` (или напрямую) из папки проекта:
```
mise run lint
```

Вывод с ошибками:
```
[lint] $ ruff check app
app/domain/track.py:3:29: F401 [*] `app.domain.user.User` imported but unused
  |
1 | from typing import NewType
2 | from app.domain.entity import EntityModel
3 | from app.domain.user import User, UserId
  |                             ^^^^ F401
  |
  = help: Remove unused import: `app.domain.user.User`
```

Ошибки можно искать на сайте по коду: [F401](https://docs.astral.sh/ruff/rules/unused-import/)

Поправить можно вручную или через команду:
```
mise run lint --fix
[lint] $ ruff check app --fix
Found 1 error (1 fixed, 0 remaining).
```

Однако не все ошибки можно починить автоматически

**Второй тип ошибок** - это проверка форматирования, ругаться будет примерно:
```
[pre-commit] $ (cd backend && mise run lint)
[lint] $ ruff format --check . && ruff check .
Would reformat: app/main.py
``` 

Нужно запустить руками:
```
ruff format .
ruff format
1 file reformatted, 16 files left unchanged
```

### Тестирование
Для запуска тестов используется [pytest](https://docs.pytest.org/en/latest/)

```
mise run test
```

### Hurl
_экспериментально пока_
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

Результат примерно такой:
```
➜  neskuchka git:(main) ✗ git commit

[pre-commit] $ mise run lint
[lint] $ ruff check app
All checks passed!
On branch main
...
```

#### Backend
```
# Переходит в папку бекенда
cd backend

# Добавляем конфиг файл `.mise.toml` в доверяемые 
mise trust

# Устанавливаем зависимости
mise run install
```

Для управления зависимостей используется [uv](https://docs.astral.sh/uv/), но пока в очень простом виде


**Caddy** (пока экспериментальный)
Для прокси используем [Caddy](https://caddyserver.com/)

Даем привилегии использовать порты 80 и 443:
```
sudo setcap cap_net_bind_service=+ep $(which caddy)
```

Запускаем caddy:

```
caddy start
```

Теперь можно бекенд доступен по адресу: `https://neskuchka.localhost`

