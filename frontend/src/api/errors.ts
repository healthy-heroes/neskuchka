export class HttpError extends Error {
	constructor(
		public status: number,
		public details: object
	) {
		super(JSON.stringify(details));
	}
}
