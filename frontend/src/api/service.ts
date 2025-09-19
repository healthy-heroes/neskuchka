import { useFetch, UseFetchOptions, UseFetchReturnValue } from '@mantine/hooks';
import { TrackWorkout, TrackWorkouts } from './types';

type ApiConfig = {
	apiUrl: string;
};

export interface RequestOptions extends UseFetchOptions {}

class ApiService {
	private readonly config: ApiConfig;

	constructor(config: ApiConfig) {
		this.config = config;
	}

	private request<T>(url: string, fetchOptions: RequestOptions): UseFetchReturnValue<T> {
		return useFetch<T>(url, fetchOptions);
	}

	/**
	 * Get the last workouts for the main track
	 */
	getMainTrackWorkouts(fetchOptions: RequestOptions = {}) {
		return this.request<TrackWorkouts>(
			`${this.config.apiUrl}/tracks/main/last_workouts`,
			fetchOptions
		);
	}

	/**
	 * Get concrete workout by id
	 */
	getWorkout(id: string, fetchOptions: RequestOptions = {}) {
		return this.request<TrackWorkout>(`${this.config.apiUrl}/workouts/${id}`, fetchOptions);
	}
}

export default ApiService;
