import { useQuery } from '@tanstack/react-query';
import { Box, Title } from '@mantine/core';
import { useApiService } from '@/api/provider';
import { WorkoutCard } from '../WorkoutCard/WorkoutCard';
import { WorkoutCardSkeleton } from '../WorkoutCard/WorkoutCardSkeleton';

export function Workouts() {
	const {
		queries: { workouts },
	} = useApiService();

	//todo: handle errors
	const { data, isSuccess, isLoading } = useQuery(workouts.getMainTrackWorkoutsQuery());

	return (
		<Box p="md">
			<Title order={2} my="md">
				Тренировки
			</Title>

			{(isLoading || !isSuccess) && <WorkoutCardSkeleton cardProps={{ mb: 'xl' }} />}

			{isSuccess &&
				data.Workouts.map((workout) => {
					return <WorkoutCard key={workout.ID} cardProps={{ mb: 'xl' }} workout={workout} />;
				})}
		</Box>
	);
}
