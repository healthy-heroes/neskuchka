import { createTrackWorkouts } from '@/api/fixtures/workout';
import { StoryPreview } from '@/components/StoryBook/StoryPreview';
import { WorkoutRows } from './WorkoutRows';

export default {
	title: 'WorkoutRows',
};

const workouts = createTrackWorkouts({ count: 6, planned: 2 });

/**
 * Весь трек глазами владельца: сверху две запланированные, ниже опубликованные.
 *
 * Здесь видны все три состояния строки сразу — приглушённая неопубликованная,
 * выделенная сегодняшняя и обычная прошедшая, — и то, что кнопки остаются
 * только у тех, чьё окно правки ещё открыто.
 */
export function Default() {
	return (
		<StoryPreview isPage>
			<WorkoutRows workouts={workouts} onDelete={() => {}} />
		</StoryPreview>
	);
}
