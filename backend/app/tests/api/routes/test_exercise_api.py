import pytest
from pytest_mock import MockerFixture

from app.api.deps import get_exercise_repository
from app.domain.exercise import ExerciseRepository

from app.tests.fixtures.domain import create_exercise


@pytest.fixture(scope="function")
def exercise_repo(client, mocker: MockerFixture):
    repo = mocker.Mock(spec=ExerciseRepository)
    client.app.dependency_overrides[get_exercise_repository] = lambda: repo

    yield repo

    client.app.dependency_overrides.pop(get_exercise_repository)
    return repo


def test_get_exercise(client, mocker: MockerFixture, exercise_repo: ExerciseRepository):
    exercise = create_exercise()

    mocker.patch.object(exercise_repo, "get_by_slug", return_value=exercise)

    response = client.get(f"/api/v1/exercises/{exercise.slug}")

    assert exercise_repo.get_by_slug.call_count == 1
    assert exercise_repo.get_by_slug.call_args[0] == (exercise.slug,)

    assert response.status_code == 200
    assert response.json() == exercise.model_dump()


def test_get_non_existent_exercise(
    client, mocker: MockerFixture, exercise_repo: ExerciseRepository
):
    mocker.patch.object(exercise_repo, "get_by_slug", return_value=None)

    response = client.get("/api/v1/exercises/non_existent_exercise")

    assert exercise_repo.get_by_slug.call_count == 1
    assert response.status_code == 404
