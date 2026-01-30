import { createFileRoute } from '@tanstack/react-router';
import { PageSkeleton } from '@/components/PageSkeleton/PageSkeleton';
import { WorkoutEdit } from '@/components/WorkoutEdit/WorkoutEdit';
import { TrackOwnerOnly } from '@/guards/TrackOwnerOnly';

export const Route = createFileRoute('/workouts/$workoutId_/edit')({
	component: RouteComponent,
});

function RouteComponent() {
	const { workoutId } = Route.useParams();

	const loadingComponent = <PageSkeleton hideHeader />;

	return (
		<TrackOwnerOnly loadingComponent={loadingComponent} redirectTo={`/workouts/${workoutId}`}>
			<WorkoutEdit workoutId={workoutId} />
		</TrackOwnerOnly>
	);
}
