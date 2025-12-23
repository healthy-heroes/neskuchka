import { useMutation } from '@tanstack/react-query';
import { Navigate, useNavigate } from '@tanstack/react-router';
import { Box, Title } from '@mantine/core';
import { useApi } from '@/api/hooks';
import { WorkoutForm } from '../WorkoutForm/WorkoutForm';

export function WorkoutCreate() {
	const navigate = useNavigate();

	const { workouts } = useApi();

	const mutation = useMutation(workouts.createWorkoutMutation());

	if (mutation.isSuccess) {
		return <Navigate to="/workouts/$workoutId" params={{ workoutId: mutation.data.Workout.ID }} />;
	}

	function handleCancel() {
		navigate({ to: '/workouts' });
	}

	const workout = mutation.data;
	return (
		<Box p="md">
			<Title order={2} my="md">
				Создание тренировки
			</Title>

			<WorkoutForm
				data={workout}
				isSubmitting={mutation.isPending}
				onSubmit={mutation.mutate}
				onCancel={handleCancel}
				error={mutation.error}
			/>
		</Box>
	);
}
