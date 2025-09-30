import ApiClient from '../client';

class Service {
	constructor(protected readonly api: ApiClient) {
		this.api = api;
	}
}

export default Service;
