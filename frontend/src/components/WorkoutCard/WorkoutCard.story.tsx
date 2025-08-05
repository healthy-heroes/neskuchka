import { Workout } from '@/types/domain';
import { StoryPreview } from '../StoryBook/StoryPreview';
import { WorkoutCard } from './WorkoutCard';
import { WorkoutCardSkeleton } from './WorkoutCardSkeleton';

export default {
	title: 'WorkoutCard',
};

const workout: Workout = {
	ID: 1,
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

export function Default() {
	return (
		<StoryPreview>
			<WorkoutCard workout={workout} />
		</StoryPreview>
	);
}

export function Skeleton() {
	return (
		<StoryPreview>
			<WorkoutCardSkeleton />
		</StoryPreview>
	);
}
