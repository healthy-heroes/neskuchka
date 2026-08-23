import { useState, type ReactNode } from 'react';
import { Text } from '@mantine/core';
import { StoryPreview } from '../StoryBook/StoryPreview';
import { ExerciseRow } from './ExerciseRow';
import sections from '../WorkoutSections/WorkoutSections.module.css';

export default {
	title: 'ExerciseRow',
};

/**
 * Строка объявлена как `display: contents` и раскладывается по гриду карточки
 * протокола, поэтому стори берёт настоящий грид из WorkoutSections, а не заводит
 * свой: копия грида разъехалась бы с оригиналом на первой же правке.
 */
function Grid({ children }: { children: ReactNode }) {
	return (
		<div className={sections.sections}>
			<ul className={sections.exercises}>{children}</ul>
		</div>
	);
}

export function Default() {
	return (
		<StoryPreview>
			<Grid>
				<ExerciseRow prescription={['10']} name="Отжимание" />
				<ExerciseRow prescription={['3х2 @ 80%']} name="Приседание со штангой" />
				<ExerciseRow prescription={['5', '5', '5']} name="Подтягивание" />
			</Grid>
		</StoryPreview>
	);
}

/** С описанием строка раскрывается; без него — просто строка. */
export function Expandable() {
	const [opened, setOpened] = useState(true);

	return (
		<StoryPreview>
			<Grid>
				<ExerciseRow
					prescription={['10']}
					name="Отжимание"
					opened={opened}
					onToggle={() => setOpened((value) => !value)}
					content={
						<Text fz="sm" c="gray.8" py="xs">
							Корпус в линию, локти вдоль тела. Опускаться до касания грудью пола.
						</Text>
					}
				/>
				<ExerciseRow prescription={['10']} name="Приседание" />
			</Grid>
		</StoryPreview>
	);
}
