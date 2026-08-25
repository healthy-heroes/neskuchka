/**
 * Вход дизайн-системы: то, что уезжает в claude.ai/design.
 *
 * Приложение этот файл не импортирует — он собирается отдельной сборкой
 * (`vite.ds.config.mjs` → `dist-ds/`) в один ESM-модуль, который конвертер
 * design-sync превращает в `window.NeskuchkaDS`. Стори рендерятся против этого
 * бандла, и агент на канвасе собирает экраны из него же.
 *
 * Правило простое: что экспортировано здесь — то существует в дизайн-системе.
 * Компонент без экспорта агент не увидит, даже если у него есть стори.
 *
 * Поверхность Mantine реэкспортируется целиком намеренно: MantineProvider должен
 * приезжать из бандла, иначе контекст темы в превью и в компонентах будет разный.
 */

import '@fontsource-variable/commissioner';
import '@fontsource-variable/oswald';
import '@mantine/core/styles.css';
import '@mantine/tiptap/styles.css';
import '@mantine/dates/styles.css';
import '@/App.css';

import dayjs from 'dayjs';

import 'dayjs/locale/ru';

// Как в App.tsx: даты в дизайн-системе должны быть русскими
dayjs.locale('ru');

export * from '@mantine/core';
export { DatesProvider } from '@mantine/dates';

/** Тема — источник правды по цветам, типографике и шкалам. */
export { theme } from '@/theme';

/* Образцы темы: рампы, шкалы, примитивы в наших умолчаниях. */
export { BadgeSpecimen } from './specimens/BadgeSpecimen';
export { ButtonSpecimen } from './specimens/ButtonSpecimen';
export { ColorRamps } from './specimens/ColorRamps';
export { SizeScale } from './specimens/SizeScale';
export { SurfaceSpecimen } from './specimens/SurfaceSpecimen';
export { TypeScale } from './specimens/TypeScale';

/* Примитивы приложения. */
export { Logo } from '@/components/Logo/Logo';
export { LogoIcon } from '@/components/Logo/LogoIcon';
export { RouteLink } from '@/components/RouteLink/RouteLink';

/* Треки и тренировки. */
export { ExerciseRow } from '@/components/ExerciseRow/ExerciseRow';
export { FeaturedWorkout } from '@/components/FeaturedWorkout/FeaturedWorkout';
export { FeaturedWorkoutSkeleton } from '@/components/FeaturedWorkout/FeaturedWorkoutSkeleton';
export { TrackAbout } from '@/pages/TrackManage/TrackAbout/TrackAbout';
export { TrackProgress } from '@/components/TrackProgress/TrackProgress';
export { TrackProgressSkeleton } from '@/components/TrackProgress/TrackProgressSkeleton';
export { WorkoutCardSkeleton } from '@/components/WorkoutCard/WorkoutCardSkeleton';
export { WorkoutForm } from '@/components/WorkoutForm/WorkoutForm';
export { WorkoutHistory } from '@/components/WorkoutHistory/WorkoutHistory';
export { WorkoutRows } from '@/pages/TrackManage/WorkoutRows/WorkoutRows';
export { Workouts } from '@/components/Workouts/Workouts';
export { WorkoutSections } from '@/components/WorkoutSections/WorkoutSections';
export { WorkoutSectionsSkeleton } from '@/components/WorkoutSections/WorkoutSectionsSkeleton';
export { WorkoutView } from '@/components/WorkoutView/WorkoutView';

/* Каркас страниц и авторизация. */
export { Header } from '@/components/Header/Header';
export { LoginConfirm } from '@/components/LoginConfirm/LoginConfirm';
export { LoginForm } from '@/components/Login/LoginForm';
export { LoginSuccess } from '@/components/Login/LoginSuccess';
export { PageSkeleton } from '@/components/PageSkeleton/PageSkeleton';
export { LandingPage } from '@/pages/Landing/Landing.page';

/**
 * Обвязка стори. Живёт в бандле, потому что стори импортируют её так же, как
 * компоненты, и вторая копия провайдеров развалила бы контекст.
 */
export { StoryPreview } from '@/components/StoryBook/StoryPreview';
export { createApiServiceMock } from '@/api/fixtures/api';
export { createAuthServiceMock } from '@/api/fixtures/auth';
export { mockTrack } from '@/api/fixtures/track';
export { createUserServiceMock, mockUser } from '@/api/fixtures/user';
export { default as createWorkout, createTrackWorkouts } from '@/api/fixtures/workout';
