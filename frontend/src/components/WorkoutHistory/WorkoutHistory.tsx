import dayjs from 'dayjs';
import { IconCheck, IconChevronRight } from '@tabler/icons-react';
import { Group, Text } from '@mantine/core';
import { Workout } from '@/types/domain';
import { isWorkoutDone } from '@/utils/completion';
import { formatIsoDateShort, formatWeekday, isToday } from '@/utils/dates';
import { RouteLink } from '../RouteLink/RouteLink';
import classes from './WorkoutHistory.module.css';

export interface WorkoutHistoryProps {
	/** Опубликованные тренировки трека, свежие первыми. */
	workouts: Array<Workout>;

	/** Тренировка из карточки выше — в истории она не повторяется. */
	featured?: Workout;
}

interface Week {
	start: dayjs.Dayjs;
	items: Array<Workout>;
}

/**
 * WorkoutHistory — прошедшие тренировки, сгруппированные по неделям.
 *
 * Порядок везде одинаковый: недели от новых к старым, внутри недели тоже.
 * Тренировка из карточки выше сюда не попадает, чтобы не двоиться.
 */
export function WorkoutHistory({ workouts, featured }: WorkoutHistoryProps) {
	const weeks = groupByWeek(workouts, featured);

	if (weeks.length === 0) {
		return null;
	}

	return (
		<div>
			{weeks.map((week) => (
				<section key={week.start.format('YYYY-MM-DD')}>
					<Group gap="sm" my={0} mt="xl" mb="sm">
						<Text span ff="heading" fz="md" fw={600} lts="0.1em" tt="uppercase" c="gray.8">
							{weekTitle(week)}
						</Text>
						<span className={classes.weekLine} />
						<Text span fz="xs" c="gray.7" className={classes.nowrap}>
							{weekMeta(week)}
						</Text>
					</Group>

					<div className={classes.rows}>
						{week.items.length === 0 ? (
							<Text component="p" my={0} px="lg" py="md" fz="sm" c="gray.7">
								{emptyWeekText(week, featured)}
							</Text>
						) : (
							week.items.map((workout) => <HistoryRow key={workout.ID} workout={workout} />)
						)}
					</div>
				</section>
			))}
		</div>
	);
}

function HistoryRow({ workout }: { workout: Workout }) {
	const done = isWorkoutDone(workout);

	return (
		<RouteLink
			to="/workouts/$workoutId"
			params={{ workoutId: workout.ID }}
			className={classes.row}
			underline="never"
		>
			<div>
				<Text ff="heading" fz="xl" fw={600} tt="uppercase">
					{formatIsoDateShort(workout.Date)}
				</Text>
				<Text fz="xs" c="gray.7">
					{formatWeekday(workout.Date)}
				</Text>
			</div>

			<Text fz="sm" c="gray.8">
				{workout.Sections.map((section, index) => (
					<span key={`${section.Title}-${index}`}>
						{index > 0 && <span className={classes.summarySep}>·</span>}
						{[section.Title, section.Protocol.Title].filter(Boolean).join(' · ')}
					</span>
				))}
			</Text>

			{done ? (
				<Group gap="xs" wrap="nowrap" c="slate.7">
					<IconCheck size={16} stroke={2.5} />
					<Text span fz="sm" fw={600}>
						Выполнено
					</Text>
				</Group>
			) : (
				<Text fz="sm" c="gray.6">
					Пропущено
				</Text>
			)}

			<IconChevronRight size={18} className={classes.rowChevron} />
		</RouteLink>
	);
}

function groupByWeek(workouts: Array<Workout>, featured?: Workout): Array<Week> {
	const rest = workouts.filter((workout) => workout.ID !== featured?.ID);

	const weeks: Array<Week> = [];
	for (const workout of rest) {
		const start = dayjs(workout.Date).startOf('week');
		const week = weeks.find((item) => item.start.isSame(start, 'day'));

		if (week) {
			week.items.push(workout);
		} else {
			weeks.push({ start, items: [workout] });
		}
	}

	// Текущая неделя показывается всегда: без неё непонятно, что на этой неделе
	// кроме показанной выше тренировки ничего не было.
	const currentStart = dayjs().startOf('week');
	if (!weeks.some((week) => week.start.isSame(currentStart, 'day'))) {
		weeks.unshift({ start: currentStart, items: [] });
	}

	return weeks;
}

/**
 * Пустая неделя объясняется по-разному: одно дело «кроме сегодняшней»,
 * другое — когда наверху висит позавчерашняя, третье — когда её там нет вовсе.
 */
function emptyWeekText(week: Week, featured?: Workout): string {
	const featuredHere = featured && dayjs(featured.Date).startOf('week').isSame(week.start, 'day');

	if (!featuredHere) {
		return 'На этой неделе тренировок пока не было';
	}

	return isToday(featured.Date)
		? 'Кроме сегодняшней на этой неделе тренировок пока не было'
		: 'Кроме показанной выше на этой неделе тренировок пока не было';
}

function weekTitle(week: Week): string {
	const currentStart = dayjs().startOf('week');

	if (week.start.isSame(currentStart, 'day')) {
		return 'Эта неделя';
	}

	if (week.start.isSame(currentStart.subtract(1, 'week'), 'day')) {
		return 'Прошлая неделя';
	}

	return weekRange(week);
}

function weekMeta(week: Week): string {
	const parts: Array<string> = [];

	if (weekTitle(week) !== weekRange(week)) {
		parts.push(weekRange(week));
	}

	if (week.items.length > 0) {
		const done = week.items.filter(isWorkoutDone).length;
		parts.push(`${done} из ${week.items.length}`);
	}

	return parts.join(' · ');
}

function weekRange(week: Week): string {
	const end = week.start.add(6, 'day');

	// Неделя может пересечь границу месяца или года: «27 июля—2 августа»,
	// а не «27—2 августа», иначе начало диапазона читается как август.
	const startFormat = week.start.isSame(end, 'month') ? 'D' : 'D MMMM';
	const endFormat = end.isSame(dayjs(), 'year') ? 'D MMMM' : 'D MMMM YYYY';

	return `${week.start.format(startFormat)}—${end.format(endFormat)}`;
}
