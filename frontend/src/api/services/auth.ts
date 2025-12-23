import { UseMutationOptions, UseQueryOptions } from '@tanstack/react-query';
import { User } from '@/types/domain';
import Service from './service';

export const AuthKeys = {
	user: ['auth', 'user'] as const,
};

export class AuthService extends Service {
	/**
	 * Get the current user if authenticated else returns error with status 401
	 */
	getUserQuery(): UseQueryOptions<User> {
		return {
			queryKey: AuthKeys.user,
			queryFn: () => this.api.get<User>(`auth/user`),
			retry: false,
			staleTime: 5 * 60 * 1000, // 5 minutes
		};
	}

	loginMutation(): UseMutationOptions<void, Error, string> {
		return {
			mutationFn: (email: string) => this.api.post<void>(`auth/login`, { email }),
		};
	}

	confirmLoginMutation(): UseMutationOptions<void, Error, string> {
		return {
			mutationFn: (token: string) => this.api.post<void>(`auth/login/confirm`, { token }),
		};
	}

	logoutMutation(): UseMutationOptions<void, Error, void> {
		return {
			mutationFn: () => this.api.post<void>(`auth/logout`, null),
		};
	}
}
