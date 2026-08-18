import createWorkout from '@/api/fixtures/workout';
import { StoryPreview } from '../StoryBook/StoryPreview';
import { WorkoutCard } from './WorkoutCard';
import { WorkoutCardSkeleton } from './WorkoutCardSkeleton';

export default {
	title: 'WorkoutCard',
};

export function Default() {
	return (
		<StoryPreview>
			<WorkoutCard workout={createWorkout()} />
		</StoryPreview>
	);
}

export function Skeleton() {
	return (
		<StoryPreview>
			<WorkoutCardSkeleton />
		</StoryPreview>
	);
}
