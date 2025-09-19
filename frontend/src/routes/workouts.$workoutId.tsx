import { createFileRoute } from '@tanstack/react-router';
import { WorkoutView } from '@/components/WorkoutView/WorkoutView';

export const Route = createFileRoute('/workouts/$workoutId')({
	component: () => {
		const { workoutId } = Route.useParams();
		return <WorkoutView workoutId={workoutId} />;
	},
});
