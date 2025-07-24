import { Card, Container, Divider, Grid, Image, List, Text, Title } from '@mantine/core';
import classes from './WorkoutCard.module.css';

export function WorkoutCard() {
	return (
		<Container size="md">
			<Card className={classes.card} shadow="sm" p="lg" radius="md" withBorder>
				<Grid>
					<Grid.Col visibleFrom="xs" span={5}>
						<Image src="https://placehold.co/400x600?text=video" alt="Workout" />
					</Grid.Col>
					<Grid.Col span={{ base: 12, xs: 7 }} p="md">
						<Title order={2} mb="md">
							23 июля
						</Title>

						<Title order={3}>Разминка</Title>
						<Text>
							<b>3 раунда</b>
							<List withPadding>
								<List.Item>5 снежных ангелов</List.Item>
								<List.Item>10 качающихся планок на локтях</List.Item>
								<List.Item>10-15 приседания широкая постановка</List.Item>
							</List>
						</Text>

						<Divider my="md" />

						<Title order={3}>Комплекс</Title>
						<Text>
							<b>По минутки 10 мин</b>
							<List withPadding>
								<List.Item>20 сек макс повт берпи / 40 сек отжимания с колен</List.Item>
								<List.Item>15+15 прыжки на одной / оставшееся время подъем на носки</List.Item>
							</List>
						</Text>
					</Grid.Col>
				</Grid>
			</Card>
		</Container>
	);
}
