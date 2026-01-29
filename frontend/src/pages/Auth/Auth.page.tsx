import { useEffect } from 'react';
import { Outlet, useNavigate } from '@tanstack/react-router';
import { useAuth } from '@/auth/hooks';
import { Header } from '@/components/Header/Header';
import { PageSkeleton } from '@/components/PageSkeleton/PageSkeleton';

export function AuthPage() {
	const { isAuthenticated, isLoading } = useAuth();
	const navigate = useNavigate();

	useEffect(() => {
		if (!isLoading && isAuthenticated) {
			navigate({ to: '/' });
		}
	}, [isLoading, isAuthenticated, navigate]);

	if (isLoading || isAuthenticated) {
		return <PageSkeleton />;
	}

	return (
		<>
			<Header />
			<Outlet />
		</>
	);
}
