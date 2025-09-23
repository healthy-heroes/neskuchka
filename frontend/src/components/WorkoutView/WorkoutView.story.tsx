import { ApiServiceMock } from '@/api/fixtures/api';
import createWorkout from '@/api/fixtures/workout';
import { StoryPreview } from '../StoryBook/StoryPreview';
import { WorkoutView } from './WorkoutView';

export default {
	title: 'WorkoutView',
};

const workout = createWorkout();

const apiService = new ApiServiceMock();
apiService.workouts.getWorkoutQuery = () => {
	return {
		queryKey: ['workout', workout.ID],
		queryFn: () => {
			return Promise.resolve({ Workout: workout });
		},
	};
};

export const Default = () => {
	return (
		<StoryPreview apiService={apiService}>
			<WorkoutView workoutId={workout.ID} />
		</StoryPreview>
	);
};
