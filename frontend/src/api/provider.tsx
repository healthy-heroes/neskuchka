import React, { createContext, useContext } from 'react';
import { QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { Api } from './api';
import { createApiClient } from './client';
import { ApiQueries } from './queries';

type ApiService = {
	queries: ApiQueries;
};

export const ApiContext = createContext<ApiService | null>(null);

export function useApiService(): ApiService {
	const context = useContext(ApiContext);
	if (context === null) {
		throw new Error('useApiService must be used within a ApiProvider');
	}

	return context;
}

type ApiProviderProps = {
	children: React.ReactNode;
};

export function ApiProvider({ children }: ApiProviderProps) {
	const api = new Api();
	const queries = new ApiQueries(api);

	//todo: add defaults
	const queryClient = createApiClient({
		defaultOptions: {},
	});

	return (
		<ApiContext.Provider value={{ queries }}>
			<QueryClientProvider client={queryClient}>
				{children}
				<ReactQueryDevtools />
			</QueryClientProvider>
		</ApiContext.Provider>
	);
}
