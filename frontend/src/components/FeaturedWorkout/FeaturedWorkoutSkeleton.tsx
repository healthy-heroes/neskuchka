import { Box, Card, Group, Skeleton } from '@mantine/core';
import { WorkoutSectionsSkeleton } from '../WorkoutSections/WorkoutSectionsSkeleton';
import classes from './FeaturedWorkout.module.css';

/** Скелетон карточки в фокусе: тёмная шапка на месте, чтобы экран не перекрашивался. */
export function FeaturedWorkoutSkeleton() {
	return (
		<Card component="article" padding={0} className={classes.card}>
			<Group component="header" justify="space-between" px="xl" py="md" bg="slate.7">
				<Group align="baseline" gap="md">
					<Skeleton height={14} width={80} />
					<Skeleton height={30} width={190} />
					<Skeleton height={14} width={70} />
				</Group>
			</Group>

			<Box px="xl" pt="xs" pb="lg">
				<WorkoutSectionsSkeleton />
			</Box>

			<Box px="xl" pb="xl">
				<Skeleton height={44} width={210} radius="md" />
			</Box>
		</Card>
	);
}
