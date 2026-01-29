import React, { createContext } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import ApiClient from './client';
import ApiService from './service';

export const ApiContext = createContext<ApiService | null>(null);

const apiService = new ApiService(new ApiClient());

//todo: add defaults
const queryClient = new QueryClient({
	defaultOptions: {
		queries: {
			retry: 3,
			retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000),
		},
	},
});

export function ApiProvider({ children }: { children: React.ReactNode }) {
	return (
		<ApiContext.Provider value={apiService}>
			<QueryClientProvider client={queryClient}>
				{children}
				<ReactQueryDevtools />
			</QueryClientProvider>
		</ApiContext.Provider>
	);
}
