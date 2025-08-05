import { Box, Title } from '@mantine/core';
import { getMainTrackWorkouts } from '@/api/methods';
import { WorkoutCard } from '../WorkoutCard/WorkoutCard';
import { WorkoutCardSkeleton } from '../WorkoutCard/WorkoutCardSkeleton';

export function Workouts() {
	const { data, loading, error } = getMainTrackWorkouts();

	return (
		<Box p="md">
			<Title order={2} my="md">
				Тренировки
			</Title>

			{(loading || error) && <WorkoutCardSkeleton cardProps={{ mb: 'xl' }} />}

			{data?.Workouts.map((workout) => {
				return <WorkoutCard cardProps={{ mb: 'xl' }} workout={workout} />;
			})}
		</Box>
	);
}
