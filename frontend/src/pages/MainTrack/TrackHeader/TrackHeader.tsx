import { Avatar, Group, Text, Title } from '@mantine/core';
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
				<Text ff="heading" fz={15} fw={500} lts="0.12em" tt="uppercase" c="copper.6" mb={6}>
					Трек
				</Text>
				<Title order={1} fz={44} lh={1.05} lts="0.01em" mb={10}>
					{Name}
				</Title>
				<Text fz={16} lh={1.55} c="gray.8" textWrap="pretty" className={classes.description}>
					{Description}
				</Text>

				{Author?.Name && (
					<Group gap={8} mt={16}>
						<Avatar size={26} radius="xl" name={Author.Name} color="copper" />
						<Text span fz={14} c="gray.8">
							{Author.Name}
						</Text>
					</Group>
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
