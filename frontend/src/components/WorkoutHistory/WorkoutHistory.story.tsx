import { createTrackWorkouts } from '@/api/fixtures/workout';
import { StoryPreview } from '../StoryBook/StoryPreview';
import { WorkoutHistory } from './WorkoutHistory';

export default {
	title: 'WorkoutHistory',
};

const workouts = createTrackWorkouts();

export function Default() {
	return (
		<StoryPreview isPage>
			<WorkoutHistory workouts={workouts} />
		</StoryPreview>
	);
}

/** Тренировка из карточки в фокусе в историю не попадает — здесь видно, что её нет. */
export function WithoutFeatured() {
	return (
		<StoryPreview isPage>
			<WorkoutHistory workouts={workouts} featured={workouts[0]} />
		</StoryPreview>
	);
}
