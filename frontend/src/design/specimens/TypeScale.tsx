import { Divider, Group, Stack, Text, Title, useMantineTheme } from '@mantine/core';

const HEADINGS = [1, 2, 3, 4, 5, 6] as const;
const SIZES = ['xs', 'sm', 'md', 'lg', 'xl'] as const;

/**
 * TypeScale — типографика темы: заголовки Oswald по `order` и текстовая шкала Commissioner.
 *
 * Uppercase и трекинг у заголовков приходят из `Title.extend` в теме, руками здесь
 * ничего не задаётся — карточка показывает поведение темы, а не своё оформление.
 */
export function TypeScale() {
	const theme = useMantineTheme();

	return (
		<Stack gap="xl">
			<Stack gap="xs">
				<Label text="Заголовки — Oswald" hint={String(theme.headings.fontFamily)} />
				{HEADINGS.map((order) => (
					<Group key={order} align="baseline" gap="md" wrap="nowrap">
						<Text fz="xs" c="gray.7" ff="monospace" w={60}>
							h{order}
						</Text>
						<Title order={order}>Тренировка дня</Title>
					</Group>
				))}
			</Stack>

			<Divider />

			<Stack gap="xs">
				<Label text="Текст — Commissioner" hint={String(theme.fontFamily)} />
				{SIZES.map((size) => (
					<Group key={size} align="baseline" gap="md" wrap="nowrap">
						<Text fz="xs" c="gray.7" ff="monospace" w={60}>
							{size}
						</Text>
						<Text fz={size}>Пять подходов по восемь повторений</Text>
					</Group>
				))}
			</Stack>

			<Divider />

			<Stack gap="xs">
				<Label text="Вторичный текст" hint='c="gray.7", не c="dimmed"' />
				<Text c="gray.7">Разминка · 3 раунда · по минутке</Text>
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
