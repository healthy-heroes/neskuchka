import { useState } from 'react';
import { Group, Text, Title } from '@mantine/core';
import { WorkoutSection } from '@/types/domain';
import { ExerciseRow } from '../ExerciseRow/ExerciseRow';
import classes from './WorkoutSections.module.css';

export interface WorkoutSectionsProps {
	sections: Array<WorkoutSection>;
}

/**
 * WorkoutSections — секции тренировки с упражнениями.
 *
 * Не «протокол»: в домене протокол — это поле секции, схема, по которой она идёт
 * («3 раунда», «21-15-9 на время»). Здесь же весь список секций.
 *
 * Грид объявлен здесь, а секции и списки внутри — display: contents, поэтому
 * колонка предписаний общая на всю тренировку: «10» в разминке стоит на той же
 * линии, что «3х2 @ 80%» в комплексе, и жёлоб не прыгает между секциями.
 * Отдельный грид на секцию всё это ломает.
 *
 * Карточку вокруг рисует вызывающий: на странице тренировки и на треке она разная.
 */
export function WorkoutSections({ sections }: WorkoutSectionsProps) {
	const [opened, setOpened] = useState<Record<string, boolean>>({});

	function toggle(key: string) {
		setOpened((state) => ({ ...state, [key]: !state[key] }));
	}

	return (
		<div className={classes.sections}>
			{sections.map((section, sectionIndex) => (
				<div key={`${section.Title}-${sectionIndex}`} className={classes.section}>
					{sectionIndex > 0 && <div className={classes.divider} />}

					<Group className={classes.fullRow} align="baseline" gap="sm">
						<Title order={2} fz="xl" lts="0.06em">
							{section.Title}
						</Title>
						{section.Protocol.Title && (
							<Text span className={classes.chip} fz="xs" fw={600} c="slate.7" bg="slate.1">
								{section.Protocol.Title}
							</Text>
						)}
					</Group>

					{section.Protocol.Description && (
						<Text
							className={classes.fullRow}
							mt="xs"
							maw="68ch"
							fz="sm"
							lh="md"
							c="gray.8"
							textWrap="pretty"
						>
							{section.Protocol.Description}
						</Text>
					)}

					<ul className={classes.exercises}>
						{section.Exercises.map((exercise, exerciseIndex) => {
							const key = `${sectionIndex}:${exerciseIndex}`;

							return (
								<ExerciseRow
									key={key}
									prescription={exercise.Prescription}
									name={exercise.Name}
									opened={opened[key]}
									onToggle={() => toggle(key)}
								/>
							);
						})}
					</ul>
				</div>
			))}
		</div>
	);
}
