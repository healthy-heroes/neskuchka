import { useQuery } from '@tanstack/react-query';
import { useApi } from '@/api/hooks';
import { TrackHeader } from '@/pages/MainTrack/TrackHeader/TrackHeader';
import { isPublished, isToday } from '@/utils/dates';
import { FeaturedWorkout } from '../FeaturedWorkout/FeaturedWorkout';
import { WorkoutCardSkeleton } from '../WorkoutCard/WorkoutCardSkeleton';
import { WorkoutHistory } from '../WorkoutHistory/WorkoutHistory';
import classes from './Workouts.module.css';

/**
 * Workouts — страница трека.
 *
 * Собрана вокруг одной задачи: человек заходит, сразу видит сегодняшнюю
 * тренировку и начинает её. Всё остальное — история, прогресс — ниже и мельче.
 */
export function Workouts() {
	const { workouts } = useApi();

	//todo: handle errors
	const { data, isPending } = useQuery(workouts.getMainTrackWorkoutsQuery());
	const { data: track } = useQuery(workouts.getMainTrackQuery());

	const published = (data?.Workouts ?? []).filter((workout) => isPublished(workout.Date));

	// В фокусе сегодняшняя, а если её нет — последняя опубликованная:
	// пустой главный экран хуже, чем вчерашняя тренировка на нём
	const featured = published.find((workout) => isToday(workout.Date)) ?? published[0];

	return (
		<div className={classes.page}>
			{track && <TrackHeader track={track} workouts={published} />}

			<div className={classes.body}>
				{isPending ? (
					<WorkoutCardSkeleton />
				) : (
					<>
						{featured && <FeaturedWorkout workout={featured} />}

						<WorkoutHistory workouts={published} featured={featured} />
					</>
				)}
			</div>
		</div>
	);
}
