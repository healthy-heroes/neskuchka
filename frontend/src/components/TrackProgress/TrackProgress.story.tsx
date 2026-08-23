import { createTrackWorkouts } from '@/api/fixtures/workout';
import { StoryPreview } from '../StoryBook/StoryPreview';
import { TrackProgress } from './TrackProgress';
import { TrackProgressSkeleton } from './TrackProgressSkeleton';

export default {
	title: 'TrackProgress',
};

const workouts = createTrackWorkouts();

export function Default() {
	return (
		<StoryPreview>
			<TrackProgress workouts={workouts} />
		</StoryPreview>
	);
}

export function WithTrackTotal() {
	return (
		<StoryPreview>
			<TrackProgress workouts={workouts} total={128} />
		</StoryPreview>
	);
}

export function Skeleton() {
	return (
		<StoryPreview>
			<TrackProgressSkeleton />
		</StoryPreview>
	);
}
