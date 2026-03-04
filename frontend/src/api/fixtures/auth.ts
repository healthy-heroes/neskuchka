import { AuthKeys } from '../services/auth';

type AuthServiceMockOptions = {
	loginError?: Error;
	logoutError?: Error;
	confirmLoginFn?: (token: string) => Promise<null>;
};

/**
 * Creates a mock AuthService for testing
 */
export function createAuthServiceMock(options: AuthServiceMockOptions = {}) {
	const { loginError, logoutError, confirmLoginFn } = options;

	return {
		loginMutation: () => ({
			mutationFn: async (_email: string): Promise<void> => {
				if (loginError) {
					throw loginError;
				}
			},
		}),

		logoutMutation: () => ({
			mutationFn: async (): Promise<void> => {
				if (logoutError) {
					throw logoutError;
				}
			},
		}),

		confirmLoginQuery: (token: string) => ({
			queryKey: AuthKeys.confirm(token),
			queryFn: async (): Promise<null> => {
				if (confirmLoginFn) {
					return confirmLoginFn(token);
				}
				return null;
			},
			retry: false,
		}),
	};
}
