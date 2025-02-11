from sqlmodel import Session

from app.domain.workout import (
    WrokoutCriteria,
)
from app.infrastructure.db.workout import WorkoutDbRepository
from app.tests.fixtures.domain import create_track, create_workout


def test_add_get_workout(session: Session):
    workout = create_workout()

    workout_repository = WorkoutDbRepository(session)
    workout_repository.add(workout)

    workout_from_db = workout_repository.get_by_id(workout.id)

    assert workout_from_db == workout


def test_get_list_workout(session: Session):
    workout_repository = WorkoutDbRepository(session)

    track = create_track()
    for _ in range(10):
        workout = create_workout(track_id=track.id)
        workout_repository.add(workout)

    workouts_from_db = workout_repository.get_list(
        track_id=track.id, criteria=WrokoutCriteria(limit=5)
    )

    assert len(workouts_from_db) == 5
