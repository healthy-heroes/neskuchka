import createWorkout from '@/api/fixtures/workout';
import { createMock, createMockApiService } from '@/api/service-mock';
import { TrackWorkouts } from '@/api/types';
import { StoryPreview } from '@/components/StoryBook/StoryPreview';
import { MainTrackPage } from './MainTrack.page';

export default {
	title: 'Pages/MainTrack',
};

const apiService = createMockApiService({
	getMainTrackWorkouts: createMock<TrackWorkouts>({
		Workouts: [createWorkout()],
	}),
});

export function Default() {
	return (
		<StoryPreview isPage apiService={apiService}>
			<MainTrackPage />
		</StoryPreview>
	);
}
