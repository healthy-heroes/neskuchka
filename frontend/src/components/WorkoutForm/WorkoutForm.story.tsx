import createWorkout from '@/api/fixtures/workout';
import { StoryPreview } from '../StoryBook/StoryPreview';
import { WorkoutForm } from './WorkoutForm';

export default {
	title: 'WorkoutForm',
};

export function Default() {
	return (
		<StoryPreview>
			<WorkoutForm onSubmit={() => {}} />
		</StoryPreview>
	);
}

export function WithInitialValues() {
	const workout = createWorkout();

	return (
		<StoryPreview>
			<WorkoutForm onSubmit={() => {}} initialValues={workout} />
		</StoryPreview>
	);
}
