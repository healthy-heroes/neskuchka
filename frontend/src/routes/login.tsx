import { createFileRoute } from '@tanstack/react-router';
import { AuthPage } from '@/pages/Auth/Auth.page';

export const Route = createFileRoute('/login')({
	component: AuthPage,
});
