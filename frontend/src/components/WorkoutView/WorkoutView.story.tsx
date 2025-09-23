import { ApiMock } from '@/api/fixtures/api';
import createWorkout from '@/api/fixtures/workout';
import { ApiQueries } from '@/api/queries';
import { StoryPreview } from '../StoryBook/StoryPreview';
import { WorkoutView } from './WorkoutView';

export default {
	title: 'WorkoutView',
};

const workout = createWorkout();

const queries = new ApiQueries(new ApiMock());
queries.workouts.getWorkoutQuery = () => {
	return {
		queryKey: ['workout', workout.ID],
		queryFn: () => {
			return Promise.resolve({ Workout: workout });
		},
	};
};

export const Default = () => {
	return (
		<StoryPreview queries={queries}>
			<WorkoutView workoutId={workout.ID} />
		</StoryPreview>
	);
};
