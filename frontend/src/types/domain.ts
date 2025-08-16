export interface Exercise {
	Slug: string;
	Name: string;
	Description: string;
}

export interface WorkoutSection {
	Title: string;
	Protocol: {
		Title: string;
		Description: string;
	};
	Exercises: Array<{
		ExerciseSlug: string;
		Description: string;
	}>;
}

export interface Workout {
	ID: number;
	Date: string;

	Sections: Array<WorkoutSection>;
}
