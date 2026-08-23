## Wrapping — required

Two providers, outermost first. Nothing renders correctly without them:

```jsx
const { DatesProvider, MantineProvider, theme } = window.NeskuchkaDS;

<DatesProvider settings={{ locale: 'ru', firstDayOfWeek: 0, weekendDays: [0] }}>
  <MantineProvider theme={theme}>{children}</MantineProvider>
</DatesProvider>
```

`theme` is a bundle export. Outside `MantineProvider` components throw
"MantineProvider was not found", and the brand palettes (`copper`, `slate`) do not
exist: the compiled stylesheet ships only Mantine's stock palettes, and the provider
injects ours at runtime from `theme`.

## The styling idiom

This is Mantine 9. Do **not** write your own classes for layout, spacing, or
typography — the library expresses all three.

- Layout: `Group` (row), `Stack` (column), `SimpleGrid`, `Container` (working width is
  fixed in the theme — write `<Container>` with no size prop).
- Typography: `Title order={1..6}` — Oswald, uppercase and letter-spacing all come from
  the theme, never restate them. Body copy: `Text` with `fz`, `fw`, `c`.
- Spacing: the `m`/`p` props and their variants (`mt`, `px`, …), `gap` on Group/Stack.
- Resetting a button: `UnstyledButton` — it also carries the focus ring.

**Scales are named, never pixels.** `fz` xs·sm·md·lg·xl = 12·14·16·18·20 px, spacing =
10·12·16·20·32, radius = 4·6·8·14·24. In hand-written CSS use the variables:
`var(--mantine-spacing-lg)`, `var(--mantine-radius-lg)`,
`var(--mantine-font-family-headings)`.

**Colors are theme keys, never hex.** Three families:

| Family | Role |
|---|---|
| `copper` | actions, accents, "today" — the primary color |
| `slate` | dark panels, the "done" status |
| `gray` | warm neutrals; overrides Mantine's stock gray |

Written as `c="gray.7"`, `bg="slate.7"`, `c="copper.6"`. Avoid `c="dimmed"` — Mantine
maps it to gray-6, but secondary text in this system is `gray.7`.

**Theme defaults — do not repeat them as props**: `Button` (radius md, weight 600),
`Card` (radius lg, bordered), `Badge` (pill radius), `Title` (uppercase + tracking),
`Container` (width 1240).

## Where the truth lives

- `styles.css` → `_ds_bundle.css` — the compiled stylesheet. Grep it for a token before
  inventing one.
- `components/<group>/<Name>/<Name>.prompt.md` — props plus real usage examples.
- `components/<group>/<Name>/<Name>.d.ts` — the prop contract.

The auto-generated one-line blurbs in the component index below are empty: the source
JSDoc is Russian and the generator strips non-Latin characters. Use this index instead.

## What is here

**foundation** — `ColorRamps`, `TypeScale`, `SizeScale`. Reference cards showing the
palettes, the type scale and the spacing/radius steps. Read them; don't compose with them.

**primitives** — `ButtonSpecimen`, `BadgeSpecimen`, `SurfaceSpecimen`: the themed
defaults of Button, Badge, Card and Paper. Same — they are specimens, not building blocks.

**Application components** — `Header`, `Logo`, `RouteLink`, `PageSkeleton`;
`TrackProgress` (the 30-day bar), `FeaturedWorkout` (today's session), `WorkoutSections`
with `ExerciseRow`, `WorkoutHistory`, `WorkoutCardSkeleton`, `WorkoutView`, `WorkoutForm`,
`Workouts` (the whole track page), `LandingPage`, `LoginForm`, `LoginConfirm`.

`ExerciseRow` is `display: contents` — it lays out into `WorkoutSections`' grid and
collapses on its own. Never use it outside that parent.

Interface copy is Russian. Keep new copy Russian too.

## Idiomatic snippet

```jsx
const { MantineProvider, DatesProvider, theme, Container, Stack, Title, Text, Card } =
  window.NeskuchkaDS;

<DatesProvider settings={{ locale: 'ru', firstDayOfWeek: 0, weekendDays: [0] }}>
  <MantineProvider theme={theme}>
    <Container>
      <Stack gap="lg">
        <Title order={2}>Тренировка дня</Title>
        <Card p="lg">
          <Text c="gray.7">Разминка · 3 раунда</Text>
        </Card>
      </Stack>
    </Container>
  </MantineProvider>
</DatesProvider>
```
