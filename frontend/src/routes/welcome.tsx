import { createFileRoute } from '@tanstack/react-router';
import { LandingPage } from '@/pages/Landing/Landing.page';

export const Route = createFileRoute('/welcome')({
	component: LandingPage,
});
