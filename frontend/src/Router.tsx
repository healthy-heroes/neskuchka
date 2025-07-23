import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { HomePage } from './pages/Home.page';
import { LandingPage } from './pages/Landing/Landing.page';

const router = createBrowserRouter([
	{
		path: '/',
		element: <HomePage />,
	},
	{
		path: '/welcome',
		element: <LandingPage />,
	},
]);

export function Router() {
	return <RouterProvider router={router} />;
}
