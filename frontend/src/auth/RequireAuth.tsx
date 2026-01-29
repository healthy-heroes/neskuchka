import { Navigate } from '@tanstack/react-router';
import { useAuth } from './hooks';

type RequireAuthProps = {
	children: React.ReactNode;
	loadingComponent: React.ReactNode;
	guestOnly?: boolean;
};

export function RequireAuth({ children, loadingComponent, guestOnly }: RequireAuthProps) {
	const { isAuthenticated, isLoading } = useAuth();

	if (isLoading) {
		return loadingComponent;
	}

	if (guestOnly) {
		isAuthenticated ? <Navigate to="/" /> : children;
	}

	return isAuthenticated ? children : <Navigate to="/login" />;
}
