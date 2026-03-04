import ApiClient from './client';
import { AuthService } from './services/auth';
import { UserService } from './services/user';
import { WorkoutsService } from './services/workouts';

export default class ApiService {
	readonly auth: AuthService;
	readonly user: UserService;
	readonly workouts: WorkoutsService;

	constructor(private readonly api: ApiClient) {
		this.api = api;

		this.auth = new AuthService(this.api);
		this.user = new UserService(this.api);
		this.workouts = new WorkoutsService(this.api);
	}
}
