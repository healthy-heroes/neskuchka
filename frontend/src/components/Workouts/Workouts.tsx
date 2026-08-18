import { useQuery } from '@tanstack/react-query';
import { useApi } from '@/api/hooks';
import { TrackHeader } from '@/pages/MainTrack/TrackHeader/TrackHeader';
import { isToday } from '@/utils/dates';
import { NoWorkoutToday, TodayWorkout } from '../TodayWorkout/TodayWorkout';
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

	const list = data?.Workouts ?? [];
	const todayWorkout = list.find((workout) => isToday(workout.Date));

	// Последняя опубликованная — на случай, когда на сегодня тренировки нет
	const lastPublished = list.find((workout) => new Date(workout.Date) <= new Date());

	return (
		<div className={classes.page}>
			{track && <TrackHeader track={track} workouts={list} />}

			<div className={classes.body}>
				{isPending ? (
					<WorkoutCardSkeleton />
				) : (
					<>
						{todayWorkout ? (
							<TodayWorkout workout={todayWorkout} />
						) : (
							<NoWorkoutToday last={lastPublished} />
						)}

						<WorkoutHistory workouts={list} />
					</>
				)}
			</div>
		</div>
	);
}
