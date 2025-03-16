from datetime import date
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
    count = 10
    for i in range(count):
        workout = create_workout(
            date=date(2025, 1, i + 1),
            track=track,
        )
        workout_repository.add(workout)

    workouts_from_db = workout_repository.get_list(
        track_id=track.id, criteria=WrokoutCriteria(limit=5)
    )

    assert len(workouts_from_db) == 5
    assert workouts_from_db[0].date == date(2025, 1, count), (
        "Последняя тренировка должна быть первой"
    )
