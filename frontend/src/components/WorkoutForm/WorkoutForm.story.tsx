import createWorkout from '@/api/fixtures/workout';
import { StoryPreview } from '../StoryBook/StoryPreview';
import { WorkoutForm } from './WorkoutForm';

export default {
	title: 'WorkoutForm',
};

const props = {
	trackId: 'track-1',
	onSubmit: () => {},
	error: null,
	onCancel: () => {},
};

export function Default() {
	return (
		<StoryPreview>
			<WorkoutForm {...props} />
		</StoryPreview>
	);
}

export function WithInitialValues() {
	const workout = createWorkout();

	return (
		<StoryPreview>
			<WorkoutForm {...props} data={workout} />
		</StoryPreview>
	);
}

export function WithSubmitting() {
	const workout = createWorkout();

	const localProps = {
		...props,
		isSubmitting: true,
	};

	return (
		<StoryPreview>
			<WorkoutForm {...localProps} data={workout} />
		</StoryPreview>
	);
}
