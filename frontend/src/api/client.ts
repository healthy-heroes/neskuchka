import { QueryClient, QueryClientConfig } from '@tanstack/react-query';

export function createApiClient(options: QueryClientConfig = {}) {
	return new QueryClient(options);
}
