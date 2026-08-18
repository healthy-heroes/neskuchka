import { Group, Skeleton } from '@mantine/core';
import classes from './WorkoutSections.module.css';

export interface WorkoutSectionsSkeletonProps {
	/** Сколько строк упражнений изобразить. */
	rows?: number;
}

/**
 * Скелетон протокола: тот же грид, что у настоящего, поэтому колонка
 * предписаний на месте и содержимое не прыгает, когда данные приедут.
 */
export function WorkoutSectionsSkeleton({ rows = 3 }: WorkoutSectionsSkeletonProps) {
	return (
		<div className={classes.sections}>
			<Group className={classes.fullRow} align="center" gap="sm">
				<Skeleton height={20} width={140} />
				<Skeleton height={18} width={90} radius="xl" />
			</Group>

			{Array.from({ length: rows }, (_, index) => (
				<div key={index} className={classes.section}>
					<Skeleton height={16} width={54} my="xs" />
					<Skeleton height={16} width={`${55 + ((index * 13) % 25)}%`} my="xs" />
					<Skeleton height={16} width={80} my="xs" />
				</div>
			))}
		</div>
	);
}
