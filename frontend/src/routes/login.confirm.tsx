import { createFileRoute } from '@tanstack/react-router';
import { LoginConfirm } from '@/components/Login/Login';

export const Route = createFileRoute('/login/confirm')({
	component: () => {
		const { token } = Route.useSearch() as { token: string };
		return <LoginConfirm token={token} />;
	},
});
