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
// routes/profile.tsx
function RouteComponent() {
  return (
    <RequireAuth loadingComponent={<PageSkeleton />}>
      <ProfilePage />
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

### Protected + Owner check (TrackOwnerOnly)

```tsx
// routes/workouts.$workoutId_.edit.tsx
function RouteComponent() {
  const { workoutId } = Route.useParams()
  
  return (
    <TrackOwnerOnly 
      loadingComponent={<PageSkeleton hideHeader />} 
      redirectTo={`/workouts/${workoutId}`}
    >
      <WorkoutEdit workoutId={workoutId} />
    </TrackOwnerOnly>
  )
}
```

```tsx
// components/WorkoutEdit/WorkoutEdit.tsx
function WorkoutEdit({ workoutId }: { workoutId: string }) {
  // Owner check уже гарантирован TrackOwnerOnly
  // Компонент загружает только данные workout
  const { workouts } = useApi()
  const { data, isLoading } = useQuery(workouts.getWorkoutQuery(workoutId))

  if (isLoading) {
    return <PageSkeleton />
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

---

## TrackOwnerOnly

Guard component для проверки что текущий пользователь — владелец трека. Используется **в route component** для owner-only страниц.

### Когда использовать

| Сценарий | Решение |
|----------|---------|
| Создание workout | `<TrackOwnerOnly>` |
| Редактирование workout | `<TrackOwnerOnly>` |
| Удаление workout | `<TrackOwnerOnly>` |

### Как работает

1. Загружает track через `getMainTrackQuery()`
2. Если `isPending` — показывает `loadingComponent`
3. Если `!IsOwner` — редирект на `redirectTo`
4. Если owner — рендерит children

### Гарантии TrackOwnerOnly

Компонент внутри `TrackOwnerOnly` получает гарантии:
- ✅ Track загружен и в кеше
- ✅ Текущий пользователь — владелец трека

### Пример

```tsx
// routes/workouts.new.tsx
function RouteComponent() {
  return (
    <TrackOwnerOnly 
      loadingComponent={<PageSkeleton hideHeader />} 
      redirectTo="/workouts"
    >
      <WorkoutCreate />
    </TrackOwnerOnly>
  )
}
```

```tsx
// routes/workouts.$workoutId_.edit.tsx
function RouteComponent() {
  const { workoutId } = Route.useParams()
  
  return (
    <TrackOwnerOnly 
      loadingComponent={<PageSkeleton hideHeader />} 
      redirectTo={`/workouts/${workoutId}`}
    >
      <WorkoutEdit workoutId={workoutId} />
    </TrackOwnerOnly>
  )
}
```

---

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
- `@/guards/TrackOwnerOnly` — guard component для owner-only страниц
- `@/auth/hooks` — хуки `useAuth()`, `useIsOwner()`
- `@/components/PageSkeleton` — компонент для loading state
