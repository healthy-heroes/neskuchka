import { IconCheck, IconChevronLeft, IconPlayerPlayFilled } from '@tabler/icons-react';
import { useQuery } from '@tanstack/react-query';
import { Box, Button, Card, Container, Group, Skeleton, Text, Title } from '@mantine/core';
import { useApi } from '@/api/hooks';
import { WorkoutSection } from '@/types/domain';
import { formatIsoDate, formatWeekday, isToday } from '@/utils/dates';
import { RouteLink } from '../RouteLink/RouteLink';
import { WorkoutSections } from '../WorkoutSections/WorkoutSections';
import { WorkoutSectionsSkeleton } from '../WorkoutSections/WorkoutSectionsSkeleton';
import classes from './WorkoutView.module.css';

interface WorkoutViewProps {
	workoutId: string;
}

export function WorkoutView({ workoutId }: WorkoutViewProps) {
	const { workouts } = useApi();

	//todo: handle errors
	const { data, isPending } = useQuery(workouts.getWorkoutQuery(workoutId));

	if (isPending || !data) {
		return (
			<Container px="xl" pt="lg" className={classes.page}>
				<Box my="lg">
					<Skeleton height={52} width={420} />
				</Box>

				<div className={classes.layout}>
					<Card className={classes.sectionsCard} px="xl" py="lg">
						<WorkoutSectionsSkeleton rows={4} />
					</Card>
				</div>
			</Container>
		);
	}

	const workout = data.Workout;

	return (
		<Container px="xl" pt="lg" className={classes.page}>
			<RouteLink to="/workouts" fz="sm" fw={600} c="gray.7" underline="never">
				<Group gap="xs" wrap="nowrap">
					<IconChevronLeft size={15} />
					<span>Нескучный спорт</span>
				</Group>
			</RouteLink>

			<Group component="header" justify="space-between" align="flex-end" gap="xl" mt="sm" mb="lg">
				<div>
					<Group gap="sm" mb="xs">
						{isToday(workout.Date) && (
							<Text
								span
								fz="xs"
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
						<Text span fz="sm" c="gray.7">
							{formatWeekday(workout.Date)}
						</Text>
					</Group>
					<Title order={1}>Тренировка {formatIsoDate(workout.Date)}</Title>
				</div>

				{/* Прохождение и отметка выполнения — второй этап, кнопки пока нерабочие */}
				<Group gap="xs" wrap="nowrap">
					<Button size="md" disabled leftSection={<IconPlayerPlayFilled size={17} />}>
						Начать тренировку
					</Button>
					<Button size="md" disabled variant="default" leftSection={<IconCheck size={17} />}>
						Выполнено
					</Button>
				</Group>
			</Group>

			<div className={classes.layout}>
				<Card className={classes.sectionsCard} px="xl" py="lg">
					<WorkoutSections sections={workout.Sections} />
				</Card>

				<aside className={classes.side}>
					<Cheatsheet sections={workout.Sections} />
				</aside>
			</div>
		</Container>
	);
}

/** Шпаргалка — те же секции, ужатые до предписаний и названий. */
function Cheatsheet({ sections }: { sections: Array<WorkoutSection> }) {
	return (
		<Card padding={0}>
			<Text
				px="md"
				py="sm"
				ff="heading"
				fz="sm"
				fw={600}
				lts="0.08em"
				tt="uppercase"
				c="gray.8"
				className={classes.cheatsheetHead}
			>
				Шпаргалка
			</Text>

			<Box px="md" py="sm" className={classes.cheatsheetBody}>
				{sections.map((section, sectionIndex) => (
					<div key={`${section.Title}-${sectionIndex}`}>
						<Group justify="space-between" align="baseline" gap="xs" mb="xs">
							<Text span fz="xs" fw={700} lts="0.05em" tt="uppercase">
								{section.Title}
							</Text>
							<Text span fz="xs" fw={600} c="slate.7" ta="right">
								{section.Protocol.Title}
							</Text>
						</Group>

						{section.Exercises.length > 0 && (
							<ul className={classes.cheatsheetList}>
								{section.Exercises.map((exercise, exerciseIndex) => (
									<li key={`${exercise.Name}-${exerciseIndex}`} className={classes.cheatsheetItem}>
										<span className={classes.cheatsheetPrescription}>
											{exercise.Prescription.map((line) => (
												<Text key={line} span ff="heading" fz="sm" fw={600} c="copper.6">
													{line}
												</Text>
											))}
										</span>
										<Text span fz="sm" lh="xs" c="gray.8">
											{exercise.Name}
										</Text>
									</li>
								))}
							</ul>
						)}
					</div>
				))}
			</Box>
		</Card>
	);
}
