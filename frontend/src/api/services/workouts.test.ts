import { describe, expect, it } from 'vitest';
import createWorkout from '@/api/fixtures/workout';
import { Workout } from '@/types/domain';
import { TrackWorkoutsPageData, unwrapPages } from './workouts';

function page(workouts: Array<Workout>, overrides: Partial<TrackWorkoutsPageData> = {}) {
	return {
		data: { Workouts: workouts, NextCursor: '', Total: workouts.length, Planned: 0, ...overrides },
	};
}

function ids(result: ReturnType<typeof unwrapPages>) {
	return result.pages.flatMap((p) => p.Workouts.map((w) => w.ID));
}

describe('unwrapPages', () => {
	it('разворачивает конверт, сохраняя порядок страниц', () => {
		const first = createWorkout({ ID: 'a', Date: '2026-08-25' });
		const second = createWorkout({ ID: 'b', Date: '2026-08-23' });

		const result = unwrapPages({
			pages: [page([first]), page([second])],
			pageParams: ['', 'cursor'],
		});

		expect(ids(result)).toEqual(['a', 'b']);
		expect(result.pages[0].Total).toBe(1);
		expect(result.pageParams).toEqual(['', 'cursor']);
	});

	it('убирает строку, попавшую в две страницы, оставляя нижнюю копию', () => {
		// тренировку перенесли под курсор, пока читали: она приехала дважды
		const stale = createWorkout({ ID: 'moved', Date: '2026-08-30' });
		const fresh = createWorkout({ ID: 'moved', Date: '2026-08-24' });
		const other = createWorkout({ ID: 'other', Date: '2026-08-22' });

		const result = unwrapPages({
			pages: [page([stale]), page([fresh, other])],
			pageParams: ['', 'cursor'],
		});

		expect(ids(result)).toEqual(['moved', 'other']);
		expect(result.pages[0].Workouts).toHaveLength(0);
		expect(result.pages[1].Workouts[0].Date).toBe('2026-08-24');
	});

	it('не трогает страницы без дублей', () => {
		const first = createWorkout({ ID: 'a' });
		const input = { pages: [page([first])], pageParams: [''] };

		const result = unwrapPages(input);

		// та же ссылка: react-query сверяет результат select глубоким сравнением
		expect(result.pages[0]).toBe(input.pages[0].data);
	});

	it('переживает пустой список', () => {
		const result = unwrapPages({ pages: [page([])], pageParams: [''] });

		expect(ids(result)).toEqual([]);
	});
});
