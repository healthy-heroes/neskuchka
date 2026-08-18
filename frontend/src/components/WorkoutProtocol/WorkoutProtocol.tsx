import { useState } from 'react';
import { WorkoutSection } from '@/types/domain';
import { ExerciseRow } from '../ExerciseRow/ExerciseRow';
import classes from './WorkoutProtocol.module.css';

export interface WorkoutProtocolProps {
	sections: Array<WorkoutSection>;
}

/**
 * WorkoutProtocol — секции тренировки с упражнениями.
 *
 * Грид объявлен здесь, а секции и списки внутри — display: contents, поэтому
 * колонка предписаний общая на всю тренировку: «10» в разминке стоит на той же
 * линии, что «3х2 @ 80%» в комплексе, и жёлоб не прыгает между секциями.
 * Отдельный грид на секцию всё это ломает.
 *
 * Карточку вокруг рисует вызывающий: на странице тренировки и на треке она разная.
 */
export function WorkoutProtocol({ sections }: WorkoutProtocolProps) {
	const [opened, setOpened] = useState<Record<string, boolean>>({});

	function toggle(key: string) {
		setOpened((state) => ({ ...state, [key]: !state[key] }));
	}

	return (
		<div className={classes.protocol}>
			{sections.map((section, sectionIndex) => (
				<div key={`${section.Title}-${sectionIndex}`} className={classes.section}>
					{sectionIndex > 0 && <div className={classes.divider} />}

					<div className={classes.sectionHead}>
						<h2 className={classes.sectionTitle}>{section.Title}</h2>
						{section.Protocol.Title && (
							<span className={classes.protocolChip}>{section.Protocol.Title}</span>
						)}
					</div>

					{section.Protocol.Description && (
						<p className={classes.sectionNote}>{section.Protocol.Description}</p>
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
