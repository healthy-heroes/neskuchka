import { QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import {
	createMemoryHistory,
	createRootRoute,
	createRouter,
	RouteComponent,
	RouterProvider,
} from '@tanstack/react-router';
import { Paper, PaperProps } from '@mantine/core';
import { createApiClient } from '@/api/client';
import { ApiMock } from '@/api/fixtures/api';
import { ApiContext } from '@/api/provider';
import { ApiQueries } from '@/api/queries';

const createStoryRouter = (component: RouteComponent) => {
	return createRouter({
		history: createMemoryHistory(),
		routeTree: createRootRoute({
			component,
		}),
	});
};

export interface StoryPreviewProps {
	children: React.ReactNode;

	paperOptions?: PaperProps;
	isPage?: boolean;

	queries?: ApiQueries;
}

export function StoryPreview(props: StoryPreviewProps) {
	const queryClient = createApiClient();

	const queries = props.queries ?? new ApiQueries(new ApiMock());

	return (
		<ApiContext.Provider value={{ queries }}>
			<QueryClientProvider client={queryClient}>
				{getPageWrapper(props)}
				<ReactQueryDevtools />
			</QueryClientProvider>
		</ApiContext.Provider>
	);
}

function getPageWrapper({ children, paperOptions, isPage }: StoryPreviewProps) {
	const defaultOptions: PaperProps = {
		shadow: 'xs',
		m: 'sm',
	};

	if (isPage) {
		defaultOptions.bd = '1px solid var(--mantine-color-gray-2)';
	} else {
		defaultOptions.p = 'sm';
	}

	const router = createStoryRouter(() => children);

	return (
		<Paper {...defaultOptions} {...paperOptions}>
			<RouterProvider router={router} />
		</Paper>
	);
}
