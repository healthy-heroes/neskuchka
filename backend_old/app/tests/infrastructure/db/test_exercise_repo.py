from sqlmodel import Session
from app.domain.exercise import ExerciseCriteria, ExerciseSlug
from app.infrastructure.db.exercise import ExerciseDbRepository
from app.tests.fixtures.domain import create_exercise


def test_add_get_exercise(session: Session):
    exercise = create_exercise()

    exercise_repository = ExerciseDbRepository(session)
    exercise_repository.add(exercise)

    exercise_from_db = exercise_repository.get_by_slug(exercise.slug)
    assert exercise_from_db == exercise


def test_find_exercises(session: Session):
    exercise_repository = ExerciseDbRepository(session)
    exercise_repository.add(create_exercise(ExerciseSlug("slug1")))
    exercise_repository.add(create_exercise(ExerciseSlug("slug2")))
    exercise_repository.add(create_exercise(ExerciseSlug("slug3")))

    exercises = exercise_repository.find(
        ExerciseCriteria(slugs=["slug1", "slug2", "slug4"])
    )
    assert len(exercises) == 2
    assert exercises[0].slug == "slug1"
    assert exercises[1].slug == "slug2"
