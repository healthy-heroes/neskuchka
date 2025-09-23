import Api from '../api';

export class ApiMock extends Api {
	get<T>(): Promise<T> {
		throw new Error('Not implemented');
	}
}
