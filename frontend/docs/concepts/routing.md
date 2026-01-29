# Routing Architecture: Thin Routes

## Концепция

Роуты в проекте — "тонкие". Они только связывают URL с компонентом, без загрузки данных и проверок бизнес-логики.

## Роуты отвечают за:

1. **Mapping URL → Component** — какой компонент рендерить
2. **Извлечение params** — `useParams()` для передачи в компонент
3. **Структура вложенности** — layouts через `<Outlet />`
4. **Auth guards** — оборачивание в `<RequireAuth>` для protected routes

## Роуты НЕ отвечают за:

- ❌ Загрузку данных (никаких `loader`)
- ❌ Проверку авторизации через `beforeLoad` с redirect
- ❌ Валидацию params
- ❌ Бизнес-логику

## Компоненты (Pages) отвечают за:

- ✅ Загрузку данных — `useQuery()` внутри компонента
- ✅ Специфичные проверки прав — `useIsOwner()` + `<Navigate />` (когда нужны данные)
- ✅ Loading/Error states
- ✅ Бизнес-логику отображения

---

## RequireAuth

Guard component для проверки авторизации. Используется **в route component**.

### Когда использовать

| Сценарий | Решение |
|----------|---------|
| Страница только для авторизованных | `<RequireAuth>` |
| Страница только для гостей | `<RequireAuth guestOnly>` |
| Проверка ownership (нужны данные) | `useIsOwner()` внутри компонента |

### Гарантии RequireAuth

Компонент внутри `RequireAuth` получает гарантии:
- ✅ Auth запрос завершён (`isLoading = false`)
- ✅ Юзер авторизован (`isAuthenticated = true`)
- ✅ `user !== null`

Это решает проблему race condition между auth и data loading.

---

## Примеры

### Protected route (только для авторизованных)

```tsx
// routes/workouts.new.tsx
function RouteComponent() {
  return (
    <RequireAuth loadingComponent={<PageSkeleton />}>
      <WorkoutCreate />
    </RequireAuth>
  )
}
```

### Guest-only route (только для гостей)

```tsx
// routes/login.tsx
function RouteComponent() {
  return (
    <RequireAuth loadingComponent={<PageSkeleton />} guestOnly>
      <AuthPage />
    </RequireAuth>
  )
}
```

### Protected + Owner check

```tsx
// routes/workouts.$workoutId_.edit.tsx
function RouteComponent() {
  const { workoutId } = Route.useParams()
  return (
    <RequireAuth loadingComponent={<PageSkeleton />}>
      <WorkoutEdit workoutId={workoutId} />
    </RequireAuth>
  )
}
```

```tsx
// components/WorkoutEdit/WorkoutEdit.tsx
function WorkoutEdit({ workoutId }: { workoutId: string }) {
  const { workouts } = useApi()
  const { data, isLoading } = useQuery(workouts.getWorkoutQuery(workoutId))
  const isOwner = useIsOwner(data?.OwnerID)

  // Auth уже гарантирован RequireAuth, ждём только данные
  if (isLoading) {
    return <PageSkeleton />
  }

  // Owner check — после загрузки данных
  if (!isOwner) {
    return <Navigate to="/workouts/$workoutId" params={{ workoutId }} />
  }

  return <WorkoutForm data={data.Workout} ... />
}
```

### Public route (без проверок)

```tsx
// routes/workouts.$workoutId.tsx
function RouteComponent() {
  const { workoutId } = Route.useParams()
  return <WorkoutView workoutId={workoutId} />  // Без RequireAuth
}
```

## Trade-offs

### Плюсы

- **Простота** — роуты тривиальны, легко понять
- **Тестируемость** — компоненты тестируются отдельно от роутинга
- **Гибкость** — компонент полностью контролирует свой lifecycle
- **Storybook** — компоненты работают без роутера

### Минусы

- Нет prefetch данных при hover на ссылку (router loaders умеют это)
- Loading state показывается ПОСЛЕ навигации (а не до)

## См. также

- `@/auth/RequireAuth` — guard component для auth проверок
- `@/auth/hooks` — хуки `useAuth()`, `useIsOwner()`
- `@/components/PageSkeleton` — компонент для loading state
