import { Navigate } from '@tanstack/react-router';
import { useAuth } from '@/auth/hooks';
import { PageSkeleton } from '@/components/PageSkeleton/PageSkeleton';

export function HomePage() {
	const { isAuthenticated, isLoading } = useAuth();

	if (isLoading) {
		return <PageSkeleton />;
	}

	if (isAuthenticated) {
		return <Navigate to="/workouts" />;
	}

	return <Navigate to="/welcome" />;
}
