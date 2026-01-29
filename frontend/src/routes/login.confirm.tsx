import { createFileRoute } from '@tanstack/react-router';
import { LoginConfirm } from '@/components/LoginConfirm/LoginConfirm';

export const Route = createFileRoute('/login/confirm')({
	component: () => {
		const { token } = Route.useSearch() as { token?: string };
		return <LoginConfirm token={token ? String(token) : undefined} />;
	},
});
