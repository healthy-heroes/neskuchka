import {
	createMemoryHistory,
	createRootRoute,
	createRouter,
	RouteComponent,
	RouterProvider,
} from '@tanstack/react-router';
import { Paper, PaperProps } from '@mantine/core';
import { ApiProvider } from '@/api/provider';
import ApiService from '@/api/service';

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

	apiService?: ApiService;
}

export function StoryPreview(props: StoryPreviewProps) {
	if (!props.apiService) {
		return getPageWrapper(props);
	}

	return <ApiProvider apiService={props.apiService}>{getPageWrapper(props)}</ApiProvider>;
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
