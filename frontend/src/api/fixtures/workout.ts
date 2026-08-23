import dayjs from 'dayjs';
import { randomId } from '@mantine/hooks';
import { Workout } from '@/types/domain';

export default function createWorkout(overrides: Partial<Workout> = {}): Workout {
	return {
		ID: randomId(),
		TrackID: 'track-1',
		Date: '2025-01-01',
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
	/** Сколько тренировок вернуть. */
	count?: number;
	/** Начинать ли с сегодняшней: без неё трек выглядит как «сегодня отдых». */
	includeToday?: boolean;
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
}: TrackWorkoutsOptions = {}): Array<Workout> {
	const today = dayjs().startOf('day');

	return Array.from({ length: count }, (_, index) => {
		const offset = includeToday ? index * 2 : index * 2 + 1;

		return createWorkout({
			ID: `workout-${index + 1}`,
			Date: today.subtract(offset, 'day').format('YYYY-MM-DD'),
		});
	});
}
