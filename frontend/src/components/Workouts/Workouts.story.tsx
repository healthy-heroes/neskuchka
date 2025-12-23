import { createApiServiceMock } from '@/api/fixtures/api';
import createWorkout from '@/api/fixtures/workout';
import { StoryPreview } from '../StoryBook/StoryPreview';
import { Workouts } from './Workouts';

export default {
	title: 'Workouts',
};

const apiService = createApiServiceMock({
	workouts: {
		getMainTrackWorkoutsQuery: () => ({
			queryKey: ['workouts'],
			queryFn: () => Promise.resolve({ Workouts: [createWorkout(), createWorkout()] }),
		}),
	},
});

export function Default() {
	return (
		<StoryPreview apiService={apiService}>
			<Workouts />
		</StoryPreview>
	);
}
