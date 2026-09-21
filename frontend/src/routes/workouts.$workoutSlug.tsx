import { createFileRoute } from '@tanstack/react-router';
import { WorkoutView } from '@/components/WorkoutView/WorkoutView';

export const Route = createFileRoute('/workouts/$workoutSlug')({
	component: () => {
		const { workoutSlug } = Route.useParams();
		return <WorkoutView workoutSlug={workoutSlug} />;
	},
});
