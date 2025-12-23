export interface User {
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
	Description: string;
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
	Date: string;

	Sections: Array<WorkoutSection>;

	Notes?: string;
}
