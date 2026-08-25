import { Button, Group, Text } from '@mantine/core';
import { RouteLink } from '@/components/RouteLink/RouteLink';
import { Workout } from '@/types/domain';
import { formatIsoDateShort, formatWeekday, isToday } from '@/utils/dates';
import classes from './WorkoutRows.module.css';

export interface WorkoutRowsProps {
	/** Тренировки трека, свежие первыми — вместе с неопубликованными. */
	workouts: Array<Workout>;

	onDelete: (workout: Workout) => void;
}

/**
 * WorkoutRows — весь трек списком, каким его правит владелец.
 *
 * От истории на публичной странице отличается тем, что показывает и то, чего
 * участники ещё не видят, и тем, что строка тут не ссылка целиком: в ней два
 * действия, и промахнуться мимо них по всей ширине нельзя.
 */
export function WorkoutRows({ workouts, onDelete }: WorkoutRowsProps) {
	return (
		<div className={classes.rows}>
			{workouts.map((workout) => (
				<WorkoutRow key={workout.ID} workout={workout} onDelete={onDelete} />
			))}
		</div>
	);
}

function WorkoutRow({
	workout,
	onDelete,
}: {
	workout: Workout;
	onDelete: WorkoutRowsProps['onDelete'];
}) {
	const today = isToday(workout.Date);
	const published = workout.IsPublished ?? true;

	// Восемь одинаковых «Удалить» подряд: без даты в имени со скринридера
	// не понять, какую тренировку удаляешь
	const date = formatIsoDateShort(workout.Date);

	return (
		<div className={classes.row}>
			<div>
				<Text
					ff="heading"
					fz="xl"
					fw={published ? 600 : 500}
					tt="uppercase"
					lts="0.02em"
					c={dateColor(today, published)}
				>
					{date}
				</Text>
				<Text fz="xs" c="gray.7">
					{formatWeekday(workout.Date)}
				</Text>
			</div>

			<Text fz="sm" c="gray.8">
				{workout.Sections.map((section, index) => (
					<span key={`${section.Title}-${index}`}>
						{index > 0 && <span className={classes.summarySep}>·</span>}
						{[section.Title, section.Protocol.Title].filter(Boolean).join(' · ')}
					</span>
				))}
			</Text>

			<Text fz="sm" c="gray.7">
				{published ? 'Опубликована' : 'Не опубликована'}
			</Text>

			{workout.IsEditable && (
				<Group gap="xs" justify="flex-end" wrap="nowrap">
					{/* renderRoot, а не component: через полиморфный проп Button теряет
					    типы роутера, и params перестают проверяться по маршруту */}
					{/* flex="none": обрезанная подпись выглядит как опечатка, а вылезший
					    за край блок сразу виден и чинится */}
					<Button
						renderRoot={(props) => (
							<RouteLink
								to="/workouts/$workoutId/edit"
								params={{ workoutId: workout.ID }}
								underline="never"
								{...props}
							/>
						)}
						variant="default"
						size="xs"
						flex="none"
						aria-label={`Изменить тренировку ${date}`}
					>
						Изменить
					</Button>
					<Button
						variant="default"
						size="xs"
						flex="none"
						c="copper.7"
						aria-label={`Удалить тренировку ${date}`}
						onClick={() => onDelete(workout)}
					>
						Удалить
					</Button>
				</Group>
			)}
		</div>
	);
}

/** Сегодняшняя выделена цветом, ещё не опубликованная — приглушена. */
function dateColor(today: boolean, published: boolean): string | undefined {
	if (today) {
		return 'copper.6';
	}

	return published ? undefined : 'gray.6';
}
