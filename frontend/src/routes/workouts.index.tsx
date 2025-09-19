import { createFileRoute } from '@tanstack/react-router';
import { Workouts } from '@/components/Workouts/Workouts';

export const Route = createFileRoute('/workouts/')({
	component: Workouts,
});
