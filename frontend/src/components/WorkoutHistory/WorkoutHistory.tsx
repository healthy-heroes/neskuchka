import dayjs from 'dayjs';
import { IconCheck, IconChevronRight } from '@tabler/icons-react';
import { Workout } from '@/types/domain';
import { isWorkoutDone } from '@/utils/completion';
import { formatIsoDateShort, formatWeekday, isToday } from '@/utils/dates';
import { RouteLink } from '../RouteLink/RouteLink';
import classes from './WorkoutHistory.module.css';

export interface WorkoutHistoryProps {
	/** Тренировки трека, свежие первыми. */
	workouts: Array<Workout>;
}

interface Week {
	start: dayjs.Dayjs;
	items: Array<Workout>;
}

/**
 * WorkoutHistory — прошедшие тренировки, сгруппированные по неделям.
 *
 * Порядок везде одинаковый: недели от новых к старым, внутри недели тоже.
 * Сегодняшняя тренировка сюда не попадает — она в карточке выше. Будущие
 * не показываем вовсе: публикуется только прошлое и сегодня.
 */
export function WorkoutHistory({ workouts }: WorkoutHistoryProps) {
	const weeks = groupByWeek(workouts);

	// Сегодняшняя тренировка живёт в карточке выше, поэтому в пустой неделе
	// текст зависит от того, есть ли она вообще
	const hasToday = workouts.some((workout) => isToday(workout.Date));

	if (weeks.length === 0) {
		return null;
	}

	return (
		<div>
			{weeks.map((week) => (
				<section key={week.start.format('YYYY-MM-DD')}>
					<div className={classes.weekHead}>
						<span className={classes.weekTitle}>{weekTitle(week)}</span>
						<span className={classes.weekLine} />
						<span className={classes.weekMeta}>{weekMeta(week)}</span>
					</div>

					<div className={classes.rows}>
						{week.items.length === 0 ? (
							<p className={classes.weekEmpty}>
								{hasToday
									? 'Кроме сегодняшней на этой неделе тренировок пока не было'
									: 'На этой неделе тренировок пока не было'}
							</p>
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
				<div className={classes.rowDate}>{formatIsoDateShort(workout.Date)}</div>
				<div className={classes.rowWeekday}>{formatWeekday(workout.Date)}</div>
			</div>

			<div className={classes.rowSummary}>
				{workout.Sections.map((section, index) => (
					<span key={`${section.Title}-${index}`}>
						{index > 0 && <span className={classes.summarySep}>·</span>}
						{[section.Title, section.Protocol.Title].filter(Boolean).join(' · ')}
					</span>
				))}
			</div>

			{done ? (
				<div className={classes.rowDone}>
					<IconCheck size={16} stroke={2.5} />
					<span>Выполнено</span>
				</div>
			) : (
				<div className={classes.rowMissed}>Пропущено</div>
			)}

			<IconChevronRight size={18} className={classes.rowChevron} />
		</RouteLink>
	);
}

function groupByWeek(workouts: Array<Workout>): Array<Week> {
	const today = dayjs().startOf('day');

	const past = workouts.filter((workout) => dayjs(workout.Date).isBefore(today, 'day'));

	const weeks: Array<Week> = [];
	for (const workout of past) {
		const start = dayjs(workout.Date).startOf('week');
		const week = weeks.find((item) => item.start.isSame(start, 'day'));

		if (week) {
			week.items.push(workout);
		} else {
			weeks.push({ start, items: [workout] });
		}
	}

	// Текущая неделя показывается всегда: без неё непонятно, что на этой неделе
	// кроме сегодняшней ничего не было.
	const currentStart = today.startOf('week');
	if (!weeks.some((week) => week.start.isSame(currentStart, 'day'))) {
		weeks.unshift({ start: currentStart, items: [] });
	}

	return weeks;
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
