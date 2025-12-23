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
};

/**
 * Creates a mock AuthService for testing
 */
export function createAuthServiceMock(options: AuthServiceMockOptions = {}) {
	const { user = mockUser, loginError, logoutError } = options;

	return {
		getUserQuery: () => ({
			queryKey: AuthKeys.user,
			queryFn: async (): Promise<User> => {
				if (user === null) {
					throw new Error('Unauthorized');
				}
				return user;
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

		confirmLoginMutation: () => ({
			mutationFn: async (_token: string): Promise<void> => {},
		}),
	};
}
