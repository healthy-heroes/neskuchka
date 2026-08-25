import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Navigate, useNavigate } from '@tanstack/react-router';
import { Box, Title } from '@mantine/core';
import { useApi } from '@/api/hooks';
import { WorkoutsKeys } from '@/api/services/workouts';
import { WorkoutCardSkeleton } from '../WorkoutCard/WorkoutCardSkeleton';
import { WorkoutForm } from '../WorkoutForm/WorkoutForm';

export function WorkoutCreate() {
	const navigate = useNavigate();

	const queryClient = useQueryClient();
	const { workouts } = useApi();

	const trackQuery = useQuery(workouts.getMainTrackQuery());

	// Списки трека кэшируются, и без сброса свежесозданная тренировка не
	// появится ни в истории, ни на экране управления
	const mutation = useMutation({
		...workouts.createWorkoutMutation(),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: WorkoutsKeys.track() });
		},
	});

	if (mutation.isSuccess) {
		return <Navigate to="/workouts/$workoutId" params={{ workoutId: mutation.data.Workout.ID }} />;
	}

	function handleCancel() {
		navigate({ to: '/workouts' });
	}

	if (trackQuery.isPending || !trackQuery.data) {
		return (
			<Box p="md">
				<Title order={2} my="md">
					Создание тренировки
				</Title>

				<WorkoutCardSkeleton cardProps={{ mb: 'xl' }} />
			</Box>
		);
	}

	const workout = mutation.data;
	return (
		<Box p="md">
			<Title order={2} my="md">
				Создание тренировки
			</Title>

			<WorkoutForm
				// we loads track
				trackId={trackQuery.data.Track.ID}
				data={workout}
				isSubmitting={mutation.isPending}
				onSubmit={mutation.mutate}
				onCancel={handleCancel}
				error={mutation.error}
			/>
		</Box>
	);
}
