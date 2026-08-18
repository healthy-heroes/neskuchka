import type { ReactNode } from 'react';
import { IconChevronDown } from '@tabler/icons-react';
import clsx from 'clsx';
import { Collapse } from '@mantine/core';
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
 * чтобы колонка предписаний была общей на всю тренировку, а не на секцию.
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
	const cellClass = clsx(classes.cell, expandable && classes.clickable);
	const onClick = expandable ? onToggle : undefined;
	const mutedProps = expandable ? ({ tabIndex: -1, 'aria-hidden': true } as const) : {};

	return (
		<li className={classes.row}>
			<button
				type="button"
				className={clsx(cellClass, classes.prescription)}
				onClick={onClick}
				{...mutedProps}
			>
				{prescription.map((line) => (
					<span key={line} className={classes.prescriptionLine}>
						{line}
					</span>
				))}
			</button>

			<button
				type="button"
				className={clsx(cellClass, classes.name)}
				aria-expanded={expandable ? opened : undefined}
				onClick={onClick}
			>
				{name}
			</button>

			<button
				type="button"
				className={clsx(cellClass, classes.hint)}
				onClick={onClick}
				{...mutedProps}
			>
				{expandable && (
					<>
						<span>{opened ? 'свернуть' : 'как делать'}</span>
						<IconChevronDown size={16} className={clsx(opened && classes.chevronUp)} />
					</>
				)}
			</button>

			{expandable && (
				<Collapse expanded={opened} className={classes.panel}>
					<div className={classes.panelInner}>{content}</div>
				</Collapse>
			)}
		</li>
	);
}
