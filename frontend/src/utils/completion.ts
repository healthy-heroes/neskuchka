import { Workout } from '@/types/domain';

/**
 * ЗАГЛУШКА на первый этап редизайна.
 *
 * Отметок о выполнении в бэкенде ещё нет: ни таблицы workout_completion,
 * ни эндпоинтов. Экраны при этом построены вокруг статуса — полоса прогресса,
 * колонка статуса в истории, состояние карточки «Сегодня», — поэтому статус
 * приходится изображать.
 *
 * Считается детерминированно по ID, чтобы вид не прыгал между рендерами,
 * и даёт примерно три четверти выполненных — как на макете, где есть и то и другое.
 *
 * Когда появится отметка выполнения, этот файл удаляется целиком, а вызовы
 * заменяются полем из API.
 */
export function isWorkoutDone(workout: Workout): boolean {
	const date = new Date(workout.Date);
	const today = new Date();
	today.setHours(0, 0, 0, 0);

	if (date >= today) {
		return false;
	}

	return hashCode(workout.ID) % 4 !== 0;
}

function hashCode(value: string): number {
	let hash = 0;
	for (let i = 0; i < value.length; i += 1) {
		hash = (hash * 31 + value.charCodeAt(i)) | 0;
	}

	return Math.abs(hash);
}
