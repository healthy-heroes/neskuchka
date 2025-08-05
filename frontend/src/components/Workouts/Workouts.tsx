import { Box, Title } from '@mantine/core';
import { WorkoutCard } from '../WorkoutCard/WorkoutCard';

export function Workouts() {
	return (
		<Box p="md">
			<Title order={2} my="md">
				Тренировки
			</Title>

			<WorkoutCard cardProps={{ mb: 'xl' }} />
			<WorkoutCard cardProps={{ mb: 'xl' }} />
		</Box>
	);
}
