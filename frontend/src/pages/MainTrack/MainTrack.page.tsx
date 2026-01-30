import { useQuery } from '@tanstack/react-query';
import { Navigate, Outlet } from '@tanstack/react-router';
import { useApi } from '@/api/hooks';
import { Header } from '@/components/Header/Header';
import { PageSkeleton } from '@/components/PageSkeleton/PageSkeleton';
import { TrackHeader } from './TrackHeader/TrackHeader';

export function MainTrackPage() {
	const api = useApi();
	const { data, isPending, isSuccess } = useQuery(api.workouts.getMainTrackQuery());

	if (isPending) {
		return <PageSkeleton />;
	}

	if (!isSuccess) {
		return <Navigate to="/" />;
	}

	return (
		<>
			<Header />
			<TrackHeader track={data} />
			<Outlet />
		</>
	);
}
