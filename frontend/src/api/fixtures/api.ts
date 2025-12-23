import ApiClient from '../client';
import ApiService from '../service';
import { AuthService } from '../services/auth';
import { WorkoutsService } from '../services/workouts';

/**
 * Mock ApiClient that throws on unexpected calls
 */
export class ApiClientMock extends ApiClient {
	get<T>(): Promise<T> {
		throw new Error('ApiClientMock.get: Not implemented');
	}

	put<T>(): Promise<T> {
		throw new Error('ApiClientMock.put: Not implemented');
	}

	post<T>(): Promise<T> {
		throw new Error('ApiClientMock.post: Not implemented');
	}
}

/**
 * Mock services that throw on unexpected calls
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
 */
export function createApiServiceMock(options: ApiServiceMockOptions = {}): ApiService {
	const strictAuth = createStrictServiceMock<AuthService>('AuthService');
	const strictWorkouts = createStrictServiceMock<WorkoutsService>('WorkoutsService');

	return {
		auth: { ...strictAuth, ...options.auth },
		workouts: { ...strictWorkouts, ...options.workouts },
	} as ApiService;
}

// Legacy exports for backwards compatibility
export const ApiMock = ApiClientMock;
export class ApiServiceMock extends ApiService {
	constructor() {
		super(new ApiClientMock());
	}
}
