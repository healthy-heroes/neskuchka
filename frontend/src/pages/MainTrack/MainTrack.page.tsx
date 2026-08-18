import { useQuery } from '@tanstack/react-query';
import { Navigate, Outlet } from '@tanstack/react-router';
import { useApi } from '@/api/hooks';
import { Header } from '@/components/Header/Header';
import { PageSkeleton } from '@/components/PageSkeleton/PageSkeleton';
import classes from './MainTrack.page.module.css';

/**
 * MainTrackPage — общий слой для всего /workouts.
 *
 * Грузит трек, дальше дети берут его из кэша мгновенно. Шапка трека живёт
 * не здесь, а на самой странице трека: на странице тренировки её быть не должно.
 */
export function MainTrackPage() {
	const api = useApi();
	const { isPending, isSuccess } = useQuery(api.workouts.getMainTrackQuery());

	if (isPending) {
		return <PageSkeleton />;
	}

	if (!isSuccess) {
		return <Navigate to="/" />;
	}

	return (
		<div className={classes.track}>
			<Header />
			<Outlet />
		</div>
	);
}
