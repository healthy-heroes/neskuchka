import { Box, Card, Group, Skeleton } from '@mantine/core';
import { WorkoutProtocolSkeleton } from '../WorkoutProtocol/WorkoutProtocolSkeleton';
import classes from './FeaturedWorkout.module.css';

/** Скелетон карточки в фокусе: тёмная шапка на месте, чтобы экран не перекрашивался. */
export function FeaturedWorkoutSkeleton() {
	return (
		<Card component="article" padding={0} className={classes.card}>
			<Group component="header" justify="space-between" px={28} py={18} bg="slate.7">
				<Group align="baseline" gap={16}>
					<Skeleton height={14} width={80} />
					<Skeleton height={30} width={190} />
					<Skeleton height={14} width={70} />
				</Group>
			</Group>

			<Box px={28} pt={10} pb={22}>
				<WorkoutProtocolSkeleton />
			</Box>

			<Box px={28} pb={28}>
				<Skeleton height={44} width={210} radius="md" />
			</Box>
		</Card>
	);
}
