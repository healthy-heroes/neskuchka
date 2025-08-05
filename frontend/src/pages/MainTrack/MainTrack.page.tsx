import { Header } from '@/components/Header/Header';
import { Workouts } from '@/components/Workouts/Workouts';
import { TrackHeader } from './TrackHeader/TrackHeader';

export function MainTrackPage() {
	return (
		<>
			<Header />
			<TrackHeader />
			<Workouts />
		</>
	);
}
