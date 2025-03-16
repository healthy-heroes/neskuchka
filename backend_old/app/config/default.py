from enum import StrEnum
from pydantic import BaseModel
from pydantic_settings import BaseSettings as PydanticBaseSettings


class Environment(StrEnum):
    LOCAL = "local"
    TEST = "test"


class DatabaseEngineArgs(BaseModel):
    # SQLAlchemy engine options
    # See: sqlalchemy.create_engine() for full documentation
    connect_args: dict | None = (
        None  # Additional arguments passed to DBAPI upon connect
    )
    echo: bool = False  # If True, engine will log all statements
    # echo_pool: bool = False  # If True, connection pool will log all checkouts/checkins
    # expire_on_commit: bool = False  # If True, all instances will be expired after each commit
    # max_overflow: int = 10  # Max number of connections to allow in connection pool "overflow"
    # pool_size: int = 10  # The size of the pool to be maintained


class DatabaseSettings(BaseModel):
    dsn: str
    engine_args: DatabaseEngineArgs


class BaseSettings(PydanticBaseSettings):
    app_name: str = "Neskuchka club"
    environment: Environment = "local"

    database: DatabaseSettings


class AppSettings(BaseSettings):
    API_V1_PREFIX: str = "/api/v1"

    database: DatabaseSettings = DatabaseSettings(
        dsn="sqlite:///database.db",
        engine_args=DatabaseEngineArgs(
            connect_args={"check_same_thread": False},
        ),
    )


def get_env() -> Environment:
    return AppSettings().environment
