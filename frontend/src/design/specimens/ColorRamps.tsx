import { Box, Group, SimpleGrid, Stack, Text, Title, useMantineTheme } from '@mantine/core';

const RAMPS = [
	{ name: 'copper', role: 'действия, акценты, «сегодня»' },
	{ name: 'slate', role: 'тёмные плашки, статус «выполнено»' },
	{ name: 'gray', role: 'тёплые нейтральные вместо серых Mantine' },
] as const;

/**
 * ColorRamps — палитра темы: три рампы по десять ступеней.
 *
 * Ступени берутся из разрешённой темы, а не из исходника theme.ts: так карточка
 * показывает ровно то, во что превращаются токены `copper.6`, `slate.7`, `gray.7`.
 */
export function ColorRamps() {
	const theme = useMantineTheme();

	return (
		<Stack gap="xl">
			{RAMPS.map(({ name, role }) => (
				<Stack key={name} gap="sm">
					<Group align="baseline" gap="sm">
						<Title order={4}>{name}</Title>
						<Text fz="sm" c="gray.7">
							{role}
						</Text>
					</Group>

					<SimpleGrid cols={{ base: 5, sm: 10 }} spacing="xs">
						{theme.colors[name].map((hex, shade) => (
							<Stack key={shade} gap={4} align="center">
								<Box w="100%" h={64} bdrs="sm" bg={`${name}.${shade}`} />
								<Text fz="xs" fw={600}>
									{shade}
								</Text>
								<Text fz="xs" c="gray.7" ff="monospace">
									{hex}
								</Text>
							</Stack>
						))}
					</SimpleGrid>
				</Stack>
			))}
		</Stack>
	);
}
