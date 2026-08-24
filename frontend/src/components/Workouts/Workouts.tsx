import { useQuery } from '@tanstack/react-query';
import { Box, Container, Text } from '@mantine/core';
import { useApi } from '@/api/hooks';
import { TrackHeader } from '@/pages/MainTrack/TrackHeader/TrackHeader';
import { isToday } from '@/utils/dates';
import { FeaturedWorkout } from '../FeaturedWorkout/FeaturedWorkout';
import { FeaturedWorkoutSkeleton } from '../FeaturedWorkout/FeaturedWorkoutSkeleton';
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

	// Неопубликованных в ответе нет: их отсекает бэкенд, а не эта страница
	const published = data?.Workouts ?? [];

	// В фокусе сегодняшняя, а если её нет — последняя опубликованная:
	// пустой главный экран хуже, чем вчерашняя тренировка на нём
	const featured = published.find((workout) => isToday(workout.Date)) ?? published[0];

	return (
		<Container px={0}>
			{track && <TrackHeader track={track} workouts={published} loading={isPending} />}

			<Box px="xl" py="xl" bg="gray.0">
				{isPending && <FeaturedWorkoutSkeleton />}

				{!isPending && !featured && (
					<Text component="p" my={0} p="xl" fz="md" lh="md" c="gray.8" className={classes.empty}>
						Тренировок пока нет — они появятся, когда их опубликуют
					</Text>
				)}

				{!isPending && featured && (
					<>
						<FeaturedWorkout workout={featured} />

						<WorkoutHistory workouts={published} featured={featured} />
					</>
				)}
			</Box>
		</Container>
	);
}
