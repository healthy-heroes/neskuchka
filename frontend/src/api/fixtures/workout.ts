import { randomId } from '@mantine/hooks';
import { Workout } from '@/types/domain';

export default function createWorkout(): Workout {
	return {
		ID: randomId(),
		Date: '2025-01-01',
		Sections: [
			{
				Title: 'Разминка',
				Protocol: {
					Title: '3 раунда',
					Description: '',
				},
				Exercises: [
					{ ExerciseSlug: 'snow-angels', Description: '5 снежных ангелов' },
					{ ExerciseSlug: 'push-ups', Description: '10 отжиманий' },
					{ ExerciseSlug: 'squats', Description: '10 приседаний' },
				],
			},
			{
				Title: 'Комплекс',
				Protocol: {
					Title: 'По минутки 10 мин',
					Description: '20 сек макс повт берпи / 40 сек отжимания с колен',
				},
				Exercises: [
					{ ExerciseSlug: 'push-ups', Description: '10 отжиманий' },
					{ ExerciseSlug: 'squats', Description: '10 приседаний' },
				],
			},
		],
	};
}
