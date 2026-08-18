import { IconCheck, IconPlayerPlayFilled } from '@tabler/icons-react';
import clsx from 'clsx';
import { Button, Group } from '@mantine/core';
import { Workout } from '@/types/domain';
import { formatIsoDate, formatWeekday, isToday } from '@/utils/dates';
import { RouteLink } from '../RouteLink/RouteLink';
import { WorkoutProtocol } from '../WorkoutProtocol/WorkoutProtocol';
import classes from './FeaturedWorkout.module.css';

export interface FeaturedWorkoutProps {
	workout: Workout;
}

/**
 * FeaturedWorkout — тренировка в фокусе трека, всегда развёрнутая.
 *
 * Обычно это сегодняшняя: человек заходит, видит, что делать сегодня, и начинает.
 * Если на сегодня тренировки нет, разворачиваем последнюю опубликованную —
 * иначе главный экран пустеет на ровном месте. Тёмная шапка при этом остаётся
 * привилегией сегодняшней: она и означает «это на сегодня».
 */
export function FeaturedWorkout({ workout }: FeaturedWorkoutProps) {
	const today = isToday(workout.Date);

	return (
		<article className={classes.card}>
			<header className={clsx(classes.head, !today && classes.headPast)}>
				<div className={classes.headMain}>
					<span className={classes.eyebrow}>{today ? 'Сегодня' : 'Последняя тренировка'}</span>
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
