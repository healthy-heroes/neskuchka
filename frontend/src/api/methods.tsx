import { useFetch, UseFetchOptions, UseFetchReturnValue } from '@mantine/hooks';
import { API_URL } from '../config';
import { Workout } from '../types/domain';

interface RequestOptions extends UseFetchOptions {}

function request<T>(url: string, fetchOptions: RequestOptions): UseFetchReturnValue<T> {
	return useFetch<T>(url, fetchOptions);
}

export interface TrackWorkouts {
	Workouts: Array<Workout>;
}
export function getMainTrackWorkouts(fetchOptions: RequestOptions = {}) {
	return request<TrackWorkouts>(`${API_URL}/tracks/main/last_workouts`, fetchOptions);
}
