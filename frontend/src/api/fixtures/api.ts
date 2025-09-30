import Api from '../client';
import ApiService from '../service';

export class ApiMock extends Api {
	get<T>(): Promise<T> {
		throw new Error('Not implemented');
	}

	put<T>(): Promise<T> {
		throw new Error('Not implemented');
	}
}

export class ApiServiceMock extends ApiService {
	constructor() {
		super(new ApiMock());
	}
}
