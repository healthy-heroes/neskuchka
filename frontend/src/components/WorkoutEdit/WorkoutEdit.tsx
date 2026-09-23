import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Navigate, useNavigate } from '@tanstack/react-router';
import { Box, Title } from '@mantine/core';
import { useApi } from '@/api/hooks';
import { WorkoutsKeys } from '@/api/services/workouts';
import { WorkoutCardSkeleton } from '../WorkoutCard/WorkoutCardSkeleton';
import { WorkoutForm } from '../WorkoutForm/WorkoutForm';

interface WorkoutEditProps {
	workoutSlug: string;
}

/**
 * Component for editing a workout
 *
 * @attention This component don't check owner of the workout
 */
export function WorkoutEdit({ workoutSlug }: WorkoutEditProps) {
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
	const { data, isSuccess, isPending } = useQuery(workouts.getWorkoutQuery(workoutSlug));

	// Слаг берём из ответа, а не из адреса: он выводится из даты, и после
	// переноса тренировки старый уже никуда не ведёт
	if (workoutUpdating.isSuccess) {
		return (
			<Navigate
				to="/workouts/$workoutSlug"
				params={{ workoutSlug: workoutUpdating.data.Workout.Slug }}
			/>
		);
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
		navigate({ to: '/workouts/$workoutSlug', params: { workoutSlug } });
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
