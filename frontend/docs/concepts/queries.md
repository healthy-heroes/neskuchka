# React Query: Паттерны и Use Cases

## Статусы запроса в TanStack Query v5

### isPending vs isLoading

В TanStack Query v5 изменилась семантика статусов:

| Статус | Описание |
|--------|----------|
| `isPending` | Данных ещё нет (`status === 'pending'`) |
| `isFetching` | Запрос выполняется прямо сейчас |
| `isLoading` | `isPending && isFetching` — первая загрузка |

### Когда что использовать

| Сценарий | Статус |
|----------|--------|
| Показать skeleton при первой загрузке | `isPending` или `isLoading` |
| Показать spinner при рефетче (данные уже есть) | `isFetching && !isPending` |
| Заблокировать UI пока нет данных | `isPending` |

### Пример

```tsx
const { data, isPending } = useQuery(api.workouts.getMainTrackQuery());

if (isPending) {
  return <PageSkeleton />;
}
```

---

## Проверка успешности запроса

### Варианты проверки data

После проверки `isPending`, есть несколько способов убедиться что данные есть:

#### 1. Проверка `!data` (рекомендуется)

```tsx
const { data, isPending } = useQuery(...);

if (isPending) return <PageSkeleton />;
if (!data) return <Navigate to="/" />;

// data: TData ✅
return <Component data={data} />;
```

**Плюсы:**
- TypeScript сужает тип
- Обрабатывает случай пустого ответа от API
- Обрабатывает случай ошибки (data будет undefined)

#### 2. Проверка `isError`

```tsx
const { data, isPending, isError } = useQuery(...);

if (isPending) return <PageSkeleton />;
if (isError) return <ErrorPage />;

// ⚠️ TypeScript НЕ сужает тип data автоматически
// data всё ещё TData | undefined
```

**Важно:** TypeScript не понимает что после `!isError` данные гарантированно есть.

#### 3. Проверка `isSuccess` (TypeScript сужает тип)

```tsx
const { data, isPending, isSuccess } = useQuery(...);

if (isPending) return <PageSkeleton />;

if (isSuccess) {
  // data: TData ✅
  return <Component data={data} />;
}

return <ErrorPage />;
```

#### 4. Проверка `status` (дискриминант)

```tsx
const { data, status } = useQuery(...);

if (status === 'pending') return <PageSkeleton />;
if (status === 'error') return <ErrorPage />;

// status === 'success', data: TData ✅
return <Component data={data} />;
```

---

## Рекомендуемый паттерн для страниц

```tsx
export function MyPage() {
  const api = useApi();
  const { data, isPending } = useQuery(api.someResource.getQuery());

  if (isPending) {
    return <PageSkeleton />;
  }

  if (!data) {
    return <Navigate to="/" />;
  }

  return <MyComponent data={data} />;
}
```

### Почему этот паттерн

1. **`isPending`** — показываем skeleton пока нет данных
2. **`!data`** — универсальная проверка:
   - Ошибка запроса → `data` будет `undefined`
   - Пустой ответ API → `data` будет `null` или `undefined`
   - TypeScript сужает тип

---

## Обработка ошибок

### Простой случай — редирект

```tsx
if (!data) {
  return <Navigate to="/" />;
}
```

### С отображением ошибки

```tsx
const { data, isPending, isError, error } = useQuery(...);

if (isPending) return <PageSkeleton />;

if (isError) {
  return <ErrorMessage error={error} />;
}

if (!data) {
  return <Navigate to="/" />;
}
```

---

## См. также

- [TanStack Query v5 Migration Guide](https://tanstack.com/query/latest/docs/framework/react/guides/migrating-to-v5)
- `@/api/hooks` — хук `useApi()` для доступа к сервисам
- `@/components/PageSkeleton` — компонент для loading state
