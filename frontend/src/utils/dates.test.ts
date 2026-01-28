import dayjs from 'dayjs';
import { describe, expect, it } from 'vitest';
import { formatIsoDate } from './dates';

import 'dayjs/locale/ru';

dayjs.locale('ru');

describe('Utils: dates.formatIsoDate', () => {
	it('should return the correct date without year if it is current year', () => {
		const date = `${new Date().getFullYear()}-03-04`;
		const formattedDate = formatIsoDate(date);
		expect(formattedDate).toBe('4 марта');
	});

	it('should return the correct date with year', () => {
		const year = new Date().getFullYear() - 2;
		const date = `${year}-02-28`;
		const formattedDate = formatIsoDate(date);
		expect(formattedDate).toBe(`28 февраля ${year}`);
	});
});
