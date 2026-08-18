import { Skeleton } from '@mantine/core';
import { WorkoutProtocolSkeleton } from '../WorkoutProtocol/WorkoutProtocolSkeleton';
import classes from './FeaturedWorkout.module.css';

/** Скелетон карточки в фокусе: тёмная шапка на месте, чтобы экран не перекрашивался. */
export function FeaturedWorkoutSkeleton() {
	return (
		<article className={classes.card}>
			<header className={classes.head}>
				<div className={classes.headMain}>
					<Skeleton height={14} width={80} />
					<Skeleton height={30} width={190} />
					<Skeleton height={14} width={70} />
				</div>
			</header>

			<div className={classes.body}>
				<WorkoutProtocolSkeleton />
			</div>

			<div className={classes.actions}>
				<Skeleton height={44} width={210} radius="md" />
			</div>
		</article>
	);
}
