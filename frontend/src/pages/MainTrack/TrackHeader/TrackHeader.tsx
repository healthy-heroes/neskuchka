import { Avatar, Box, Group, Text, Title } from '@mantine/core';
import { TrackData } from '@/api/services/workouts';
import classes from './TrackHeader.module.css';

interface TrackHeaderProps {
	track: TrackData;
}

export function TrackHeader({ track }: TrackHeaderProps) {
	return (
		<Box p="lg" py="xl" bg="blue.1">
			<Title order={1} mb="xl">
				{track.Track.Name}
			</Title>

			<Text size="xl">Тренируйтесь с нами — где бы вы ни находились!</Text>
			<Text size="xl">
				Идеальная программа, чтобы поддерживать форму дома, без специального оборудования.
			</Text>

			<Title order={6} mt="xl" c="gray.6" className={classes.authorTitle}>
				Автор
			</Title>
			<Group gap="xs">
				<Avatar size="sm" src="/img/avatar.jpg" />
				<Text size="md" span>
					Илья Карягин
				</Text>
			</Group>
		</Box>
	);
}
