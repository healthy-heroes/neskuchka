import dayjs from 'dayjs';
import { randomId } from '@mantine/hooks';
import { Workout, WorkoutExercise, WorkoutSection } from '@/types/domain';

const idPrefix = 'new';

export type WorkoutFormData = Workout & {
	Sections: Array<WorkoutSectionFormData>;
};

export type WorkoutSectionFormData = WorkoutSection & {
	_key: string;

	Exercises: Array<WorkoutExerciseFormData>;
};

export type WorkoutExerciseFormData = WorkoutExercise & {
	_key: string;
};

export function convertToFormData(data: Workout): WorkoutFormData {
	return {
		...data,
		Sections: data.Sections.map((section) => ({
			...section,
			_key: randomId(idPrefix),
			Exercises: section.Exercises.map((exercise) => ({
				...exercise,
				_key: randomId(idPrefix),
			})),
		})),
	};
}

export function convertToDomainData(data: WorkoutFormData): Workout {
	return {
		...data,
		Sections: data.Sections.map((section) => ({
			...section,
			Exercises: section.Exercises.map((exercise) => ({
				...exercise,
			})),
		})),
	};
}

// Helpers for creating initial values
export function makeInitialValues(): WorkoutFormData {
	return {
		ID: randomId(idPrefix),
		Date: dayjs().format('YYYY-MM-DD'),
		Sections: [makeSection('Разминка'), makeSection('Комплекс')],
	};
}

export function makeSection(title: string = 'Комплекс'): WorkoutSectionFormData {
	return {
		_key: randomId('new'),
		Title: title,
		Protocol: {
			Title: '',
			Description: '',
		},
		Exercises: [makeExercise()],
	};
}

export function makeExercise(): WorkoutExerciseFormData {
	return {
		_key: randomId('new'),
		ExerciseSlug: '',
		Description: '',
	};
}
