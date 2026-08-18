import { StoryPreview } from '../StoryBook/StoryPreview';
import { WorkoutCardSkeleton } from './WorkoutCardSkeleton';

export default {
	title: 'WorkoutCardSkeleton',
};

export function Default() {
	return (
		<StoryPreview>
			<WorkoutCardSkeleton />
		</StoryPreview>
	);
}
