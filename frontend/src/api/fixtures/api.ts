import ApiService from '../service';
import { AuthService } from '../services/auth';
import { UserService } from '../services/user';
import { WorkoutsService } from '../services/workouts';
import { createAuthServiceMock } from './auth';
import { createUserServiceMock } from './user';

/**
 * Creates a strict service mock using Proxy
 * Returns provided mock methods, throws error for unmocked ones
 */
function createStrictServiceMock<T extends object>(name: string, mockMethods: Partial<T> = {}): T {
	return new Proxy(mockMethods as T, {
		get(target, prop) {
			if (prop in target) {
				return target[prop as keyof T];
			}
			return () => {
				throw new Error(`${name}.${String(prop)}: Not implemented`);
			};
		},
	});
}

type ApiServiceMockOptions = {
	auth?: Partial<AuthService>;
	user?: Partial<UserService>;
	workouts?: Partial<WorkoutsService>;
};

/**
 * Creates a mock ApiService where unmocked methods throw errors
 * Auth and user services are mocked by default (unauthorized user)
 *
 * @example
 * const mock = createApiServiceMock({
 *   user: createUserServiceMock({ user: mockUser }),
 *   auth: createAuthServiceMock(),
 * });
 *
 * @example
 * const mock = createApiServiceMock({
 *   workouts: {
 *     getWorkoutQuery: () => ({ queryKey: ['workout'], queryFn: async () => workout }),
 *   },
 * });
 */
export function createApiServiceMock(options: ApiServiceMockOptions = {}): ApiService {
	const defaultAuth = createAuthServiceMock();
	const defaultUser = createUserServiceMock({ user: null });

	return {
		auth: createStrictServiceMock<AuthService>('AuthService', options.auth ?? defaultAuth),
		user: createStrictServiceMock<UserService>('UserService', options.user ?? defaultUser),
		workouts: createStrictServiceMock<WorkoutsService>('WorkoutsService', options.workouts),
	} as ApiService;
}
