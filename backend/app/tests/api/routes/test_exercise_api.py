from pytest_mock import MockerFixture

from app.api.deps import get_exercise_repository
from app.domain.exercise import Exercise, ExerciseRepository, ExerciseSlug


def test_get_exercise(client, mocker: MockerFixture):
    test_exercise = Exercise(
        slug=ExerciseSlug("test_exercise"),
        name="Test exercise",
        description="Test description",
    )

    repo = ExerciseRepository()
    mocker.patch.object(repo, "get_by_slug", return_value=test_exercise)
    client.app.dependency_overrides[get_exercise_repository] = lambda: repo

    response = client.get("/api/v1/exercises/test_exercise")

    assert repo.get_by_slug.call_count == 1
    assert repo.get_by_slug.call_args[0] == (test_exercise.slug,)

    assert response.status_code == 200
    assert response.json() == test_exercise.model_dump()


def test_get_non_existent_exercise(client, mocker: MockerFixture):
    repo = ExerciseRepository()
    mocker.patch.object(repo, "get_by_slug", return_value=None)
    client.app.dependency_overrides[get_exercise_repository] = lambda: repo

    response = client.get("/api/v1/exercises/non_existent_exercise")

    assert repo.get_by_slug.call_count == 1
    assert response.status_code == 404
