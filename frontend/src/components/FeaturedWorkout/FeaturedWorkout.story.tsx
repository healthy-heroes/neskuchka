import { createTrackWorkouts } from '@/api/fixtures/workout';
import { StoryPreview } from '../StoryBook/StoryPreview';
import { FeaturedWorkout } from './FeaturedWorkout';
import { FeaturedWorkoutSkeleton } from './FeaturedWorkoutSkeleton';

export default {
	title: 'FeaturedWorkout',
};

const [today, earlier] = createTrackWorkouts({ count: 2 });

export function Today() {
	return (
		<StoryPreview isPage>
			<FeaturedWorkout workout={today} />
		</StoryPreview>
	);
}

export function LastPublished() {
	return (
		<StoryPreview isPage>
			<FeaturedWorkout workout={earlier} />
		</StoryPreview>
	);
}

export function Skeleton() {
	return (
		<StoryPreview isPage>
			<FeaturedWorkoutSkeleton />
		</StoryPreview>
	);
}
