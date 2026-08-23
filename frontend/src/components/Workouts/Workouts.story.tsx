import { createApiServiceMock } from '@/api/fixtures/api';
import { mockTrack } from '@/api/fixtures/track';
import { createTrackWorkouts } from '@/api/fixtures/workout';
import { TrackData, TrackWorkoutsData, WorkoutsKeys } from '@/api/services/workouts';
import { StoryPreview } from '../StoryBook/StoryPreview';
import { Workouts } from './Workouts';

export default {
	title: 'Workouts',
};

// Восемь, а не четырнадцать: страница трека целиком должна влезать в кадр
// карточки дизайн-системы (потолок захвата — 2000px), да и на карточке
// история в полгода нечитаема.
const workouts = createTrackWorkouts({ count: 8 });

const apiService = createApiServiceMock({
	workouts: {
		getMainTrackQuery: () => ({
			queryKey: WorkoutsKeys.track(),
			queryFn: () => Promise.resolve({ data: { Track: mockTrack.Track, IsOwner: false } }),
			select: (response: { data: TrackData }): TrackData => response.data,
		}),
		getMainTrackWorkoutsQuery: () => ({
			queryKey: WorkoutsKeys.workouts(),
			queryFn: () => Promise.resolve({ data: { Workouts: workouts } }),
			select: (response: { data: TrackWorkoutsData }): TrackWorkoutsData => response.data,
		}),
	},
});

export function Default() {
	return (
		<StoryPreview isPage apiService={apiService}>
			<Workouts />
		</StoryPreview>
	);
}
