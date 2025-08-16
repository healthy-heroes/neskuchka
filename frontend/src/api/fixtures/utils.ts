export function createIncrementalIdFactory() {
	let ID = 1;
	return () => ID++;
}
