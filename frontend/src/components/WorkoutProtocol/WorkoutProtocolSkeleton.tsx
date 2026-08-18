import { Group, Skeleton } from '@mantine/core';
import classes from './WorkoutProtocol.module.css';

export interface WorkoutProtocolSkeletonProps {
	/** Сколько строк упражнений изобразить. */
	rows?: number;
}

/**
 * Скелетон протокола: тот же грид, что у настоящего, поэтому колонка
 * предписаний на месте и содержимое не прыгает, когда данные приедут.
 */
export function WorkoutProtocolSkeleton({ rows = 3 }: WorkoutProtocolSkeletonProps) {
	return (
		<div className={classes.protocol}>
			<Group className={classes.fullRow} align="baseline" gap={12}>
				<Skeleton height={20} width={140} />
				<Skeleton height={18} width={90} radius="xl" />
			</Group>

			{Array.from({ length: rows }, (_, index) => (
				<div key={index} className={classes.section}>
					<Skeleton height={16} width={54} my={11} />
					<Skeleton height={16} width={`${55 + ((index * 13) % 25)}%`} my={11} />
					<Skeleton height={16} width={80} my={11} />
				</div>
			))}
		</div>
	);
}
