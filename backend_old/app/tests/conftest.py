import pytest
from sqlmodel import SQLModel, Session
from fastapi.testclient import TestClient

from app.api.deps import (
    get_exercise_repository,
    get_track_repository,
    get_workout_repository,
)
from app.config.test import TestAppSettings
from app.infrastructure.db.database import db
from app.main import app
from app.tests.fixtures.domain import (
    ExerciseRepositoryTest,
    TrackRepositoryTest,
    WorkoutRepositoryTest,
)


settings = TestAppSettings()


# Init testing database
@pytest.fixture(scope="session", autouse=True)
def init_db():
    # если импортить рано, то настройки инциализируются раньше выставления ENVIRONMENT
    # from app.infrastructure.db.database import db

    print("init_db:Starting up")
    SQLModel.metadata.create_all(db.engine)

    yield

    print("init_db:Shutting down")
    SQLModel.metadata.drop_all(db.engine)
    db.dispose()


@pytest.fixture(name="client")
def client():
    yield TestClient(app)
    app.dependency_overrides.clear()


@pytest.fixture(name="session", scope="function")
def session():
    # TODO: remove copypaste from database.py
    with Session(db.engine) as session:
        yield session


# repositories
@pytest.fixture(scope="function")
def exercise_repo(client):
    repo = ExerciseRepositoryTest()
    client.app.dependency_overrides[get_exercise_repository] = lambda: repo

    yield repo

    client.app.dependency_overrides.pop(get_exercise_repository)
    return repo


@pytest.fixture(scope="function")
def track_repo(client):
    repo = TrackRepositoryTest()
    client.app.dependency_overrides[get_track_repository] = lambda: repo

    yield repo

    client.app.dependency_overrides.pop(get_track_repository)


@pytest.fixture(scope="function")
def workout_repo(client):
    repo = WorkoutRepositoryTest()
    client.app.dependency_overrides[get_workout_repository] = lambda: repo

    yield repo

    client.app.dependency_overrides.pop(get_workout_repository)
