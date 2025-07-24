import { StoryPreview } from '../StoryBook/StoryPreview';
import { WorkoutCard } from './WorkoutCard';

export default {
	title: 'WorkoutCard',
};

export function Default() {
	return (
		<StoryPreview>
			<WorkoutCard />
		</StoryPreview>
	);
}
