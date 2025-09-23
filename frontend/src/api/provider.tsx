import React, { createContext, useContext } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import ApiClient from './client';
import ApiService from './service';

type ApiContextProps = {
	service: ApiService;
};

export const ApiContext = createContext<ApiContextProps | null>(null);

export function useApiService(): ApiContextProps {
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
	const api = new ApiClient();
	const service = new ApiService(api);

	//todo: add defaults
	const queryClient = new QueryClient({
		defaultOptions: {},
	});

	return (
		<ApiContext.Provider value={{ service }}>
			<QueryClientProvider client={queryClient}>
				{children}
				<ReactQueryDevtools />
			</QueryClientProvider>
		</ApiContext.Provider>
	);
}
