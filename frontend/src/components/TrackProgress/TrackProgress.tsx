import clsx from 'clsx';
import { Text } from '@mantine/core';
import { Workout } from '@/types/domain';
import { isWorkoutDone } from '@/utils/completion';
import { formatIsoDate, isToday } from '@/utils/dates';
import classes from './TrackProgress.module.css';

const WINDOW_DAYS = 30;

export interface TrackProgressProps {
	/** Тренировки трека, свежие первыми. */
	workouts: Array<Workout>;
	/** Всего тренировок в треке за всё время. */
	total?: number;

	/** Ужатый вариант для боковой колонки. */
	compact?: boolean;
}

/**
 * TrackProgress — полоса «Последние 30 дней».
 *
 * Окно фиксированное: полоса показывает только последние 30 дней и не растёт
 * с возрастом трека. Полный объём трека уходит в подпись под полосой.
 *
 * Статус выполнения пока изображается заглушкой, см. utils/completion.
 */
export function TrackProgress({ workouts, total, compact = false }: TrackProgressProps) {
	const segments = workoutsInWindow(workouts).map((workout) => ({
		date: workout.Date,
		today: isToday(workout.Date),
		done: isWorkoutDone(workout),
	}));

	const doneCount = segments.filter((segment) => segment.done).length;
	const lastDone = segments.find((segment) => segment.done);

	// В окне пусто — показывать «0 из 0» и голый жёлоб незачем
	if (segments.length === 0) {
		return null;
	}

	return (
		<div className={clsx(classes.root, compact && classes.compact)}>
			<div className={classes.head}>
				<Text span fz={13} fw={600} lts="0.06em" tt="uppercase" c="gray.7">
					Последние 30 дней
				</Text>
				<Text span fz={14} c="gray.8" className={classes.count}>
					<Text span ff="heading" fz={compact ? 18 : 20} c="copper.6">
						{doneCount}
					</Text>{' '}
					из {segments.length}
					{!compact && ' тренировок'}
				</Text>
			</div>

			<div className={classes.bar}>
				{segments
					.slice()
					.reverse()
					.map((segment) => (
						<span
							key={segment.date}
							className={clsx(
								classes.segment,
								segment.today && classes.segmentToday,
								!segment.today && !segment.done && classes.segmentMissed
							)}
						/>
					))}
			</div>

			{(lastDone || total) && (
				<Text component="p" mt={10} fz={13} c="gray.7">
					{lastDone && <>последняя — {formatIsoDate(lastDone.date)}</>}
					{lastDone && total ? ' · ' : ''}
					{total ? <>всего в треке {total}</> : null}
				</Text>
			)}
		</div>
	);
}

function workoutsInWindow(workouts: Array<Workout>): Array<Workout> {
	const from = new Date();
	from.setDate(from.getDate() - WINDOW_DAYS);

	return workouts.filter((workout) => new Date(workout.Date) >= from);
}
