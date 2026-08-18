import { IconCheck, IconChevronLeft, IconPlayerPlayFilled } from '@tabler/icons-react';
import { useQuery } from '@tanstack/react-query';
import { Box, Button, Group } from '@mantine/core';
import { useApi } from '@/api/hooks';
import { Workout, WorkoutSection } from '@/types/domain';
import { formatIsoDate, formatWeekday, isToday } from '@/utils/dates';
import { RouteLink } from '../RouteLink/RouteLink';
import { TrackProgress } from '../TrackProgress/TrackProgress';
import { WorkoutCardSkeleton } from '../WorkoutCard/WorkoutCardSkeleton';
import { WorkoutProtocol } from '../WorkoutProtocol/WorkoutProtocol';
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
				<WorkoutCardSkeleton />
			</Box>
		);
	}

	const workout = data.Workout;

	return (
		<Box className={classes.page}>
			<RouteLink to="/workouts" className={classes.back} underline="never">
				<IconChevronLeft size={15} />
				<span>Нескучный спорт</span>
			</RouteLink>

			<header className={classes.title}>
				<div>
					<Group gap={12} mb={8}>
						{isToday(workout.Date) && <span className={classes.todayPill}>Сегодня</span>}
						<span className={classes.meta}>{formatWeekday(workout.Date)}</span>
					</Group>
					<h1 className={classes.heading}>Тренировка {formatIsoDate(workout.Date)}</h1>
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
			</header>

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
			<div className={classes.cheatsheetHead}>Шпаргалка</div>
			<div className={classes.cheatsheetBody}>
				{sections.map((section, sectionIndex) => (
					<div key={`${section.Title}-${sectionIndex}`}>
						<div className={classes.cheatsheetSection}>
							<span className={classes.cheatsheetTitle}>{section.Title}</span>
							<span className={classes.cheatsheetProtocol}>{section.Protocol.Title}</span>
						</div>

						{section.Exercises.length > 0 && (
							<ul className={classes.cheatsheetList}>
								{section.Exercises.map((exercise, exerciseIndex) => (
									<li key={`${exercise.Name}-${exerciseIndex}`} className={classes.cheatsheetItem}>
										<span className={classes.cheatsheetPrescription}>
											{exercise.Prescription.map((line) => (
												<span key={line}>{line}</span>
											))}
										</span>
										<span className={classes.cheatsheetName}>{exercise.Name}</span>
									</li>
								))}
							</ul>
						)}
					</div>
				))}
			</div>
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
			<div className={classes.cheatsheetBody}>
				{neighbours.map(({ label, workout }) => (
					<RouteLink
						key={workout.ID}
						to="/workouts/$workoutId"
						params={{ workoutId: workout.ID }}
						className={classes.neighbour}
						underline="never"
					>
						<span className={classes.neighbourDate}>{formatIsoDate(workout.Date)}</span>
						<span className={classes.neighbourLabel}>{label}</span>
					</RouteLink>
				))}
			</div>
		</div>
	);
}
