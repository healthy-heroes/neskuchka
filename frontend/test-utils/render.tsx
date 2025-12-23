import React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import {
	RenderOptions,
	renderHook as rtlRenderHook,
	render as testingLibraryRender,
} from '@testing-library/react';
import { MantineProvider } from '@mantine/core';
import { ApiContext } from '../src/api/provider';
import ApiService from '../src/api/service';
import { theme } from '../src/theme';

/**
 * Creates a QueryClient configured for tests
 */
export function createTestQueryClient() {
	return new QueryClient({
		defaultOptions: {
			queries: {
				retry: false,
				gcTime: 0,
				staleTime: Infinity, // Don't refetch in tests
			},
		},
	});
}

type TestWrapperOptions = {
	apiService?: ApiService | null;
	queryClient?: QueryClient;
};

function createWrapper(options: TestWrapperOptions = {}) {
	const queryClient = options.queryClient ?? createTestQueryClient();

	return function TestWrapper({ children }: { children: React.ReactNode }) {
		return (
			<QueryClientProvider client={queryClient}>
				<ApiContext.Provider value={options.apiService ?? null}>
					<MantineProvider theme={theme} env="test">
						{children}
					</MantineProvider>
				</ApiContext.Provider>
			</QueryClientProvider>
		);
	};
}

export function render(
	ui: React.ReactNode,
	options: TestWrapperOptions & Omit<RenderOptions, 'wrapper'> = {}
) {
	const { apiService, queryClient, ...renderOptions } = options;

	return testingLibraryRender(ui, {
		wrapper: createWrapper({ apiService, queryClient }),
		...renderOptions,
	});
}

export function renderHook<TResult, TProps>(
	hook: (props: TProps) => TResult,
	options: TestWrapperOptions = {}
) {
	return rtlRenderHook(hook, {
		wrapper: createWrapper(options),
	});
}
