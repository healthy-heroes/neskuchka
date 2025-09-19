import createWorkout from '@/api/fixtures/workout';
import { createMock, createMockApiService } from '@/api/service-mock';
import { TrackWorkout } from '@/api/types';
import { StoryPreview } from '../StoryBook/StoryPreview';
import { WorkoutView } from './WorkoutView';

export default {
	title: 'WorkoutView',
};

const apiService = createMockApiService({
	getWorkout: createMock<TrackWorkout>({
		Workout: createWorkout(),
	}),
});

export const Default = () => {
	return (
		<StoryPreview apiService={apiService}>
			<WorkoutView workoutId="1" />
		</StoryPreview>
	);
};
