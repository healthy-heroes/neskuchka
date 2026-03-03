import { UseMutationOptions, UseQueryOptions } from '@tanstack/react-query';
import Service from './service';

export const AuthKeys = {
	confirm: (token: string) => ['auth', 'confirm', token] as const,
};

type LoginPayload = {
	email: string;
};

type ConfirmLoginPayload = {
	token: string;
};

export class AuthService extends Service {
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
