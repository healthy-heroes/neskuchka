import { Card, Group, Paper, Stack, Text, Title } from '@mantine/core';

/**
 * SurfaceSpecimen — плоскости, из которых собраны экраны.
 *
 * Card берёт радиус и рамку из темы, поэтому пропов оформления у него нет.
 * Paper показан в двух ролях, которые реально встречаются: светлая плашка
 * на тёплом нейтральном фоне и тёмная плашка на slate.
 */
export function SurfaceSpecimen() {
	return (
		<Stack gap="xl">
			<Stack gap="xs">
				<Label text="Card — умолчания темы" hint="radius lg, withBorder" />
				<Card p="lg">
					<Title order={4}>Комплекс</Title>
					<Text c="gray.7" mt="xs">
						По минутке 10 мин · 20 сек берпи / 40 сек отжимания с колен
					</Text>
				</Card>
			</Stack>

			<Stack gap="xs">
				<Label text="Paper — светлая плашка" hint='withBorder radius="lg" bg="gray.0"' />
				<Paper withBorder radius="lg" p="md" bg="gray.0">
					<Group justify="space-between" align="baseline">
						<Text fz="xs" fw={600} lts="0.06em" tt="uppercase" c="gray.7">
							Последние 30 дней
						</Text>
						<Text fz="sm" fw={600}>
							12 из 30
						</Text>
					</Group>
				</Paper>
			</Stack>

			<Stack gap="xs">
				<Label text="Тёмная шапка" hint='bg="slate.7"' />
				<Paper radius="lg" p="md" bg="slate.7">
					<Group align="baseline" gap="md">
						<Text span ff="heading" fz="sm" fw={500} lts="0.14em" tt="uppercase" c="copper.3">
							Сегодня
						</Text>
						<Title order={3} c="white">
							Среда, 12 марта
						</Title>
					</Group>
				</Paper>
			</Stack>
		</Stack>
	);
}

function Label({ text, hint }: { text: string; hint: string }) {
	return (
		<Group align="baseline" gap="sm">
			<Text fz="sm" fw={600} lts="0.06em" tt="uppercase">
				{text}
			</Text>
			<Text fz="xs" c="gray.7" ff="monospace">
				{hint}
			</Text>
		</Group>
	);
}
