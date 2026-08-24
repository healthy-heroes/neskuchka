import dayjs from 'dayjs';
import { randomId } from '@mantine/hooks';
import { Workout } from '@/types/domain';

export default function createWorkout(overrides: Partial<Workout> = {}): Workout {
	const date = overrides.Date ?? '2025-01-01';

	return {
		ID: randomId(),
		TrackID: 'track-1',
		Date: date,

		// Состояние считает бэкенд, и фикстура считает его по тем же правилам:
		// тренировка видна с её дня, а править её можно ещё сутки после
		IsPublished: !dayjs(date).isAfter(dayjs(), 'day'),
		CanEdit: !dayjs(date).isBefore(dayjs().subtract(1, 'day'), 'day'),

		Sections: [
			{
				Title: 'Разминка',
				Protocol: {
					Title: '3 раунда',
					Description: '',
				},
				Exercises: [
					{ ExerciseSlug: 'snow-angels', Prescription: ['5'], Name: 'Снежный ангел' },
					{ ExerciseSlug: 'push-ups', Prescription: ['10'], Name: 'Отжимание' },
					{ ExerciseSlug: 'squats', Prescription: ['10'], Name: 'Приседание' },
				],
			},
			{
				Title: 'Комплекс',
				Protocol: {
					Title: 'По минутки 10 мин',
					Description: '20 сек макс повт берпи / 40 сек отжимания с колен',
				},
				Exercises: [
					{ ExerciseSlug: 'push-ups', Prescription: ['10'], Name: 'Отжимание' },
					{ ExerciseSlug: 'squats', Prescription: ['10'], Name: 'Приседание' },
				],
			},
		],
		...overrides,
	};
}

export interface TrackWorkoutsOptions {
	/** Сколько опубликованных тренировок вернуть. */
	count?: number;
	/** Начинать ли с сегодняшней: без неё трек выглядит как «сегодня отдых». */
	includeToday?: boolean;
	/** Сколько запланированных добавить впереди — их видит только владелец. */
	planned?: number;
}

/**
 * Тренировки трека, свежие первыми — по одной раз в два дня.
 *
 * ID и даты детерминированы намеренно. Статус выполнения считается по хешу ID
 * (см. utils/completion), а design-sync сравнивает два рендера одной стори
 * скриншотами: со случайными ID полоса прогресса каждый раз разная и сравнение
 * не сходится никогда.
 */
export function createTrackWorkouts({
	count = 14,
	includeToday = true,
	planned = 0,
}: TrackWorkoutsOptions = {}): Array<Workout> {
	const today = dayjs().startOf('day');

	// Запланированные идут первыми: список везде отсортирован свежими вверх
	const upcoming = Array.from({ length: planned }, (_, index) => {
		const daysAhead = (planned - index) * 3;

		return createWorkout({
			ID: `workout-planned-${planned - index}`,
			Date: today.add(daysAhead, 'day').format('YYYY-MM-DD'),
		});
	});

	const published = Array.from({ length: count }, (_, index) => {
		const offset = includeToday ? index * 2 : index * 2 + 1;

		return createWorkout({
			ID: `workout-${index + 1}`,
			Date: today.subtract(offset, 'day').format('YYYY-MM-DD'),
		});
	});

	return [...upcoming, ...published];
}
