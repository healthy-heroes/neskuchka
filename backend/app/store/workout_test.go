package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractSlugsFromWorkouts(t *testing.T) {
	tests := []struct {
		name     string
		workouts []*Workout
		expected []ExerciseSlug
	}{
		{
			name:     "Empty workouts slice",
			workouts: []*Workout{},
			expected: []ExerciseSlug{},
		},
		{
			name: "Single workout with no exercises",
			workouts: []*Workout{
				{
					Sections: []WorkoutSection{},
				},
			},
			expected: []ExerciseSlug{},
		},
		{
			name: "Single workout with one exercise",
			workouts: []*Workout{
				{
					Sections: []WorkoutSection{
						{
							Title: "Section 1",
							Exercises: []WorkoutExercise{
								{
									ExerciseSlug: "squat",
								},
							},
						},
					},
				},
			},
			expected: []ExerciseSlug{"squat"},
		},
		{
			name: "Multiple workouts with multiple exercises",
			workouts: []*Workout{
				{
					Sections: []WorkoutSection{
						{
							Title: "Section 1",
							Exercises: []WorkoutExercise{
								{
									ExerciseSlug: "squat",
								},
							},
						},
					},
				},
				{
					Sections: []WorkoutSection{
						{
							Title: "Section 1",
							Exercises: []WorkoutExercise{
								{
									ExerciseSlug: "bench-press",
								},
							},
						},
					},
				},
				{
					Sections: []WorkoutSection{
						{
							Title: "Section 1",
							Exercises: []WorkoutExercise{
								{
									ExerciseSlug: "deadlift",
								},
							},
						},
					},
				},
			},
			expected: []ExerciseSlug{"squat", "bench-press", "deadlift"},
		},
		{
			name: "Multiple workouts with duplicate exercises",
			workouts: []*Workout{
				{
					Sections: []WorkoutSection{
						{
							Title: "Section 1",
							Exercises: []WorkoutExercise{
								{
									ExerciseSlug: "squat",
								},
							},
						},
					},
				},
				{
					Sections: []WorkoutSection{
						{
							Title: "Section 1",
							Exercises: []WorkoutExercise{
								{
									ExerciseSlug: "squat", // Duplicate
								},
							},
						},
					},
				},
				{
					Sections: []WorkoutSection{
						{
							Title: "Section 1",
							Exercises: []WorkoutExercise{
								{
									ExerciseSlug: "bench-press",
								},
							},
						},
					},
				},
			},
			expected: []ExerciseSlug{"squat", "bench-press"},
		},
		{
			name: "Multiple workouts with exercises in multiple sections",
			workouts: []*Workout{
				{
					Sections: []WorkoutSection{
						{
							Title: "Section 1",
							Exercises: []WorkoutExercise{
								{
									ExerciseSlug: "squat",
								},
							},
						},
						{
							Title: "Section 2",
							Exercises: []WorkoutExercise{
								{
									ExerciseSlug: "deadlift",
								},
							},
						},
					},
				},
				{
					Sections: []WorkoutSection{
						{
							Title: "Section 1",
							Exercises: []WorkoutExercise{
								{
									ExerciseSlug: "bench-press",
								},
							},
						},
					},
				},
			},
			expected: []ExerciseSlug{"squat", "deadlift", "bench-press"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slugs := ExtractSlugsFromWorkouts(tt.workouts)
			assert.ElementsMatch(t, tt.expected, slugs)
		})
	}
}

func TestClearUnknownExercisesSlugs(t *testing.T) {
	tests := []struct {
		name       string
		sections   []WorkoutSection
		knownSlugs map[ExerciseSlug]bool
		expected   []WorkoutSection
	}{
		{
			name: "All exercises known",
			sections: []WorkoutSection{
				{
					Title: "Section 1",
					Exercises: []WorkoutExercise{
						{ExerciseSlug: "squat"},
						{ExerciseSlug: "bench-press"},
					},
				},
				{
					Title: "Section 2",
					Exercises: []WorkoutExercise{
						{ExerciseSlug: "deadlift"},
					},
				},
			},
			knownSlugs: map[ExerciseSlug]bool{
				"squat":       true,
				"bench-press": true,
				"deadlift":    true,
			},
			expected: []WorkoutSection{
				{
					Title: "Section 1",
					Exercises: []WorkoutExercise{
						{ExerciseSlug: "squat"},
						{ExerciseSlug: "bench-press"},
					},
				},
				{
					Title: "Section 2",
					Exercises: []WorkoutExercise{
						{ExerciseSlug: "deadlift"},
					},
				},
			},
		},
		{
			name: "Some exercises unknown",
			sections: []WorkoutSection{
				{
					Title: "Section 1",
					Exercises: []WorkoutExercise{
						{ExerciseSlug: "squat"},
						{ExerciseSlug: "unknown-ex"},
					},
				},
				{
					Title: "Section 2",
					Exercises: []WorkoutExercise{
						{ExerciseSlug: "bench-press"},
					},
				},
			},
			knownSlugs: map[ExerciseSlug]bool{
				"squat":       true,
				"bench-press": true,
			},
			expected: []WorkoutSection{
				{
					Title: "Section 1",
					Exercises: []WorkoutExercise{
						{ExerciseSlug: "squat"},
						{ExerciseSlug: ""},
					},
				},
				{
					Title: "Section 2",
					Exercises: []WorkoutExercise{
						{ExerciseSlug: "bench-press"},
					},
				},
			},
		},
		{
			name: "All exercises unknown",
			sections: []WorkoutSection{
				{
					Title: "Section 1",
					Exercises: []WorkoutExercise{
						{ExerciseSlug: "foo"},
						{ExerciseSlug: "bar"},
					},
				},
			},
			knownSlugs: map[ExerciseSlug]bool{
				"baz": true,
			},
			expected: []WorkoutSection{
				{
					Title: "Section 1",
					Exercises: []WorkoutExercise{
						{ExerciseSlug: ""},
						{ExerciseSlug: ""},
					},
				},
			},
		},
		{
			name: "No exercises in sections",
			sections: []WorkoutSection{
				{
					Title:     "Section 1",
					Exercises: []WorkoutExercise{},
				},
				{
					Title:     "Section 2",
					Exercises: []WorkoutExercise{},
				},
			},
			knownSlugs: map[ExerciseSlug]bool{
				"squat": true,
			},
			expected: []WorkoutSection{
				{
					Title:     "Section 1",
					Exercises: []WorkoutExercise{},
				},
				{
					Title:     "Section 2",
					Exercises: []WorkoutExercise{},
				},
			},
		},
		{
			name:     "Empty sections slice",
			sections: []WorkoutSection{},
			knownSlugs: map[ExerciseSlug]bool{
				"squat": true,
			},
			expected: []WorkoutSection{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workout := Workout{
				Sections: tt.sections,
			}

			workout.ClearUnknownExercisesSlugs(tt.knownSlugs)
			assert.Equal(t, tt.expected, workout.Sections)
		})
	}
}
