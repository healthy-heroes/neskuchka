import { IconCheck, IconPlayerPlayFilled } from '@tabler/icons-react';
import { Box, Button, Card, Group, Text } from '@mantine/core';
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
		<Card component="article" padding={0} className={classes.card}>
			<Group
				component="header"
				justify="space-between"
				px={28}
				py={18}
				bg={today ? 'slate.7' : 'gray.0'}
				className={today ? undefined : classes.headPast}
			>
				<Group align="baseline" gap={16}>
					<Text
						span
						ff="heading"
						fz={14}
						fw={500}
						lts="0.14em"
						tt="uppercase"
						c={today ? undefined : 'copper.6'}
						className={today ? classes.eyebrowToday : undefined}
					>
						{today ? 'Сегодня' : 'Последняя тренировка'}
					</Text>
					<Text span ff="heading" fz={30} fw={600} tt="uppercase" c={today ? 'white' : 'gray.9'}>
						{formatIsoDate(workout.Date)}
					</Text>
					<Text span fz={14} c={today ? 'slate.2' : 'gray.7'}>
						{formatWeekday(workout.Date)}
					</Text>
				</Group>
			</Group>

			<Box px={28} pt={10} pb={22}>
				<WorkoutProtocol sections={workout.Sections} />
			</Box>

			{/* Прохождение и отметка выполнения — второй этап, кнопки пока нерабочие */}
			<Group px={28} pb={28} gap={12}>
				<Button h={44} fz={15} disabled leftSection={<IconPlayerPlayFilled size={16} />}>
					Начать тренировку
				</Button>
				<Button h={44} fz={15} disabled variant="default" leftSection={<IconCheck size={16} />}>
					Отметить выполненной
				</Button>

				<RouteLink
					to="/workouts/$workoutId"
					params={{ workoutId: workout.ID }}
					ml="auto"
					fz={15}
					fw={600}
					c="copper.6"
					underline="never"
				>
					Разбор упражнений →
				</RouteLink>
			</Group>
		</Card>
	);
}
