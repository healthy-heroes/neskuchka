import ApiService from '../service';
import { AuthService } from '../services/auth';
import { WorkoutsService } from '../services/workouts';

/**
 * Creates a strict service mock using Proxy
 * Any unmocked method will throw an error with clear message
 */
function createStrictServiceMock<T extends object>(name: string): T {
	return new Proxy({} as T, {
		get(_, prop) {
			return () => {
				throw new Error(`${name}.${String(prop)}: Not implemented`);
			};
		},
	});
}

type ApiServiceMockOptions = {
	auth?: Partial<AuthService>;
	workouts?: Partial<WorkoutsService>;
};

/**
 * Creates a mock ApiService where unmocked methods throw errors
 *
 * @example
 * const mock = createApiServiceMock({
 *   auth: createAuthServiceMock({ user: mockUser }),
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
	const strictAuth = createStrictServiceMock<AuthService>('AuthService');
	const strictWorkouts = createStrictServiceMock<WorkoutsService>('WorkoutsService');

	return {
		auth: { ...strictAuth, ...options.auth },
		workouts: { ...strictWorkouts, ...options.workouts },
	} as ApiService;
}
