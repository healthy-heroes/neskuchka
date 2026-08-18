import { IconCheck, IconPlayerPlayFilled } from '@tabler/icons-react';
import { Button, Group } from '@mantine/core';
import { Workout } from '@/types/domain';
import { formatIsoDate, formatWeekday } from '@/utils/dates';
import { RouteLink } from '../RouteLink/RouteLink';
import { WorkoutProtocol } from '../WorkoutProtocol/WorkoutProtocol';
import classes from './TodayWorkout.module.css';

export interface TodayWorkoutProps {
	workout: Workout;
}

/**
 * TodayWorkout — карточка сегодняшней тренировки, главный блок трека.
 *
 * Всегда развёрнута: человек заходит, видит, что делать сегодня, и начинает.
 * Аккордеона здесь нет, остальные тренировки живут строками в истории.
 */
export function TodayWorkout({ workout }: TodayWorkoutProps) {
	return (
		<article className={classes.card}>
			<header className={classes.head}>
				<div className={classes.headMain}>
					<span className={classes.eyebrow}>Сегодня</span>
					<span className={classes.date}>{formatIsoDate(workout.Date)}</span>
					<span className={classes.weekday}>{formatWeekday(workout.Date)}</span>
				</div>
			</header>

			<div className={classes.body}>
				<WorkoutProtocol sections={workout.Sections} />
			</div>

			{/* Прохождение и отметка выполнения — второй этап, кнопки пока нерабочие */}
			<div className={classes.actions}>
				<Group gap={12} align="center" w="100%">
					<Button h={44} fz={15} disabled leftSection={<IconPlayerPlayFilled size={16} />}>
						Начать тренировку
					</Button>
					<Button h={44} fz={15} disabled variant="default" leftSection={<IconCheck size={16} />}>
						Отметить выполненной
					</Button>

					<RouteLink
						to="/workouts/$workoutId"
						params={{ workoutId: workout.ID }}
						className={classes.breakdown}
						underline="never"
					>
						Разбор упражнений →
					</RouteLink>
				</Group>
			</div>
		</article>
	);
}

/** Состояние «на сегодня тренировки нет»: показываем последнюю опубликованную. */
export function NoWorkoutToday({ last }: { last?: Workout }) {
	return (
		<article className={classes.card}>
			<div className={classes.empty}>
				<p className={classes.emptyText}>
					На сегодня тренировки нет — новая появится, когда её опубликуют
				</p>

				{last && (
					<RouteLink
						to="/workouts/$workoutId"
						params={{ workoutId: last.ID }}
						className={classes.emptyLink}
						underline="never"
					>
						Открыть {formatIsoDate(last.Date)} →
					</RouteLink>
				)}
			</div>
		</article>
	);
}
