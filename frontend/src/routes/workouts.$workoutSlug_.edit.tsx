import { createFileRoute } from '@tanstack/react-router';
import { RequireAuth } from '@/auth/RequireAuth';
import { PageSkeleton } from '@/components/PageSkeleton/PageSkeleton';
import { WorkoutEdit } from '@/components/WorkoutEdit/WorkoutEdit';
import { TrackOwnerOnly } from '@/guards/TrackOwnerOnly';

export const Route = createFileRoute('/workouts/$workoutSlug_/edit')({
	component: RouteComponent,
});

function RouteComponent() {
	const { workoutSlug } = Route.useParams();

	const loadingComponent = <PageSkeleton hideHeader />;

	return (
		<RequireAuth loadingComponent={loadingComponent}>
			<TrackOwnerOnly loadingComponent={loadingComponent} redirectTo={`/workouts/${workoutSlug}`}>
				<WorkoutEdit workoutSlug={workoutSlug} />
			</TrackOwnerOnly>
		</RequireAuth>
	);
}
