# Миграции схемы БД (goose) — план реализации

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Заменить `CREATE TABLE IF NOT EXISTS` на старте версионированными goose-миграциями, чтобы схему существующей прод-базы можно было менять.

**Architecture:** SQL-миграции лежат в `backend/app/storage/db/migrations/` и embed'ятся в пакет `db`. `NewSqliteEngine` после прагм прогоняет `goose` Provider API (`Up`) вместо `createSchema`. Миграция `0001` — текущая схема дословно, поэтому на существующей базе она no-op и только фиксирует версию.

**Tech Stack:** Go 1.26.6, `github.com/pressly/goose/v3` v3.27.3 (Provider API), `modernc.org/sqlite` (без cgo), sqlx, zerolog, testify.

**Spec:** `docs/superpowers/specs/2026-08-19-db-migrations-design.md`

## Global Constraints

- Ветка `db-migrations`. Пуш и создание PR — только по явной просьбе пользователя (репозиторий публичный).
- Команды — через `mise run` из корня репо: тесты `mise run //backend:tests`, линтеры `mise run //backend:checks`. Пустой вывод и нулевой exit code = успех; при падении печатается полный лог. Напрямую `go test` не запускать.
- Зависимость пиним точно: `go get github.com/pressly/goose/v3@v3.27.3`. Драйвер SQLite остаётся `modernc.org/sqlite`, cgo-зависимостей не появляется.
- Политика up-only: секций `-- +goose Down` в миграциях нет вообще.
- Коммит-месседжи — английские, императив, без conventional-префиксов (стиль репо: «Return 404 for a missing main track»), в конце трейлер `Co-Authored-By: Claude Opus 5 <noreply@anthropic.com>`.

---

### Task 1: goose-миграции вместо createSchema

**Files:**
- Create: `backend/app/storage/db/migrations/0001_init.sql`
- Create: `backend/app/storage/db/migrate.go`
- Modify: `backend/app/storage/db/sqlite.go` (удалить константы схем на строках 11–55 и `createSchema` на 100–115, заменить вызов на 71–74)
- Modify: `backend/go.mod`, `backend/go.sum` (через `go get`)
- Test: `backend/app/storage/db/sqlite_test.go` (новый файл)

**Interfaces:**
- Consumes: `Engine` (embed `*sqlx.DB`) из `engine.go`; `NewSqliteEngine(fileSource string, logger zerolog.Logger) (*Engine, error)` из `sqlite.go`.
- Produces: приватный метод `(e *Engine) migrate(logger zerolog.Logger) error` в пакете `db`; embed-переменная `migrationsFS embed.FS`. Наружу интерфейс пакета не меняется — Task 2 и 3 опираются на тот же `NewSqliteEngine`.

- [ ] **Step 1: Написать падающий тест**

Создать `backend/app/storage/db/sqlite_test.go`:

```go
package db

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestNewSqliteEngine_AppliesMigrations(t *testing.T) {
	engine, err := NewSqliteEngine(":memory:", zerolog.Nop())
	require.NoError(t, err)
	defer engine.Close()

	// goose ведёт версию схемы в собственной таблице
	var version int
	err = engine.Get(&version, "SELECT MAX(version_id) FROM goose_db_version")
	require.NoError(t, err)
	require.GreaterOrEqual(t, version, 1)

	var tables []string
	err = engine.Select(&tables,
		"SELECT name FROM sqlite_master WHERE type='table' AND name IN ('user','track','workout','avatar') ORDER BY name")
	require.NoError(t, err)
	require.Equal(t, []string{"avatar", "track", "user", "workout"}, tables)
}
```

- [ ] **Step 2: Убедиться, что тест падает**

Run: `mise run //backend:tests`
Expected: FAIL, в логе `TestNewSqliteEngine_AppliesMigrations` с ошибкой вида `no such table: goose_db_version`. Остальные тесты зелёные.

- [ ] **Step 3: Добавить зависимость goose**

```bash
cd backend && go get github.com/pressly/goose/v3@v3.27.3 && go mod tidy
```

- [ ] **Step 4: Создать миграцию 0001 — текущая схема дословно**

