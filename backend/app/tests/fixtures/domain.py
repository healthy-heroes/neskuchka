from datetime import date
import uuid

from app.domain.exercise import (
    Exercise,
    ExerciseCriteria,
    ExerciseRepository,
    ExerciseSlug,
)
from app.domain.track import Track, TrackId, TrackRepository
from app.domain.user import User
from app.domain.workout import (
    Workout,
    WorkoutExercise,
    WorkoutProtocol,
    WorkoutSection,
    WorkoutRepository,
    WorkoutId,
    WrokoutCriteria,
)


# Exercise
def create_exercise(slug: ExerciseSlug | None = None) -> Exercise:
    return Exercise(
        slug=slug or ExerciseSlug(str(uuid.uuid4())),
        name="Some exercise",
        description="Some exercise description",
    )


class ExerciseRepositoryTest(ExerciseRepository):
    def __init__(self):
        self.exercises = {}

    def add(self, exercise: Exercise) -> Exercise:
        self.exercises[exercise.slug] = exercise
        return exercise

    def get_all(self) -> list[Exercise]:
        return list(self.exercises.values())

    def get_by_slug(self, slug: ExerciseSlug) -> Exercise | None:
        return self.exercises.get(slug)

    def find(self, criteria: ExerciseCriteria) -> list[Exercise]:
        return [self.exercises[slug] for slug in criteria.slugs]


# User
def create_user() -> User:
    login = "test-user"
    return User(
        id=User.create_id(login),
        login=login,
        name="Test user",
        email="test-user@example.com",
    )


# Track
def create_track() -> Track:
    return Track(
        name="Main test track",
        owner_id=create_user().id,
    )


class TrackRepositoryTest(TrackRepository):
    def __init__(self):
        self.tracks = {}

    def add(self, track: Track) -> Track:
        self.tracks[track.id] = track
        return track

    def get_by_id(self, track_id: TrackId) -> Track | None:
        return self.tracks.get(track_id)

    def get_main_track(self) -> Track | None:
        return list(self.tracks.values())[0] if self.tracks else None


# Workout
def create_workout_exercise() -> WorkoutExercise:
    return WorkoutExercise(
        exercise_slug=create_exercise().slug,
        repetitions=10,
        repetitions_text="10 reps",
        weight=10,
        weight_text="10 kg",
    )


def create_workout_section(
    exercises: list[WorkoutExercise] | None = None,
) -> WorkoutSection:
    return WorkoutSection(
        title="Test section",
        protocol=WorkoutProtocol(
            title="Test protocol",
        ),
        exercises=exercises or [create_workout_exercise()],
    )


def create_workout(
    track_id: TrackId = None, sections: list[WorkoutSection] | None = None
) -> Workout:
    return Workout(
        date=date.today(),
        track_id=track_id or create_track().id,
        sections=sections or [create_workout_section()],
    )


class WorkoutRepositoryTest(WorkoutRepository):
    def __init__(self):
        self.workouts = {}

    def add(self, workout: Workout) -> Workout:
        self.workouts[workout.id] = workout
        return workout

    def get_by_id(self, workout_id: WorkoutId) -> Workout | None:
        return self.workouts.get(workout_id)

    def get_list(self, track_id: TrackId, criteria: WrokoutCriteria) -> list[Workout]:
        return list(self.workouts.values())[: criteria.limit]
