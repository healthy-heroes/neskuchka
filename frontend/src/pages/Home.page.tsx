import { useEffect } from 'react';
import { useNavigate } from '@tanstack/react-router';
import { useAuth } from '@/auth/hooks';
import { PageSkeleton } from '@/components/PageSkeleton/PageSkeleton';

export function HomePage() {
	const { isAuthenticated, isLoading } = useAuth();
	const navigate = useNavigate();

	useEffect(() => {
		if (!isLoading) {
			navigate({ to: isAuthenticated ? '/workouts' : '/welcome' });
		}
	}, [isLoading, isAuthenticated, navigate]);

	return <PageSkeleton />;
}
