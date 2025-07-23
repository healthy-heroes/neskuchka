import { IconLockOpen2, IconMoodSmileBeam, IconSquareRoundedCheck } from '@tabler/icons-react';
import { Card, Container, SimpleGrid, Text, useMantineTheme } from '@mantine/core';
import classes from './Features.module.css';

const featuresData = [
	{
		title: 'Бесплатно',
		description: 'Доступ к бесплатным тренировкам – занимайтесь в удобное время, в любом месте!',
		icon: IconLockOpen2,
	},
	{
		title: 'Просто и эффективно',
		description:
			'Простые и понятные упражнения, подходящие для любого уровня подготовки. Помогут улучшить гибкость, силу и выносливость.',
		icon: IconMoodSmileBeam,
	},
	{
		title: 'Удобно',
		description:
			'Занимайтесь дома, на работе, на отдыхе – где вам удобно! Все, что нужно – это телефон или планшет.',
		icon: IconSquareRoundedCheck,
	},
];

export function Features() {
	const theme = useMantineTheme();
	const features = featuresData.map((feature) => (
		<Card key={feature.title} shadow="md" radius="md" className={classes.card} padding="xl">
			<feature.icon size={50} stroke={1.5} color={theme.colors.blue[6]} />
			<Text fz="lg" fw={500} className={classes.cardTitle} mt="md">
				{feature.title}
			</Text>
			<Text fz="sm" c="dimmed" mt="sm">
				{feature.description}
			</Text>
		</Card>
	));

	return (
		<Container size="lg" py="xl">
			<SimpleGrid cols={{ base: 1, md: 3 }} spacing="xl">
				{features}
			</SimpleGrid>
		</Container>
	);
}
