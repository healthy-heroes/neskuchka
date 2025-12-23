import React, { createContext } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import ApiClient from './client';
import ApiService from './service';

export const ApiContext = createContext<ApiService | null>(null);

const service = new ApiService(new ApiClient());

export function ApiProvider({ children }: { children: React.ReactNode }) {
	//todo: add defaults
	const queryClient = new QueryClient({
		defaultOptions: {},
	});

	return (
		<ApiContext.Provider value={service}>
			<QueryClientProvider client={queryClient}>
				{children}
				<ReactQueryDevtools />
			</QueryClientProvider>
		</ApiContext.Provider>
	);
}
