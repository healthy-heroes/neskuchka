# Примеры использования API

Документация по [Resterm](https://github.com/unkn0wn-root/resterm/blob/main/docs/resterm.md)

## Быстрый старт
Скопировать `resterm.env.example.json` в `resterm.env.json` и заполнить токены.

Получить токены:
```
cd backend
mise run token --email admin@example.com
```
Запустить установленный ранее Resterm:
```
resterm
```

## Полезное
- можно сгенерить выхлоп через `resterm --from-curl`

## Грабли
- не смог нормально отправить multipart. Сделано напрямую через curl, но этот запрос не видно в меню
