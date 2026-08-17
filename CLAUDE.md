# CLAUDE.md

Monorepo: `backend/` (Go) + `frontend/` (TypeScript/React), собирается в один бинарь —
бэкенд раздаёт собранный фронт через `embed.FS`.

Настройка окружения, docker, Storybook — в `README.md`.

## Commands

Всё через `mise run`, а не напрямую `go` / `pnpm`. Из корня — с префиксом, из подкаталога — без:

```
mise run checks --fix          # линтеры + typecheck по всему монорепо
mise run tests                 # тесты по всему монорепо
mise run //backend:tests       # только бэкенд
mise run //frontend:checks     # только фронтенд
cd frontend && mise run checks # то же самое из подкаталога
```

Прочее: `//backend:seed` (пересоздаёт `bin/app.db`), `//backend:token`, `//backend:coverage`,
`//frontend:app` (`--back` — прокси на локальный бэкенд), `//frontend:build`,
`//frontend:storybook`.

Git-хуков нет — коммит и пуш ничего не гоняют. Всё вместе (`checks` + `tests` +
`//frontend:build`) запускает `mise run ci`; это команда для человека перед пушем.

При работе гоняй точечно то, что затронул, а не `ci` целиком: `//backend:tests`,
`//frontend:tests`, `//frontend:checks`. Задачи молчат, пока всё зелёное — пустой вывод
и нулевой exit code означают успех, при падении печатается полный лог. Покрытие в обычный
прогон тестов не входит, для него есть `//backend:coverage`.

## Workflow

Работаем через pull request'ы на GitHub. В `main` напрямую не коммитим.

1. Ветка под задачу от свежего `main`.
2. Коммиты в ветку, ревью диффа локально через **revdiff** — до пуша.
3. Пуш и `gh pr create`.
4. Мержим **squash**'ем — один PR становится одним коммитом в `main`, история линейная.
   Обычный мерж-коммит — только когда коммиты ветки ценны по отдельности
   (пример: PR #39, где каждый апдейт зависимости откатывается независимо).
5. После мержа remote-ветка удаляется сама (`delete_branch_on_merge`), локально:
   ```
   git checkout main && git pull
   git branch -D <branch>
   ```
   `-D`, а не `-d`: после squash'а коммит ветки не предок `main`, и `-d` ругается.

Пуш и создание PR — только по явной просьбе: репозиторий публичный.

## Backend

Слои: **API handlers → Domain (`Store`) → Storage (`dataStorage`)**.

- `domain/` — агрегаты (`User`, `Track`, `Workout`, `Exercise`, `Protocol`) и `Store` с методами
  по ним. Интерфейс `storage` объявлен в `domain/store.go` — инверсия зависимостей,
  реализация в `storage/datastorage/`.
- Поведение живёт на агрегате (`Workout.ApplyUpdate`, `Workout.Ref`), не на `Store`.
- **Аутентификация — на уровне middleware** (`session.Authenticator`).
  **Авторизация (проверка прав) — в домене**, в методах `Store`, не в хендлерах.
- Доменные ошибки: `ErrNotFound`, `ErrForbidden` (`domain/errors.go`); хендлеры мапят их
  через `httpx.RenderDomainError`.
- Хендлеры тонкие: распарсить body → `toDomain()` → вызвать `Store` → отрендерить ответ.
- API в RPC-стиле: клиент шлёт объект целиком, включая ID, в теле запроса.
- Роутер — chi, префикс `/api/v1`, все роуты собраны в `api/main.go`.

## Frontend

React + Mantine + TanStack Router + TanStack Query. Файловый роутинг в `src/routes/`
(`routeTree.gen.ts` генерируется, руками не править).

- Состояние авторизации — хук `useAuth()` из `@/auth/hooks`.
- Редиректы по авторизации — в компонентах страниц через `<Navigate to="..." />`,
  не в `beforeLoad` на уровне роута.
- Гарды оборачивают компонент страницы: `<RequireAuth>` — только для залогиненных,
  `<TrackOwnerOnly>` — только для владельца трека.
- Загрузка: `<PageSkeleton />`. Ссылки между роутами: `<RouteLink />`.
- Права проверяются по полю `IsOwner: bool` из ответа API, а не сравнением ID на клиенте.

Подробнее: `frontend/docs/concepts/routing.md`, `frontend/docs/concepts/queries.md`.

### Придержанные мажоры

Два пакета намеренно не на последней версии — не забывать при следующем обновлении:

- **eslint 9**, не 10. `eslint-plugin-react` зовёт `context.getFilename()`, убранный в 10-ке,
  и падает на загрузке правила. `eslint-config-mantine` тоже упирается в `^9.9.1`.
- **typescript 6**, не 7. `typescript-eslint` бросает исключение на TS 7,
  см. typescript-eslint#10940 — поддержку обещают с TS 7.1.

### TanStack Query (v5)

- `isPending` вместо `isLoading`.
- `data?.field` предпочтительнее `isSuccess && data.field` для простых случаев.
- `select` — чтобы разворачивать ответ API (`{ data: T }` → `T`).
- Общие данные грузит layout-компонент, дети берут их из кэша мгновенно.

## Задачи и бэклог

`docs/neskuchka.taskpaper` — формат TaskPaper. Секции: `INBOX` (неразобранное),
`BACKLOG` (всё остальное, плоско), `THINK` (вопросы без решения), `Archive` (`@done`).
Между `INBOX` и `BACKLOG` — проекты: задача больше чем на один пункт, расписанная
подзадачами. Заводится, когда задачу начинают расписывать или брать в работу;
после завершения удаляется целиком, а не архивируется.

Теги ставим в конце строки, в порядке «область → тип → модификатор»:

- **область** — `@backend` `@frontend` `@infra`
- **тип** — `@feature` (новая возможность), `@bug`, `@debt` (починить сломанное), `@arch`
- **модификатор** — `@high` `@blocked` `@flag` `@done`

Больше тегов не заводим: срез по чему угодно ещё берётся поиском.

Файл в `.gitignore`: репозиторий публичный, а бэклог пока черновой. Когда причешется —
убрать строчку из `.gitignore` и закоммитить.