Создать `backend/app/storage/db/migrations/0001_init.sql`. Тела таблиц — точная копия констант из `sqlite.go` (важно: `IF NOT EXISTS` сохраняем, на существующей прод-базе миграция должна пройти no-op'ом):

```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS user (
    id TEXT PRIMARY KEY NOT NULL,
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS track (
    id TEXT PRIMARY KEY NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    owner_id TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS workout (
    id TEXT PRIMARY KEY NOT NULL,
    date TEXT NOT NULL,
    track_id TEXT NOT NULL,
    sections TEXT NOT NULL,
    notes TEXT,
    schema_version INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS avatar (
    user_id TEXT PRIMARY KEY NOT NULL,
    mime_type TEXT NOT NULL,
    data BLOB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

- [ ] **Step 5: Создать migrate.go**

Создать `backend/app/storage/db/migrate.go`:

```go
package db

import (
	"context"
	"embed"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// migrate докатывает схему до актуальной версии embed'нутыми миграциями.
func (e *Engine) migrate(logger zerolog.Logger) error {
	dir, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to open migrations dir: %w", err)
	}

	provider, err := goose.NewProvider(goose.DialectSQLite3, e.DB.DB, dir)
	if err != nil {
		return fmt.Errorf("failed to create migration provider: %w", err)
	}

	results, err := provider.Up(context.Background())
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	if len(results) > 0 {
		logger.Info().Msgf("applied %d migrations", len(results))
	}

	return nil
}
```

Примечание: `e.DB` — это embed'нутый `*sqlx.DB`, его поле `.DB` — голый `*sql.DB`, который и нужен goose.

- [ ] **Step 6: Переключить NewSqliteEngine на migrate, удалить createSchema**

В `backend/app/storage/db/sqlite.go`:
1. Удалить целиком блок констант `userSchema`, `trackSchema`, `workoutSchema`, `avatarSchema` (строки 11–55) вместе с обрамляющим `const (...)`.
2. Удалить функцию `createSchema` (строки 100–115).
3. Заменить в `NewSqliteEngine` блок вызова `createSchema`:

```go
	if err := engine.createSchema(); err != nil {
		logger.Error().Err(err).Msg("failed to create sqlite schema")
		return nil, err
	}
```

на:

```go
	if err := engine.migrate(logger); err != nil {
		logger.Error().Err(err).Msg("failed to migrate sqlite schema")
		_ = engine.Close()
		return nil, err
	}
```

Блок `engine.setup()` в этом таске не трогать — им занимается Task 3.

- [ ] **Step 7: Прогнать тесты и линтеры**

Run: `mise run //backend:tests && mise run //backend:checks`
Expected: пустой вывод, exit 0. Весь существующий набор (datastorage, api, auth) проходит — схема теперь приезжает миграцией по тому же пути `NewSqliteEngine`.

- [ ] **Step 8: Commit**

```bash
git add backend/app/storage/db backend/go.mod backend/go.sum
git commit -m "Replace startup CREATE TABLEs with goose migrations

Co-Authored-By: Claude Opus 5 <noreply@anthropic.com>"
```

---

### Task 2: Тесты сценариев существующей базы

**Files:**
- Test: `backend/app/storage/db/sqlite_test.go` (дописать два теста)

**Interfaces:**
- Consumes: `NewSqliteEngine(fileSource string, logger zerolog.Logger) (*Engine, error)` из Task 1.
- Produces: только тесты, нового API нет.

Оба теста — регрессионные стражи поверх уже готового поведения, поэтому ожидание — PASS сразу. Если какой-то упал — стоп, это баг в Task 1 (расхождение DDL миграции со старой схемой или проблема goose с файловой базой), чинить там, а не подгонять тест.

- [ ] **Step 1: Тест идемпотентности рестарта**

Дописать в `backend/app/storage/db/sqlite_test.go` (импорты дополнить: `path/filepath`, `github.com/jmoiron/sqlx`):

```go
func TestNewSqliteEngine_ReopenIsIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app.db")

	engine, err := NewSqliteEngine(path, zerolog.Nop())
	require.NoError(t, err)
	_, err = engine.Exec("INSERT INTO user(id, email, name) VALUES('u1', 'u1@example.com', 'U1')")
	require.NoError(t, err)
	require.NoError(t, engine.Close())

	// повторное открытие = рестарт приложения: миграции уже применены
	engine, err = NewSqliteEngine(path, zerolog.Nop())
	require.NoError(t, err)
	defer engine.Close()

	var count int
	require.NoError(t, engine.Get(&count, "SELECT COUNT(*) FROM user"))
	require.Equal(t, 1, count)
}
```

- [ ] **Step 2: Тест первого запуска на до-goose'овой базе**

Это сценарий первого деплоя на прод: таблицы есть, `goose_db_version` нет. Дописать:

```go
func TestNewSqliteEngine_AdoptsPreGooseDatabase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app.db")

	// база, созданная старым createSchema: таблицы есть, goose_db_version нет
	raw, err := sqlx.Connect("sqlite", path)
	require.NoError(t, err)
	oldSchema := []string{
		`CREATE TABLE IF NOT EXISTS user (
			id TEXT PRIMARY KEY NOT NULL,
			email TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS track (
			id TEXT PRIMARY KEY NOT NULL,
			slug TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			owner_id TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS workout (
			id TEXT PRIMARY KEY NOT NULL,
			date TEXT NOT NULL,
			track_id TEXT NOT NULL,
			sections TEXT NOT NULL,
			notes TEXT,
			schema_version INTEGER NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS avatar (
			user_id TEXT PRIMARY KEY NOT NULL,
			mime_type TEXT NOT NULL,
			data BLOB NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
	}
	for _, ddl := range oldSchema {
		_, err = raw.Exec(ddl)
		require.NoError(t, err)
	}
	_, err = raw.Exec("INSERT INTO user(id, email, name) VALUES('u1', 'u1@example.com', 'U1')")
	require.NoError(t, err)
	require.NoError(t, raw.Close())

	engine, err := NewSqliteEngine(path, zerolog.Nop())
	require.NoError(t, err)
	defer engine.Close()

	// миграция 0001 прошла no-op'ом и зафиксировала версию
	var version int
	require.NoError(t, engine.Get(&version, "SELECT MAX(version_id) FROM goose_db_version"))
	require.Equal(t, 1, version)

	// данные пережили adoption
	var count int
	require.NoError(t, engine.Get(&count, "SELECT COUNT(*) FROM user"))
	require.Equal(t, 1, count)
}
```

- [ ] **Step 3: Прогнать тесты**

Run: `mise run //backend:tests`
Expected: пустой вывод, exit 0 (оба новых теста проходят).

- [ ] **Step 4: Commit**

```bash
git add backend/app/storage/db/sqlite_test.go
git commit -m "Cover restart and pre-goose database adoption with tests

Co-Authored-By: Claude Opus 5 <noreply@anthropic.com>"
```

---

### Task 3: Ошибка setup() должна останавливать старт

**Files:**
- Modify: `backend/app/storage/db/sqlite.go` (блок `engine.setup()` в `NewSqliteEngine`)

**Interfaces:**
- Consumes: `(e *Engine) setup() error` — уже существует, менять его не надо.
- Produces: изменение поведения `NewSqliteEngine`: ошибка прагм теперь возвращается наверх.

Сейчас ошибка `setup()` только логируется и старт продолжается с недонастроенной базой (нет WAL, busy_timeout, foreign_keys). Отдельного юнит-теста нет намеренно: чтобы заставить прагму упасть, пришлось бы подсовывать фейковый драйвер — не окупается для двухстрочного фикса. Страховка — зелёный прогон всего набора.

- [ ] **Step 1: Вернуть ошибку из блока setup**

В `NewSqliteEngine` заменить:

```go
	if err := engine.setup(); err != nil {
		logger.Error().Err(err).Msg("failed to setup sqlite engine")
	}
```

на:

```go
	if err := engine.setup(); err != nil {
		logger.Error().Err(err).Msg("failed to setup sqlite engine")
		return nil, err
	}
```

(Закрывать engine здесь не нужно: `setup()` при ошибке прагмы сам делает `e.Close()`.)

- [ ] **Step 2: Прогнать тесты и линтеры**

Run: `mise run //backend:tests && mise run //backend:checks`
Expected: пустой вывод, exit 0.

- [ ] **Step 3: Commit**

```bash
git add backend/app/storage/db/sqlite.go
git commit -m "Stop startup when sqlite pragmas fail to apply

Co-Authored-By: Claude Opus 5 <noreply@anthropic.com>"
```

---

### Task 4: Документация и бэклог

**Files:**
- Modify: `CLAUDE.md` (секция Backend)
- Modify: `docs/neskuchka.taskpaper` (файл в `.gitignore`, в коммит не попадёт)

**Interfaces:**
- Consumes: ничего из кода.
- Produces: только документация.

- [ ] **Step 1: Записать порядок работы с миграциями в CLAUDE.md**

В секцию `## Backend` (после пункта про доменные ошибки, перед «Хендлеры тонкие») добавить пункт:

```markdown
- Схема БД — goose-миграции в `backend/app/storage/db/migrations/`, embed'ятся в бинарь
  и применяются на старте. Изменение схемы = новый файл `NNNN_<name>.sql` со следующим
  номером и секцией `-- +goose Up`. Down-секций не пишем: откат — восстановлением из бэкапа.
```

- [ ] **Step 2: Заархивировать задачу в бэклоге**

В `docs/neskuchka.taskpaper`: строку `	- Миграция базы @backend @feature` удалить из `BACKLOG:` и добавить в `Archive:` как:

```
	- Миграция базы @backend @feature @done @project(BACKLOG)
```

- [ ] **Step 3: Commit**

```bash
git add CLAUDE.md
git commit -m "Document the goose migrations workflow in CLAUDE.md

Co-Authored-By: Claude Opus 5 <noreply@anthropic.com>"
```
