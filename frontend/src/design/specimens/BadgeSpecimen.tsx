import { Badge, Group, Stack, Text } from '@mantine/core';

const VARIANTS = ['filled', 'light', 'outline', 'dot'] as const;

/**
 * BadgeSpecimen — бейджи в теме.
 *
 * Пилюльный радиус и трекинг лейбла заданы в `Badge.extend`; в разметке
 * остаётся только цвет и вариант.
 */
export function BadgeSpecimen() {
	return (
		<Stack gap="xl">
			<Stack gap="xs">
				<Label text="Варианты" />
				<Group gap="md">
					{VARIANTS.map((variant) => (
						<Badge key={variant} variant={variant}>
							3 раунда
						</Badge>
					))}
				</Group>
			</Stack>

			<Stack gap="xs">
				<Label text="Цвета" />
				<Group gap="md">
					<Badge color="copper">Сегодня</Badge>
					<Badge color="slate">Выполнено</Badge>
					<Badge color="gray" variant="light">
						Черновик
					</Badge>
				</Group>
			</Stack>
		</Stack>
	);
}

function Label({ text }: { text: string }) {
	return (
		<Text fz="sm" fw={600} lts="0.06em" tt="uppercase">
			{text}
		</Text>
	);
}
