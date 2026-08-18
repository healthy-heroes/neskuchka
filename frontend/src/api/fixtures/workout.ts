import { randomId } from '@mantine/hooks';
import { Workout } from '@/types/domain';

export default function createWorkout(): Workout {
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
	};
}
