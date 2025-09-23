import axios from 'axios';

export class Api {
	private readonly apiPath: string;

	constructor() {
		this.apiPath = '/api/v1';
	}

	/**
	 * Base helper for getting data from the API
	 * wraps axios.get, handles errors and transforms the response for convenience
	 */
	async get<T>(path: string): Promise<T> {
		//todo: handle errors
		return axios.get<T>(`${this.apiPath}/${path}`).then((response) => response.data);
	}
}

export default Api;
