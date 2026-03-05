import { createFileRoute } from '@tanstack/react-router';
import { RequireAuth } from '@/auth/RequireAuth';
import { PageSkeleton } from '@/components/PageSkeleton/PageSkeleton';
import { SettingsPage } from '@/pages/Settings/Settings.page';

export const Route = createFileRoute('/settings')({
	component: () => (
		<RequireAuth loadingComponent={<PageSkeleton />}>
			<SettingsPage />
		</RequireAuth>
	),
});
