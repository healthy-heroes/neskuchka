import { StoryPreview } from '../StoryBook/StoryPreview';
import { Workouts } from './Workouts';

export default {
	title: 'Workouts',
	component: Workouts,
};

export function Default() {
	return (
		<StoryPreview>
			<Workouts />
		</StoryPreview>
	);
}
