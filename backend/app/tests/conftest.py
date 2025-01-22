import os
import pytest
from sqlmodel import SQLModel

from app.infrastructure.db.session import engine


def pytest_configure(config):
    os.environ["ENVIRONMENT"] = "test"


# Init testing database
@pytest.fixture(scope="session", autouse=True)
def init_db():
    print("init_db:Starting up")
    # TODO: пока создается
    # SQLModel.metadata.create_all(bind=engine)

    yield

    print("init_db:Shutting down")
    SQLModel.metadata.drop_all(engine)
