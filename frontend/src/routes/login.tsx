import { createFileRoute } from '@tanstack/react-router';
import { RequireAuth } from '@/auth/RequireAuth';
import { PageSkeleton } from '@/components/PageSkeleton/PageSkeleton';
import { AuthPage } from '@/pages/Auth/Auth.page';

export const Route = createFileRoute('/login')({
	component: RouteComponent,
});

function RouteComponent() {
	return (
		<RequireAuth loadingComponent={<PageSkeleton />} guestOnly>
			<AuthPage />
		</RequireAuth>
	);
}
