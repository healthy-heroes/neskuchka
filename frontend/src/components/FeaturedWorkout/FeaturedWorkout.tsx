import { IconCheck, IconPlayerPlayFilled } from '@tabler/icons-react';
import { Box, Button, Card, Group, Text, Title } from '@mantine/core';
import { Workout } from '@/types/domain';
import { formatIsoDate, formatWeekday, isToday } from '@/utils/dates';
import { RouteLink } from '../RouteLink/RouteLink';
import { WorkoutSections } from '../WorkoutSections/WorkoutSections';
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
				px="xl"
				py="md"
				bg={today ? 'slate.7' : 'gray.0'}
				className={today ? undefined : classes.headPast}
			>
				<Group align="baseline" gap="md">
					<Text
						span
						ff="heading"
						fz="sm"
						fw={500}
						lts="0.14em"
						tt="uppercase"
						c={today ? 'copper.3' : 'copper.6'}
					>
						{today ? 'Сегодня' : 'Последняя тренировка'}
					</Text>
					<Title order={3} c={today ? 'white' : 'gray.9'}>
						{formatIsoDate(workout.Date)}
					</Title>
					<Text span fz="sm" c={today ? 'slate.2' : 'gray.7'}>
						{formatWeekday(workout.Date)}
					</Text>
				</Group>
			</Group>

			<Box px="xl" pt="xs" pb="lg">
				<WorkoutSections sections={workout.Sections} />
			</Box>

			{/* Прохождение и отметка выполнения — второй этап, кнопки пока нерабочие */}
			<Group px="xl" pb="xl" gap="sm">
				<Button h={44} fz="sm" disabled leftSection={<IconPlayerPlayFilled size={16} />}>
					Начать тренировку
				</Button>
				<Button h={44} fz="sm" disabled variant="default" leftSection={<IconCheck size={16} />}>
					Отметить выполненной
				</Button>

				<RouteLink
					to="/workouts/$workoutId"
					params={{ workoutId: workout.ID }}
					ml="auto"
					fz="sm"
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
