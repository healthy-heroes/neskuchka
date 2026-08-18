import { IconCheck, IconChevronLeft, IconPlayerPlayFilled } from '@tabler/icons-react';
import { useQuery } from '@tanstack/react-query';
import { Box, Button, Group, Skeleton, Text, Title } from '@mantine/core';
import { useApi } from '@/api/hooks';
import { Workout, WorkoutSection } from '@/types/domain';
import { formatIsoDate, formatWeekday, isToday } from '@/utils/dates';
import { RouteLink } from '../RouteLink/RouteLink';
import { TrackProgress } from '../TrackProgress/TrackProgress';
import { TrackProgressSkeleton } from '../TrackProgress/TrackProgressSkeleton';
import { WorkoutProtocol } from '../WorkoutProtocol/WorkoutProtocol';
import { WorkoutProtocolSkeleton } from '../WorkoutProtocol/WorkoutProtocolSkeleton';
import classes from './WorkoutView.module.css';

interface WorkoutViewProps {
	workoutId: string;
}

export function WorkoutView({ workoutId }: WorkoutViewProps) {
	const { workouts } = useApi();

	//todo: handle errors
	const { data, isPending } = useQuery(workouts.getWorkoutQuery(workoutId));
	const { data: trackWorkouts } = useQuery(workouts.getMainTrackWorkoutsQuery());

	if (isPending || !data) {
		return (
			<Box className={classes.page}>
				<Box my={24}>
					<Skeleton height={52} width={420} />
				</Box>

				<div className={classes.layout}>
					<div className={classes.protocolCard}>
						<WorkoutProtocolSkeleton rows={4} />
					</div>
					<TrackProgressSkeleton compact />
				</div>
			</Box>
		);
	}

	const workout = data.Workout;

	return (
		<Box className={classes.page}>
			<RouteLink to="/workouts" fz={14} fw={600} c="gray.7" underline="never">
				<Group gap={8} wrap="nowrap">
					<IconChevronLeft size={15} />
					<span>Нескучный спорт</span>
				</Group>
			</RouteLink>

			<Group component="header" justify="space-between" align="flex-end" gap={40} mt={14} mb={24}>
				<div>
					<Group gap={12} mb={8}>
						{isToday(workout.Date) && (
							<Text
								span
								fz={12}
								fw={700}
								lts="0.1em"
								tt="uppercase"
								c="white"
								bg="copper.6"
								className={classes.pill}
							>
								Сегодня
							</Text>
						)}
						<Text span fz={14} c="gray.7">
							{formatWeekday(workout.Date)}
						</Text>
					</Group>
					<Title order={1} lts="0.01em">
						Тренировка {formatIsoDate(workout.Date)}
					</Title>
				</div>

				{/* Прохождение и отметка выполнения — второй этап, кнопки пока нерабочие */}
				<Group gap={10} wrap="nowrap">
					<Button size="md" h={46} disabled leftSection={<IconPlayerPlayFilled size={17} />}>
						Начать тренировку
					</Button>
					<Button size="md" h={46} disabled variant="default" leftSection={<IconCheck size={17} />}>
						Выполнено
					</Button>
				</Group>
			</Group>

			<div className={classes.layout}>
				<div className={classes.protocolCard}>
					<WorkoutProtocol sections={workout.Sections} />
				</div>

				<aside className={classes.side}>
					<Cheatsheet sections={workout.Sections} />

					{trackWorkouts && (
						<TrackProgress
							compact
							workouts={trackWorkouts.Workouts}
							total={trackWorkouts.Workouts.length}
						/>
					)}

					{trackWorkouts && <Neighbours workouts={trackWorkouts.Workouts} current={workout} />}
				</aside>
			</div>
		</Box>
	);
}

/** Шпаргалка — тот же протокол, ужатый до предписаний и названий. */
function Cheatsheet({ sections }: { sections: Array<WorkoutSection> }) {
	return (
		<div className={classes.sideCard}>
			<Text
				px={18}
				py={13}
				ff="heading"
				fz={15}
				fw={600}
				lts="0.08em"
				tt="uppercase"
				c="gray.8"
				className={classes.cheatsheetHead}
			>
				Шпаргалка
			</Text>

			<Box px={18} py={14} className={classes.cheatsheetBody}>
				{sections.map((section, sectionIndex) => (
					<div key={`${section.Title}-${sectionIndex}`}>
						<Group justify="space-between" align="baseline" gap={10} mb={6}>
							<Text span fz={13} fw={700} lts="0.05em" tt="uppercase">
								{section.Title}
							</Text>
							<Text span fz={13} fw={600} c="slate.7" ta="right">
								{section.Protocol.Title}
							</Text>
						</Group>

						{section.Exercises.length > 0 && (
							<ul className={classes.cheatsheetList}>
								{section.Exercises.map((exercise, exerciseIndex) => (
									<li key={`${exercise.Name}-${exerciseIndex}`} className={classes.cheatsheetItem}>
										<span className={classes.cheatsheetPrescription}>
											{exercise.Prescription.map((line) => (
												<Text key={line} span ff="heading" fz={14} fw={600} c="copper.6">
													{line}
												</Text>
											))}
										</span>
										<Text span fz={14} lh={1.4} c="gray.8">
											{exercise.Name}
										</Text>
									</li>
								))}
							</ul>
						)}
					</div>
				))}
			</Box>
		</div>
	);
}

/** Соседние тренировки. Статуса выполнения пока нет — только даты. */
function Neighbours({ workouts, current }: { workouts: Array<Workout>; current: Workout }) {
	const index = workouts.findIndex((workout) => workout.ID === current.ID);
	if (index === -1) {
		return null;
	}

	const neighbours = [
		{ label: 'Следующая', workout: workouts[index - 1] },
		{ label: 'Предыдущая', workout: workouts[index + 1] },
	].filter((item) => item.workout);

	if (neighbours.length === 0) {
		return null;
	}

	return (
		<div className={classes.sideCard}>
			<Box px={18} py={14}>
				{neighbours.map(({ label, workout }) => (
					<RouteLink
						key={workout.ID}
						to="/workouts/$workoutId"
						params={{ workoutId: workout.ID }}
						c="inherit"
						underline="never"
					>
						<Group justify="space-between" align="baseline" gap={10} py={10}>
							<Text span fz={14} c="gray.8">
								{formatIsoDate(workout.Date)}
							</Text>
							<Text span fz={13} c="gray.7">
								{label}
							</Text>
						</Group>
					</RouteLink>
				))}
			</Box>
		</div>
	);
}
