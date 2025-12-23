import ApiClient from './client';
import { AuthService } from './services/auth';
import { WorkoutsService } from './services/workouts';

export default class ApiService {
	readonly auth: AuthService;
	readonly workouts: WorkoutsService;

	constructor(private readonly api: ApiClient) {
		this.api = api;

		this.auth = new AuthService(this.api);
		this.workouts = new WorkoutsService(this.api);
	}
}
