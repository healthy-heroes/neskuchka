import { Navigate } from '@tanstack/react-router';
import { useAuth } from '@/auth/hooks';
import { PageSkeleton } from '@/components/PageSkeleton/PageSkeleton';

export function HomePage() {
	const { isAuthenticated, isPending } = useAuth();

	if (isPending) {
		return <PageSkeleton />;
	}

	return isAuthenticated ? <Navigate to="/workouts" /> : <Navigate to="/welcome" />;
}
