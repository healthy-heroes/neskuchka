import { User } from '@/types/domain';
import { AuthKeys } from '../services/auth';

export const mockUser: User = {
	ID: 'user-123',
	Name: 'Test User',
};

type AuthServiceMockOptions = {
	user?: User | null;
	loginError?: Error;
	logoutError?: Error;
	confirmLoginFn?: (token: string) => Promise<void>;
};

type UserResponse = {
	data: User;
};

/**
 * Creates a mock AuthService for testing
 */
export function createAuthServiceMock(options: AuthServiceMockOptions = {}) {
	const { user = mockUser, loginError, logoutError, confirmLoginFn } = options;

	return {
		getUserQuery: () => ({
			queryKey: AuthKeys.user,
			queryFn: async (): Promise<UserResponse> => {
				if (user === null) {
					throw new Error('Unauthorized');
				}
				return { data: user };
			},
			retry: false,
		}),

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
			queryFn: async (): Promise<void> => {
				if (confirmLoginFn) {
					await confirmLoginFn(token);
				}
			},
			retry: false,
		}),
	};
}
