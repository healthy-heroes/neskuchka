import {
	InfiniteData,
	QueryKey,
	UseInfiniteQueryOptions,
	UseMutationOptions,
	UseQueryOptions,
} from '@tanstack/react-query';
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

export interface TrackWorkoutData {
	Workout: Workout;
}

/** Страница трека для владельца: строки плюс счётчики шапки списка. */
export interface TrackWorkoutsPageData {
	Workouts: Array<Workout>;

	Total: number;
	Planned: number;
}

export interface UpdateTrackPayload {
	Name: string;
	Description: string;
}

export const WorkoutsKeys = {
	track: () => ['track:main'] as const,
	workouts: () => [...WorkoutsKeys.track(), 'workouts'],
	workout: (id: string) => [...WorkoutsKeys.track(), 'workout', id],

	// Свой ключ: админский список отдаёт ещё и неопубликованное, в кэше
	// публичного списка ему делать нечего
	allWorkouts: () => [...WorkoutsKeys.track(), 'workouts:all'],
};

/** Сколько строк тянем за раз — столько же ждёт бэкенд по умолчанию. */
export const MANAGE_PAGE_SIZE = 8;

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
	getWorkoutQuery(
		id: string
	): UseQueryOptions<ApiResponse<TrackWorkoutData>, Error, TrackWorkoutData> {
		return {
			queryKey: WorkoutsKeys.workout(id),
			queryFn: () => this.api.get<ApiResponse<TrackWorkoutData>>(`tracks/main/workouts/${id}`),
			select: (response) => response.data,
		};
	}

	updateWorkoutMutation(): UseMutationOptions<TrackWorkoutData, Error, Workout> {
		return {
			mutationFn: async (workout: Workout) => {
				const response = await this.api.put<ApiResponse<TrackWorkoutData>, Workout>(
					`tracks/main/workouts/${workout.ID}`,
					workout
				);
				return response.data;
			},
		};
	}

	createWorkoutMutation(): UseMutationOptions<TrackWorkoutData, Error, Workout> {
		return {
			mutationFn: async (workout: Workout) => {
				const response = await this.api.post<ApiResponse<TrackWorkoutData>, Workout>(
					`tracks/main/workouts`,
					workout
				);
				return response.data;
			},
		};
	}

	/**
	 * Весь трек владельцу, страницами — вместе с неопубликованным.
	 *
	 * Счётчики приезжают в каждой странице: они считаются по треку, а не по
	 * загруженному куску, поэтому «Показаны 8 из 41» не врёт с первой страницы.
	 */
	getAllTrackWorkoutsQuery(): UseInfiniteQueryOptions<
		ApiResponse<TrackWorkoutsPageData>,
		Error,
		InfiniteData<ApiResponse<TrackWorkoutsPageData>>,
		QueryKey,
		number
	> {
		return {
			queryKey: WorkoutsKeys.allWorkouts(),
			queryFn: ({ pageParam }) =>
				this.api.get<ApiResponse<TrackWorkoutsPageData>>(
					`tracks/main/workouts?limit=${MANAGE_PAGE_SIZE}&offset=${pageParam}`
				),
			initialPageParam: 0,
			getNextPageParam: (lastPage, allPages) => {
				const loaded = allPages.reduce((sum, page) => sum + page.data.Workouts.length, 0);

				return loaded < lastPage.data.Total ? loaded : undefined;
			},
		};
	}

	updateTrackMutation(): UseMutationOptions<TrackData, Error, UpdateTrackPayload> {
		return {
			mutationFn: async (payload: UpdateTrackPayload) => {
				const response = await this.api.put<ApiResponse<TrackData>, UpdateTrackPayload>(
					`tracks/main`,
					payload
				);
				return response.data;
			},
		};
	}

	deleteWorkoutMutation(): UseMutationOptions<void, Error, string> {
		return {
			mutationFn: (id: string) => this.api.delete<void>(`tracks/main/workouts/${id}`),
		};
	}
}
