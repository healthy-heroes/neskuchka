import type { ReactNode } from 'react';
import { IconPlayerPlayFilled } from '@tabler/icons-react';
import { Button, Group, Stack, Text } from '@mantine/core';

const VARIANTS = ['filled', 'light', 'outline', 'subtle', 'default'] as const;
const SIZES = ['xs', 'sm', 'md', 'lg', 'xl'] as const;

/**
 * ButtonSpecimen — кнопки в теме: варианты, ступени размера и состояния.
 *
 * Радиус и вес шрифта приходят из `Button.extend`, цвет — из `primaryColor`,
 * поэтому пропов оформления здесь нет: карточка показывает умолчания темы.
 */
export function ButtonSpecimen() {
	return (
		<Stack gap="xl">
			<Row label="Варианты">
				{VARIANTS.map((variant) => (
					<Button key={variant} variant={variant}>
						Начать
					</Button>
				))}
			</Row>

			<Row label="Размеры">
				{SIZES.map((size) => (
					<Button key={size} size={size}>
						Начать
					</Button>
				))}
			</Row>

			<Row label="Состояния">
				<Button leftSection={<IconPlayerPlayFilled size={18} />}>С иконкой</Button>
				<Button loading>Загрузка</Button>
				<Button disabled>Недоступна</Button>
			</Row>

			<Row label="Вторичный цвет">
				<Button color="slate">Отметить выполненной</Button>
				<Button color="slate" variant="light">
					Отметить выполненной
				</Button>
			</Row>
		</Stack>
	);
}

function Row({ label, children }: { label: string; children: ReactNode }) {
	return (
		<Stack gap="xs">
			<Text fz="sm" fw={600} lts="0.06em" tt="uppercase">
				{label}
			</Text>
			<Group gap="md" align="center">
				{children}
			</Group>
		</Stack>
	);
}
