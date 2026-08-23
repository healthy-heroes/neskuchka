import createWorkout from '@/api/fixtures/workout';
import { StoryPreview } from '../StoryBook/StoryPreview';
import { WorkoutSections } from './WorkoutSections';
import { WorkoutSectionsSkeleton } from './WorkoutSectionsSkeleton';

export default {
	title: 'WorkoutSections',
};

const { Sections } = createWorkout({ ID: 'workout-sections' });

export function Default() {
	return (
		<StoryPreview>
			<WorkoutSections sections={Sections} />
		</StoryPreview>
	);
}

export function Skeleton() {
	return (
		<StoryPreview>
			<WorkoutSectionsSkeleton />
		</StoryPreview>
	);
}
