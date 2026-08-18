import type { ReactNode } from 'react';
import { IconChevronDown } from '@tabler/icons-react';
import clsx from 'clsx';
import { Collapse, Group, Text, UnstyledButton } from '@mantine/core';
import classes from './ExerciseRow.module.css';

export interface ExerciseRowProps {
	prescription: Array<string>;
	name: string;

	/**
	 * Содержимое раскрытой панели — описание, на что смотреть, вариации, видео.
	 * Приедет из справочника упражнений; пока его нет, строка не раскрывается.
	 */
	content?: ReactNode;

	opened?: boolean;
	onToggle?: () => void;
}

/**
 * ExerciseRow — строка упражнения: предписание, название, подсказка.
 *
 * Своего грида не заводит: раскладывается по трём колонкам родительской карточки,
 * чтобы колонка предписаний была общей на всю тренировку, а не на секцию. По той же
 * причине <li> остаётся сырым: List.Item добавляет две обёртки, и они схлопнут грид.
 */
export function ExerciseRow({
	prescription,
	name,
	content,
	opened = false,
	onToggle,
}: ExerciseRowProps) {
	const expandable = Boolean(content);

	// Три ячейки кликаются одинаково, но фокус и озвучка — только на названии,
	// иначе на каждую строку приходится три одинаковых элемента управления.
	const cellClass = clsx(classes.cell, !expandable && classes.static);
	const onClick = expandable ? onToggle : undefined;
	const mutedProps = expandable ? ({ tabIndex: -1, 'aria-hidden': true } as const) : {};

	return (
		<li className={classes.row}>
			<UnstyledButton
				className={clsx(cellClass, classes.prescription)}
				onClick={onClick}
				{...mutedProps}
			>
				{prescription.map((line) => (
					<Text key={line} span ff="heading" fz="md" fw={600} lts="0.02em" c="copper.6">
						{line}
					</Text>
				))}
			</UnstyledButton>

			<UnstyledButton
				className={cellClass}
				fz="md"
				lh="sm"
				aria-expanded={expandable ? opened : undefined}
				onClick={onClick}
			>
				{name}
			</UnstyledButton>

			<UnstyledButton className={cellClass} onClick={onClick} {...mutedProps}>
				{expandable && (
					<Group gap="xs" wrap="nowrap">
						<Text span fz="xs" fw={600} c="gray.7">
							{opened ? 'свернуть' : 'как делать'}
						</Text>
						<IconChevronDown size={16} className={clsx(opened && classes.chevronUp)} />
					</Group>
				)}
			</UnstyledButton>

			{expandable && (
				<Collapse expanded={opened} className={classes.panel}>
					<Group align="flex-start" gap="lg" pt="xs" pb="lg">
						{content}
					</Group>
				</Collapse>
			)}
		</li>
	);
}
