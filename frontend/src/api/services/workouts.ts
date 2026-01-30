import { UseMutationOptions, UseQueryOptions } from '@tanstack/react-query';
import { Track, Workout } from '@/types/domain';
import Service from './service';

//todo: move out of this file
interface ApiResponse<T> {
	data: T;
}

export interface TrackData {
	Track: Track;
	IsOwner: boolean;
}

export interface TrackWorkoutsData {
	Workouts: Array<Workout>;
}

export interface TrackWorkout {
	Workout: Workout;
}

export const WorkoutsKeys = {
	track: () => ['track:main'] as const,
	workouts: () => [...WorkoutsKeys.track(), 'workouts'],
	workout: (id: string) => [...WorkoutsKeys.track(), 'workout', id],
};

export class WorkoutsService extends Service {
	/**
	 * Get the main track
	 */
	getMainTrackQuery(): UseQueryOptions<ApiResponse<TrackData>, Error, TrackData> {
		return {
			queryKey: WorkoutsKeys.track(),
			queryFn: () => this.api.get<ApiResponse<TrackData>>(`tracks/main`),
			select: (response) => response.data,
		};
	}

	/**
	 * Get the last workouts for the main track
	 */
	getMainTrackWorkoutsQuery(): UseQueryOptions<
		ApiResponse<TrackWorkoutsData>,
		Error,
		TrackWorkoutsData
	> {
		return {
			queryKey: WorkoutsKeys.workouts(),
			queryFn: () => this.api.get<ApiResponse<TrackWorkoutsData>>(`tracks/main/last_workouts`),
			select: (response) => response.data,
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
