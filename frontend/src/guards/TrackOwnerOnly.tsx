import { useQuery } from '@tanstack/react-query';
import { Navigate } from '@tanstack/react-router';
import { useApi } from '@/api/hooks';

type TrackOwnerOnlyProps = {
	children: React.ReactNode;
	loadingComponent: React.ReactNode;

	redirectTo: string;
};

export function TrackOwnerOnly({ children, loadingComponent, redirectTo }: TrackOwnerOnlyProps) {
	const api = useApi();
	const { data, isPending, isSuccess } = useQuery(api.workouts.getMainTrackQuery());

	if (isPending) {
		return loadingComponent;
	}

	if (!isSuccess || !data.IsOwner) {
		return <Navigate to={redirectTo} />;
	}

	return children;
}
