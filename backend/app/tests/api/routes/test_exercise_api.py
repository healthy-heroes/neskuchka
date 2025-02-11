from app.domain.exercise import ExerciseRepository
from app.tests.fixtures.domain import create_exercise


def test_get_exercise(client, exercise_repo: ExerciseRepository):
    exercise = create_exercise()
    exercise_repo.add(exercise)

    response = client.get(f"/api/v1/exercises/{exercise.slug}")

    assert response.status_code == 200
    assert response.json() == exercise.model_dump()


def test_get_non_existent_exercise(client, exercise_repo: ExerciseRepository):
    response = client.get("/api/v1/exercises/non_existent_exercise")

    assert response.status_code == 404
