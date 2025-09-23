import Api from './api';
import { WorkoutsQueries } from './workouts';

export class ApiQueries {
	readonly workouts: WorkoutsQueries;

	constructor(private readonly api: Api) {
		this.api = api;

		this.workouts = new WorkoutsQueries(this.api);
	}
}
