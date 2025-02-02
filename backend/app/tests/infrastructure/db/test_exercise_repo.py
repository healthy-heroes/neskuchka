from sqlmodel import Session
from app.domain.exercise import Exercise, ExerciseSlug
from app.infrastructure.db.exercise import ExerciseDbRepository


def test_add_get_exercise(session: Session):
    exercise_slug = "test_exercise"
    exercise = Exercise(
        slug=ExerciseSlug(exercise_slug),
        name="Test exercise",
        description="Test description",
    )

    exercise_repository = ExerciseDbRepository(session)
    exercise_repository.add(exercise)

    exercise_from_db = exercise_repository.get_by_slug(exercise_slug)
    assert exercise_from_db == exercise
