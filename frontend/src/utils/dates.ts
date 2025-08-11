import dayjs from 'dayjs';

export function formatIsoDate(isoDate: string) {
	let formattedDate = dayjs(isoDate).format('D MMMM');

	const date = new Date(isoDate);
	if (date.getFullYear() !== new Date().getFullYear()) {
		formattedDate += ` ${date.getFullYear()}`;
	}

	return formattedDate;
}
