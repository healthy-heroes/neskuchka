import { createFileRoute } from '@tanstack/react-router';
import { RequireAuth } from '@/auth/RequireAuth';
import { PageSkeleton } from '@/components/PageSkeleton/PageSkeleton';
import { WorkoutCreate } from '@/components/WorkoutCreate/WorkoutCreate';
import { TrackOwnerOnly } from '@/guards/TrackOwnerOnly';

export const Route = createFileRoute('/workouts/new')({
	component: RouteComponent,
});

function RouteComponent() {
	const loadingComponent = <PageSkeleton hideHeader />;

	return (
		<RequireAuth loadingComponent={loadingComponent}>
			<TrackOwnerOnly loadingComponent={loadingComponent} redirectTo="/workouts">
				<WorkoutCreate />
			</TrackOwnerOnly>
		</RequireAuth>
	);
}
