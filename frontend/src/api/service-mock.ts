import { UseFetchReturnValue } from '@mantine/hooks';
import ApiService from './service';

export function createMockApiService(overrides: Partial<ApiService> = {}): ApiService {
	return {
		...overrides,
	} as unknown as ApiService;
}

export function createMock<T>(data: T, loading: boolean = false, error: Error | null = null) {
	return function (): UseFetchReturnValue<T> {
		return {
			data,
			loading,
			error,
			refetch: () => Promise.reject(new Error('Not implemented')),
			abort: () => {
				throw new Error('Not implemented');
			},
		};
	};
}
