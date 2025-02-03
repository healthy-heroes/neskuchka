import pytest
from sqlmodel import SQLModel, Session
from fastapi.testclient import TestClient

from app.config.test import TestAppSettings
from app.infrastructure.db.database import db
from app.main import app


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


# client fixture
@pytest.fixture(name="client")
def client():
    yield TestClient(app)
    app.dependency_overrides.clear()


@pytest.fixture(name="session")
def session():
    # TODO: remove copypaste from database.py
    with Session(db.engine) as session:
        yield session
