import { createApiServiceMock } from '@/api/fixtures/api';
import createWorkout from '@/api/fixtures/workout';
import { StoryPreview } from '../StoryBook/StoryPreview';
import { WorkoutView } from './WorkoutView';

export default {
	title: 'WorkoutView',
};

const workout = createWorkout();

const apiService = createApiServiceMock({
	workouts: {
		getWorkoutQuery: () => ({
			queryKey: ['workout', workout.ID],
			queryFn: () => Promise.resolve({ data: { Workout: workout } }),
			select: (response) => response.data,
		}),
	},
});

export const Default = () => {
	return (
		<StoryPreview apiService={apiService}>
			<WorkoutView workoutId={workout.ID} />
		</StoryPreview>
	);
};
