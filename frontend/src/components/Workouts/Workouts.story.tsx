import { ApiMock } from '@/api/fixtures/api';
import createWorkout from '@/api/fixtures/workout';
import { ApiQueries } from '@/api/queries';
import { StoryPreview } from '../StoryBook/StoryPreview';
import { Workouts } from './Workouts';

export default {
	title: 'Workouts',
};

const queries = new ApiQueries(new ApiMock());
queries.workouts.getMainTrackWorkoutsQuery = () => {
	return {
		queryKey: ['workouts'],
		queryFn: () => {
			return Promise.resolve({ Workouts: [createWorkout(), createWorkout()] });
		},
	};
};

export function Default() {
	return (
		<StoryPreview queries={queries}>
			<Workouts />
		</StoryPreview>
	);
}
