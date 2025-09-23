import ApiClient from './client';
import { WorkoutsService } from './services/workouts';

export default class ApiService {
	readonly workouts: WorkoutsService;

	constructor(private readonly api: ApiClient) {
		this.api = api;

		this.workouts = new WorkoutsService(this.api);
	}
}
