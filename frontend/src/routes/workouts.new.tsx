import { createFileRoute } from '@tanstack/react-router';
import { WorkoutCreate } from '@/components/WorkoutCreate/WorkoutCreate';

export const Route = createFileRoute('/workouts/new')({
	component: RouteComponent,
});

function RouteComponent() {
	return <WorkoutCreate />;
}
