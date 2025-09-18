import { createFileRoute } from '@tanstack/react-router';
import { MainTrackPage } from '@/pages/MainTrack/MainTrack.page';

export const Route = createFileRoute('/workouts')({
	component: MainTrackPage,
});
