import { Box, Title } from '@mantine/core';
import { useApi } from '@/api/provider';
import { WorkoutCard } from '../WorkoutCard/WorkoutCard';
import { WorkoutCardSkeleton } from '../WorkoutCard/WorkoutCardSkeleton';

export function Workouts() {
	const apiService = useApi();
	const { data, loading, error } = apiService.getMainTrackWorkouts();

	return (
		<Box p="md">
			<Title order={2} my="md">
				Тренировки
			</Title>

			{(loading || error) && <WorkoutCardSkeleton cardProps={{ mb: 'xl' }} />}

			{data?.Workouts.map((workout) => {
				return <WorkoutCard key={workout.ID} cardProps={{ mb: 'xl' }} workout={workout} />;
			})}
		</Box>
	);
}
