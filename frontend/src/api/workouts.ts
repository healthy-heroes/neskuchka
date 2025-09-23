import { UseQueryOptions } from '@tanstack/react-query';
import { Workout } from '@/types/domain';
import Api from './api';

export interface TrackWorkouts {
	Workouts: Array<Workout>;
}

export interface TrackWorkout {
	Workout: Workout;
}

export class WorkoutsQueries {
	static readonly keys = {
		all: ['workouts'] as const,
		byTrack: () => [...WorkoutsQueries.keys.all, 'track:main'],
		workout: (id: string) => [...WorkoutsQueries.keys.all, 'workout', id],
	};

	constructor(private readonly api: Api) {
		this.api = api;
	}

	/**
	 * Get the last workouts for the main track
	 */
	getMainTrackWorkoutsQuery(): UseQueryOptions<TrackWorkouts> {
		return {
			queryKey: WorkoutsQueries.keys.byTrack(),
			queryFn: () => this.api.get<TrackWorkouts>(`tracks/main/last_workouts`),
		};
	}

	/**
	 * Get concrete workout by id
	 */
	getWorkoutQuery(id: string): UseQueryOptions<TrackWorkout> {
		return {
			queryKey: WorkoutsQueries.keys.workout(id),
			queryFn: () => this.api.get<TrackWorkout>(`workouts/${id}`),
		};
	}
}
