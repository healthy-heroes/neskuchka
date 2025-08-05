import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { LandingPage } from './pages/Landing/Landing.page';
import { MainTrackPage } from './pages/MainTrack/MainTrack.page';

const router = createBrowserRouter([
	{
		path: '/',
		element: <LandingPage />,
	},
	{
		path: '/welcome',
		element: <LandingPage />,
	},
	{
		path: '/main',
		element: <MainTrackPage />,
	},
]);

export function Router() {
	return <RouterProvider router={router} />;
}
