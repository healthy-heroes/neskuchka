import uuid

import pytest
from sqlmodel import Session

from app.domain.workout import Workout, WorkoutExercise, WorkoutProtocol, WorkoutSection
from app.domain.track import TrackId
from app.domain.exercise import ExerciseSlug
from app.infrastructure.db.workout import WorkoutDbRepository


def test_add_get_workout(session: Session):
    track_id = TrackId(str(uuid.uuid4()))
    exercise_slug = ExerciseSlug("test_exercise")

    workout_exercise = WorkoutExercise(
        exercise_slug=exercise_slug,
        repetitions=3,
        repetitions_text="make 3 reps or more",
    )

    workout_section = WorkoutSection(
        title="Test section",
        protocol=WorkoutProtocol(
            title="Tabata 8x20:10",
            description="Make 8 reps of 20 seconds, rest 10 seconds",
        ),
        exercises=[workout_exercise],
    )

    workout = Workout(
        date="2025-01-30",
        track_id=track_id,
        sections=[workout_section],
    )

    workout_repository = WorkoutDbRepository(session)
    workout_repository.add(workout)

    workout_from_db = workout_repository.get_by_id(workout.id)

    assert workout_from_db == workout


@pytest.mark.skip(reason="Not implemented")
def test_get_list_workout(session: Session):
    # todo data generators
    pass
