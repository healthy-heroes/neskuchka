import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Navigate, useNavigate } from '@tanstack/react-router';
import { Box, Title } from '@mantine/core';
import { useApi } from '@/api/hooks';
import { WorkoutsKeys } from '@/api/services/workouts';
import { WorkoutCardSkeleton } from '../WorkoutCard/WorkoutCardSkeleton';
import { WorkoutForm } from '../WorkoutForm/WorkoutForm';

interface WorkoutEditProps {
	workoutId: string;
}

/**
 * Component for editing a workout
 *
 * @attention This component don't check owner of the workout
 */
export function WorkoutEdit({ workoutId }: WorkoutEditProps) {
	const navigate = useNavigate();
	const queryClient = useQueryClient();
	const { workouts } = useApi();

	// Как и при создании: без сброса правка не доедет до списков трека
	const workoutUpdating = useMutation({
		...workouts.updateWorkoutMutation(),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: WorkoutsKeys.track() });
		},
	});
	const { data, isSuccess, isPending } = useQuery(workouts.getWorkoutQuery(workoutId));

	if (workoutUpdating.isSuccess) {
		return <Navigate to="/workouts/$workoutId" params={{ workoutId }} />;
	}

	if (isPending || !isSuccess) {
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
				trackId={workout.TrackID}
				data={workout}
				isSubmitting={workoutUpdating.isPending}
				onSubmit={workoutUpdating.mutate}
				onCancel={handleCancel}
				error={workoutUpdating.error}
			/>
		</Box>
	);
}
