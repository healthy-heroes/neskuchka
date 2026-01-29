import { UseMutationOptions, UseQueryOptions } from '@tanstack/react-query';
import { Workout } from '@/types/domain';
import Service from './service';

export interface TrackWorkouts {
	Workouts: Array<Workout>;
}

export interface TrackWorkout {
	Workout: Workout;
}

export const WorkoutsKeys = {
	all: ['workouts'] as const,
	byTrack: () => [...WorkoutsKeys.all, 'track:main'],
	workout: (id: string) => [...WorkoutsKeys.all, 'workout', id],
};

export class WorkoutsService extends Service {
	/**
	 * Get the last workouts for the main track
	 */
	getMainTrackWorkoutsQuery(): UseQueryOptions<TrackWorkouts> {
		return {
			queryKey: WorkoutsKeys.byTrack(),
			queryFn: () => this.api.get<TrackWorkouts>(`tracks/main/last_workouts`),
		};
	}

	/**
	 * Get concrete workout by id
	 */
	getWorkoutQuery(id: string): UseQueryOptions<TrackWorkout> {
		return {
			queryKey: WorkoutsKeys.workout(id),
			queryFn: () => this.api.get<TrackWorkout>(`tracks/main/workouts/${id}`),
		};
	}

	updateWorkoutMutation(): UseMutationOptions<TrackWorkout, Error, Workout> {
		return {
			mutationFn: (workout: Workout) => {
				return this.api.put<TrackWorkout, Workout>(`tracks/main/workouts/${workout.ID}`, workout);
			},
		};
	}

	createWorkoutMutation(): UseMutationOptions<TrackWorkout, Error, Workout> {
		return {
			mutationFn: (workout: Workout) => {
				return this.api.post<TrackWorkout, Workout>(`tracks/main/workouts`, workout);
			},
		};
	}
}
