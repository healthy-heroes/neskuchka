import clsx from 'clsx';
import { Skeleton } from '@mantine/core';
import classes from './TrackProgress.module.css';

const SEGMENTS = 12;

export interface TrackProgressSkeletonProps {
	compact?: boolean;
}

/** Скелетон полосы прогресса: столько же сегментов, сколько влезает в окно 30 дней. */
export function TrackProgressSkeleton({ compact = false }: TrackProgressSkeletonProps) {
	return (
		<div className={clsx(classes.root, compact && classes.compact)}>
			<div className={classes.head}>
				<Skeleton height={13} width={130} />
				<Skeleton height={14} width={70} />
			</div>

			<div className={classes.bar}>
				{Array.from({ length: SEGMENTS }, (_, index) => (
					<Skeleton key={index} height={10} radius={3} style={{ flex: 1 }} />
				))}
			</div>

			<Skeleton height={13} width="60%" mt={10} />
		</div>
	);
}
