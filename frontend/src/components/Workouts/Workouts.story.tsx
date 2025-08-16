import createWorkout from '@/api/fixtures/workout';
import { createMock, createMockApiService } from '@/api/service-mock';
import { TrackWorkouts } from '@/api/types';
import { StoryPreview } from '../StoryBook/StoryPreview';
import { Workouts } from './Workouts';

export default {
	title: 'Workouts',
	component: Workouts,
};

const apiService = createMockApiService({
	getMainTrackWorkouts: createMock<TrackWorkouts>({
		Workouts: [createWorkout()],
	}),
});

export function Default() {
	return (
		<StoryPreview apiService={apiService}>
			<Workouts />
		</StoryPreview>
	);
}
