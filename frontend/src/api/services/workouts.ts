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

	/** Пусто на последней странице. */
	NextCursor: string;

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

/**
 * Разворачивает конверт `{ data }` у каждой страницы и убирает строки,
 * попавшие в две страницы сразу.
 *
 * Курсор защищает от вставок и удалений, но не от правки самой даты: если
 * тренировку перенесли ниже курсора, пока читали, она приедет второй раз.
 * Окно правки держит это в узде — всё, что ниже вчерашнего дня, изменить
 * уже нельзя, — так что дубль возможен, только пока курсор сам внутри окна.
 * Но React-ключи от этого ломаются по-настоящему, а стоит проверка недорого.
 *
 * Из двух копий остаётся нижняя: правка даты сдвинула строку именно туда, и
 * там же её свежая версия.
 *
 * Обратный случай — строка уехала выше курсора и не попала ни в одну
 * страницу — отсюда не виден и лечится только перезапросом.
 *
 * Вынесено из опций намеренно: react-query прогоняет select на каждый рендер и
 * сверяет результат глубоким сравнением, а новая функция каждый раз лишает его
 * мемоизации.
 */
export function unwrapPages(
	data: InfiniteData<ApiResponse<TrackWorkoutsPageData>>
): InfiniteData<TrackWorkoutsPageData> {
	const seen = new Set<string>();
	const pages: Array<TrackWorkoutsPageData> = [];

	// С конца, чтобы «первой встреченной» оказалась нижняя копия
	for (let i = data.pages.length - 1; i >= 0; i -= 1) {
		const page = data.pages[i].data;
		const kept = page.Workouts.filter((workout) => {
			if (seen.has(workout.ID)) {
				return false;
			}

			seen.add(workout.ID);
			return true;
		});

		pages.unshift(kept.length === page.Workouts.length ? page : { ...page, Workouts: kept });
	}

	return { pages, pageParams: data.pageParams };
}

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
	 * Пагинация курсорная: трек правят, пока его читают, и offset после каждой
	 * вставки или удаления сдвигает всё, что ниже, — одна строка приходит дважды,
	 * другая теряется. Курсор непрозрачный, конец списка сервер отмечает пустым.
	 *
	 * Счётчики приезжают в каждой странице: они считаются по треку, а не по
	 * загруженному куску, поэтому «Показаны 8 из 41» не врёт с первой страницы.
	 */
	getAllTrackWorkoutsQuery(): UseInfiniteQueryOptions<
		ApiResponse<TrackWorkoutsPageData>,
		Error,
		InfiniteData<TrackWorkoutsPageData>,
		QueryKey,
		string
	> {
		return {
			queryKey: WorkoutsKeys.allWorkouts(),
			queryFn: ({ pageParam }) =>
				this.api.get<ApiResponse<TrackWorkoutsPageData>>(
					`tracks/main/workouts?limit=${MANAGE_PAGE_SIZE}&after=${encodeURIComponent(pageParam)}`
				),
			initialPageParam: '',
			getNextPageParam: (lastPage) => lastPage.data.NextCursor || undefined,
			select: unwrapPages,
			// Иначе каждый возврат во вкладку перезапрашивает все загруженные
			// страницы подряд
			staleTime: 30 * 1000,
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
