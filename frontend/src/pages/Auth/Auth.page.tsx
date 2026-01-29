import { Navigate, Outlet } from '@tanstack/react-router';
import { useAuth } from '@/auth/hooks';
import { Header } from '@/components/Header/Header';
import { PageSkeleton } from '@/components/PageSkeleton/PageSkeleton';

export function AuthPage() {
	const { isAuthenticated, isLoading } = useAuth();

	if (isLoading) {
		return <PageSkeleton />;
	}

	if (isAuthenticated) {
		return <Navigate to="/" />;
	}

	return (
		<>
			<Header />
			<Outlet />
		</>
	);
}
