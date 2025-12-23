import { useMutation, useQuery } from '@tanstack/react-query';
import { Navigate, useNavigate } from '@tanstack/react-router';
import { Box, Title } from '@mantine/core';
import { useApi } from '@/api/hooks';
import { WorkoutCardSkeleton } from '../WorkoutCard/WorkoutCardSkeleton';
import { WorkoutForm } from '../WorkoutForm/WorkoutForm';

interface WorkoutEditProps {
	workoutId: string;
}

export function WorkoutEdit({ workoutId }: WorkoutEditProps) {
	const navigate = useNavigate();
	const { workouts } = useApi();

	const mutation = useMutation(workouts.updateWorkoutMutation());
	const { data, isSuccess, isLoading } = useQuery(workouts.getWorkoutQuery(workoutId));

	if (mutation.isSuccess) {
		return <Navigate to="/workouts/$workoutId" params={{ workoutId }} />;
	}

	if (isLoading || !isSuccess) {
		return (
			<Box p="md">
				<Title order={2} my="md">
					Редактирование тренировки
				</Title>

				<WorkoutCardSkeleton cardProps={{ mb: 'xl' }} />
			</Box>
		);
	}

	function handleCancel() {
		navigate({ to: '/workouts/$workoutId', params: { workoutId } });
	}

	const workout = data.Workout;
	return (
		<Box p="md">
			<Title order={2} my="md">
				Редактирование тренировки
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
