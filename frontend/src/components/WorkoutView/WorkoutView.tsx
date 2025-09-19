import { IconArrowLeft } from '@tabler/icons-react';
import { Box, Divider, Grid, Image, List, Text, Title } from '@mantine/core';
import { useApi } from '@/api/provider';
import { formatIsoDate } from '@/utils/dates';
import { RouteLink } from '../RouteLink/RouteLink';
import { WorkoutCardSkeleton } from '../WorkoutCard/WorkoutCardSkeleton';
import classes from './WorkoutView.module.css';

interface WorkoutViewProps {
	workoutId: string;
}

export function WorkoutView({ workoutId }: WorkoutViewProps) {
	const apiService = useApi();
	const { data, loading, error } = apiService.getWorkout(workoutId);

	if (loading || error || !data) {
		return (
			<Box p="md">
				<Title order={2} my="md">
					Тренировка
				</Title>

				{(loading || error) && <WorkoutCardSkeleton cardProps={{ mb: 'xl' }} />}
			</Box>
		);
	}

	const workout = data.Workout;

	return (
		<Box p="md">
			<RouteLink to="..">
				<IconArrowLeft className={classes.linkIcon} size={14} /> к списку тренировок
			</RouteLink>

			<Title order={2} size="h1" mb="md">
				Тренировка от {formatIsoDate(workout.Date)}
			</Title>

			<Grid>
				<Grid.Col span={{ base: 12, xs: 7 }} p="md">
					{workout.Sections.map((section, index) => {
						const key = `${workout.ID}-${index}`;
						return (
							<div key={key}>
								<Title order={4} className={classes.sectionTitle}>
									{section.Title}
								</Title>
								<div>
									<b>{section.Protocol.Title}</b>
									{section.Protocol.Description && <span>{section.Protocol.Description}</span>}
									<List withPadding>
										{section.Exercises.map((e, index) => {
											return <List.Item key={`${key}-${index}`}>{e.Description}</List.Item>;
										})}
									</List>
								</div>

								<Divider my="md" />
							</div>
						);
					})}

					{workout.Notes && <Text>{workout.Notes}</Text>}
				</Grid.Col>
				<Grid.Col span={{ base: 12, xs: 5 }} p="md">
					<Image src="https://placehold.co/400x600?text=video" alt="Workout" />
				</Grid.Col>
			</Grid>
		</Box>
	);
}
