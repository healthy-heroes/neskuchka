from sqlmodel import Session
from app.infrastructure.db.exercise import ExerciseDbRepository
from app.tests.fixtures.domain import create_exercise


def test_add_get_exercise(session: Session):
    exercise = create_exercise()

    exercise_repository = ExerciseDbRepository(session)
    exercise_repository.add(exercise)

    exercise_from_db = exercise_repository.get_by_slug(exercise.slug)
    assert exercise_from_db == exercise
