import { Box, Divider, Group, Stack, Text, useMantineTheme } from '@mantine/core';

const STEPS = ['xs', 'sm', 'md', 'lg', 'xl'] as const;

/**
 * SizeScale — ступени отступов и радиусов.
 *
 * Обе шкалы именованные: в разметке пишут `p="lg"` и `radius="xl"`, а не пиксели.
 * Карточка нужна, чтобы было видно, какому шагу какая величина соответствует.
 */
export function SizeScale() {
	const theme = useMantineTheme();

	return (
		<Stack gap="xl">
			<Stack gap="xs">
				<Heading text="Отступы" hint='p="lg", gap="xs", var(--mantine-spacing-lg)' />
				{STEPS.map((step) => (
					<Group key={step} align="center" gap="md" wrap="nowrap">
						<Text fz="xs" c="gray.7" ff="monospace" w={60}>
							{step}
						</Text>
						<Box h={20} bdrs="xs" bg="copper.4" w={`var(--mantine-spacing-${step})`} />
						<Text fz="xs" c="gray.7" ff="monospace">
							{toPx(theme.spacing[step])}
						</Text>
					</Group>
				))}
			</Stack>

			<Divider />

			<Stack gap="xs">
				<Heading text="Радиусы" hint='radius="lg", var(--mantine-radius-lg)' />
				<Group gap="md" align="flex-end">
					{STEPS.map((step) => (
						<Stack key={step} gap={4} align="center">
							<Box
								w={72}
								h={72}
								bdrs={step}
								bg="slate.2"
								bd="1px solid var(--mantine-color-slate-4)"
							/>
							<Text fz="xs" fw={600}>
								{step}
							</Text>
							<Text fz="xs" c="gray.7" ff="monospace">
								{toPx(theme.radius[step])}
							</Text>
						</Stack>
					))}
				</Group>
			</Stack>
		</Stack>
	);
}

function Heading({ text, hint }: { text: string; hint: string }) {
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

/**
 * `calc(0.625rem * var(--mantine-scale))` → `10px`.
 *
 * Тема хранит ступени в rem и заворачивает их в calc с масштабом, а на карточке
 * нужна величина, которую можно сверить с макетом.
 */
function toPx(value: string) {
	const rem = value.match(/([\d.]+)rem/);

	return rem ? `${Math.round(parseFloat(rem[1]) * 16)}px` : value;
}
