import axios from 'axios';

export default class ApiClient {
	private readonly apiPath: string;

	constructor() {
		this.apiPath = '/api/v1';
	}

	/**
	 * Helper for making GET requests to the API
	 * wraps axios.get
	 */
	async get<T>(path: string): Promise<T> {
		//todo: handle errors
		return axios.get<T>(`${this.apiPath}/${path}`).then((response) => response.data);
	}

	/**
	 * Helper for making PUT requests to the API
	 * wraps axios.put
	 */
	async put<T>(path: string, data: any): Promise<T> {
		return axios.put<T>(`${this.apiPath}/${path}`, data).then((response) => response.data);
	}
}
