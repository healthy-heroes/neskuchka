import { createFileRoute } from '@tanstack/react-router';
import { WorkoutEdit } from '@/components/WorkoutEdit/WorkoutEdit';

export const Route = createFileRoute('/workouts/$workoutId_/edit')({
	component: RouteComponent,
});

function RouteComponent() {
	const { workoutId } = Route.useParams();
	return <WorkoutEdit workoutId={workoutId} />;
}
