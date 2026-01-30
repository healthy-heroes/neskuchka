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
	async get<TResponse>(path: string): Promise<TResponse> {
		//todo: handle errors
		return axios.get<TResponse>(`${this.apiPath}/${path}`).then((response) => response.data);
	}

	/**
	 * Helper for making PUT requests to the API
	 * wraps axios.put
	 */
	async put<TResponse, TPayload = unknown>(path: string, payload: TPayload): Promise<TResponse> {
		return axios
			.put<TResponse>(`${this.apiPath}/${path}`, payload)
			.then((response) => response.data);
	}

	/**
	 * Helper for making POST requests to the API
	 * wraps axios.post
	 */
	async post<TResponse, TPayload = unknown>(path: string, payload: TPayload): Promise<TResponse> {
		return axios
			.post<TResponse>(`${this.apiPath}/${path}`, payload)
			.then((response) => response.data);
	}
}
