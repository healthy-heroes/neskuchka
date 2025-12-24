import { Outlet } from '@tanstack/react-router';
import { Header } from '@/components/Header/Header';

export function AuthPage() {
	return (
		<>
			<Header />
			<Outlet />
		</>
	);
}
