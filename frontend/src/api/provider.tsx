import React, { createContext, useContext } from 'react';
import ApiService from './service';

export const ApiContext = createContext<ApiService | null>(null);

export function useApi(): ApiService {
	const context = useContext(ApiContext);
	if (context === null) {
		throw new Error('useApi must be used within a ApiProvider');
	}

	return context;
}

type ApiProviderProps = {
	apiService: ApiService;

	children: React.ReactNode;
};

export function ApiProvider({ apiService, children }: ApiProviderProps) {
	return <ApiContext.Provider value={apiService}>{children}</ApiContext.Provider>;
}
