import { Avatar } from '@mantine/core';
import { TrackData } from '@/api/services/workouts';
import { TrackProgress } from '@/components/TrackProgress/TrackProgress';
import { TrackProgressSkeleton } from '@/components/TrackProgress/TrackProgressSkeleton';
import { Workout } from '@/types/domain';
import classes from './TrackHeader.module.css';

interface TrackHeaderProps {
	track: TrackData;

	/** Нужны блоку прогресса: он считает окно последних 30 дней. */
	workouts: Array<Workout>;

	/** Трек приезжает раньше списка тренировок, полосе пока нечего показывать. */
	loading?: boolean;
}

export function TrackHeader({ track, workouts, loading = false }: TrackHeaderProps) {
	const { Name, Description, Author } = track.Track;

	return (
		<header className={classes.header}>
			<div className={classes.about}>
				<div className={classes.eyebrow}>Трек</div>
				<h1 className={classes.title}>{Name}</h1>
				<p className={classes.description}>{Description}</p>

				{Author?.Name && (
					<div className={classes.author}>
						<Avatar size={26} radius="xl" name={Author.Name} color="copper" />
						<span>{Author.Name}</span>
					</div>
				)}
			</div>

			{loading ? (
				<TrackProgressSkeleton />
			) : (
				<TrackProgress workouts={workouts} total={workouts.length} />
			)}
		</header>
	);
}
