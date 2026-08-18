import dayjs from 'dayjs';

/** Короткая дата для плотных списков: «17 авг». */
export function formatIsoDateShort(isoDate: string) {
	// В русской локали короткий месяц идёт с точкой — в макете её нет
	return dayjs(isoDate).format('D MMM').replace(/\.$/, '');
}

/** День недели целиком: «четверг». */
export function formatWeekday(isoDate: string) {
	return dayjs(isoDate).format('dddd');
}

/** Опубликована ли тренировка: будущие на экранах не показываем вовсе. */
export function isPublished(isoDate: string) {
	return !dayjs(isoDate).isAfter(dayjs(), 'day');
}

/** Сегодняшняя ли дата — по календарному дню, без времени. */
export function isToday(isoDate: string) {
	return dayjs(isoDate).isSame(dayjs(), 'day');
}

export function formatIsoDate(isoDate: string) {
	let formattedDate = dayjs(isoDate).format('D MMMM');

	const date = new Date(isoDate);
	if (date.getFullYear() !== new Date().getFullYear()) {
		formattedDate += ` ${date.getFullYear()}`;
	}

	return formattedDate;
}
