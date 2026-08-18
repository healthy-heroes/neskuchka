export interface User {
	ID: string;
	Name: string;
	Avatar?: string;
}

export interface UserSettings {
	Name: string;
	Email: string;
}

export interface Track {
	ID: string;
	Name: string;
	Description: string;

	Author: TrackAuthor;
}

export interface TrackAuthor {
	ID: string;
	Name: string;
}

export interface Exercise {
	Slug: string;
	Name: string;
	Description: string;
}

export interface WorkoutExercise {
	ExerciseSlug: string;

	/** Что делать: «10», «3х2 @ 80%». Несколько строк — несколько подходов. */
	Prescription: string[];
	/** Название упражнения. Временно лежит на тренировке, пока нет справочника упражнений. */
	Name: string;
}

export interface WorkoutSection {
	Title: string;
	Protocol: {
		Title: string;
		Description: string;
	};
	Exercises: Array<WorkoutExercise>;
}

export interface Workout {
	ID: string;
	TrackID: string;
	Date: string;

	Sections: Array<WorkoutSection>;

	Notes?: string;
}
