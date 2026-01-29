import { createFileRoute } from '@tanstack/react-router';
import { HomePage } from '@/pages/Home.page';

export const Route = createFileRoute('/')({
	component: HomePage,
});
