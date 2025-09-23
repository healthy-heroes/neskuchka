import { ApiServiceMock } from '@/api/fixtures/api';
import createWorkout from '@/api/fixtures/workout';
import { StoryPreview } from '../StoryBook/StoryPreview';
import { Workouts } from './Workouts';

export default {
	title: 'Workouts',
};

const apiService = new ApiServiceMock();
apiService.workouts.getMainTrackWorkoutsQuery = () => {
	return {
		queryKey: ['workouts'],
		queryFn: () => {
			return Promise.resolve({ Workouts: [createWorkout(), createWorkout()] });
		},
	};
};

export function Default() {
	return (
		<StoryPreview apiService={apiService}>
			<Workouts />
		</StoryPreview>
	);
}
