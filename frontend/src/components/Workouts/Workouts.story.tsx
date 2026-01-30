import { createApiServiceMock } from '@/api/fixtures/api';
import createWorkout from '@/api/fixtures/workout';
import { TrackWorkoutsData } from '@/api/services/workouts';
import { StoryPreview } from '../StoryBook/StoryPreview';
import { Workouts } from './Workouts';

export default {
	title: 'Workouts',
};

const apiService = createApiServiceMock({
	workouts: {
		getMainTrackWorkoutsQuery: () => ({
			queryKey: ['workouts'],
			queryFn: () => Promise.resolve({ data: { Workouts: [createWorkout(), createWorkout()] } }),
			select: (response: { data: TrackWorkoutsData }): TrackWorkoutsData => response.data,
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
