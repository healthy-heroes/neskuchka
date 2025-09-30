import { Outlet } from '@tanstack/react-router';
import { Header } from '@/components/Header/Header';
import { TrackHeader } from './TrackHeader/TrackHeader';

export function MainTrackPage() {
	return (
		<>
			<Header />
			<TrackHeader />
			<Outlet />
		</>
	);
}
