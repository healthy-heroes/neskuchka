import { UseMutationOptions, UseQueryOptions } from '@tanstack/react-query';
import axios from 'axios';
import { User } from '@/types/domain';
import Service from './service';

export const AuthKeys = {
	user: ['auth', 'user'] as const,
	confirm: (token: string) => ['auth', 'confirm', token] as const,
};

type UserResponse = {
	data: User;
};

type LoginPayload = {
	email: string;
};

type ConfirmLoginPayload = {
	token: string;
};

export class AuthService extends Service {
	/**
	 * Get the current user if authenticated, null if not authenticated (401)
	 */
	getUserQuery(): UseQueryOptions<UserResponse | null> {
		return {
			queryKey: AuthKeys.user,
			queryFn: async () => {
				try {
					return await this.api.get<UserResponse>(`auth/user`);
				} catch (error) {
					if (axios.isAxiosError(error) && error.response?.status === 401) {
						return null;
					}
					throw error;
				}
			},
			staleTime: 5 * 60 * 1000, // 5 minutes
		};
	}

	loginMutation(): UseMutationOptions<void, Error, string> {
		return {
			mutationFn: (email: string) => this.api.post<void, LoginPayload>(`auth/login`, { email }),
		};
	}

	/**
	 * Confirm a login attempt by providing a token
	 *
	 * @note Using query instead of mutation because need fight with double useEffect hooks in dev
	 * @note Returns null instead of void because TanStack Query v5 doesn't allow undefined as query data
	 */
	confirmLoginQuery(token: string): UseQueryOptions<null> {
		return {
			queryKey: AuthKeys.confirm(token),
			queryFn: async () => {
				await this.api.post<void, ConfirmLoginPayload>(`auth/login/confirm`, { token });
				return null;
			},
			retry: false,
		};
	}

	logoutMutation(): UseMutationOptions<void, Error, void> {
		return {
			mutationFn: () => this.api.post<void>(`auth/logout`, null),
		};
	}
}
