import { User } from '@/types/domain';
import { UserKeys } from '../services/user';

export const mockUser: User = {
	ID: 'user-123',
	Name: 'Test User',
};

type UserServiceMockOptions = {
	user?: User | null;
};

type UserResponse = {
	data: User;
};

/**
 * Creates a mock UserService for testing
 */
export function createUserServiceMock(options: UserServiceMockOptions = {}) {
	const { user = mockUser } = options;

	return {
		getUserQuery: () => ({
			queryKey: UserKeys.me,
			queryFn: async (): Promise<UserResponse | null> => {
				if (user === null) {
					return null;
				}
				return { data: user };
			},
			retry: false,
		}),
	};
}
