import { createFileRoute } from '@tanstack/react-router';
import { RequireAuth } from '@/auth/RequireAuth';
import { PageSkeleton } from '@/components/PageSkeleton/PageSkeleton';
import { TrackOwnerOnly } from '@/guards/TrackOwnerOnly';
import { TrackManagePage } from '@/pages/TrackManage/TrackManage.page';

export const Route = createFileRoute('/workouts/manage')({
	component: RouteComponent,
});

function RouteComponent() {
	const loadingComponent = <PageSkeleton hideHeader />;

	return (
		<RequireAuth loadingComponent={loadingComponent}>
			<TrackOwnerOnly loadingComponent={loadingComponent} redirectTo="/workouts">
				<TrackManagePage />
			</TrackOwnerOnly>
		</RequireAuth>
	);
}
