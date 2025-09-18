import { IconArrowRight } from '@tabler/icons-react';
import { Button, Card, CardProps, Divider, Grid, Image, List, Title } from '@mantine/core';
import { Workout } from '@/types/domain';
import { formatIsoDate } from '@/utils/dates';
import classes from './WorkoutCard.module.css';

export interface WorkoutCardProps {
	cardProps?: CardProps;

	workout: Workout;
}

export function WorkoutCard({ cardProps, workout }: WorkoutCardProps) {
	return (
		<Card shadow="sm" p="lg" radius="md" withBorder {...cardProps}>
			<Grid>
				<Grid.Col visibleFrom="xs" span={5}>
					<Image src="https://placehold.co/400x600?text=video" alt="Workout" />
				</Grid.Col>
				<Grid.Col span={{ base: 12, xs: 7 }} p="md">
					<Title order={2} size="h1" mb="md" className={classes.title}>
						{formatIsoDate(workout.Date)}
					</Title>

					{workout.Sections.map((section, index) => {
						const key = `${workout.ID}-${index}`;
						return (
							<div key={key}>
								<Title order={4} className={classes.sectionTitle}>
									{section.Title}
								</Title>
								<div>
									<b>{section.Protocol.Title}</b>
									<List withPadding>
										{section.Exercises.map((e, index) => {
											return <List.Item key={`${key}-${index}`}>{e.ExerciseSlug}</List.Item>;
										})}
									</List>
								</div>

								<Divider my="md" />
							</div>
						);
					})}

					<Button
						variant="light"
						rightSection={<IconArrowRight size={14} />}
						component="a"
						href="#"
					>
						Подробности
					</Button>
				</Grid.Col>
			</Grid>
		</Card>
	);
}
