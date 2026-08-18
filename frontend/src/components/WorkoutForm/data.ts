import dayjs from 'dayjs';
import { randomId } from '@mantine/hooks';
import { Workout, WorkoutExercise, WorkoutSection } from '@/types/domain';

const idPrefix = 'new';

export type WorkoutFormData = Omit<Workout, 'Sections'> & {
	Sections: Array<WorkoutSectionFormData>;
};

export type WorkoutSectionFormData = Omit<WorkoutSection, 'Exercises'> & {
	_key: string;

	Exercises: Array<WorkoutExerciseFormData>;
};

export type WorkoutExerciseFormData = Omit<WorkoutExercise, 'Prescription'> & {
	_key: string;

	/** В форме предписания редактируются одним полем, по подходу на строку. */
	Prescription: string;
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
				Prescription: exercise.Prescription.join('\n'),
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
				Prescription: splitPrescription(exercise.Prescription),
			})),
		})),
	};
}

/** Пустые строки отбрасываем: упражнение без предписания — нормальный случай. */
function splitPrescription(value: string): string[] {
	return value
		.split('\n')
		.map((line) => line.trim())
		.filter((line) => line !== '');
}

// Helpers for creating initial values
export function makeInitialValues(trackId: string): WorkoutFormData {
	return {
		ID: randomId(idPrefix),
		TrackID: trackId,
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
		Prescription: '',
		Name: '',
	};
}
