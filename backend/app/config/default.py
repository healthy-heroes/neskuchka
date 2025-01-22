from enum import StrEnum
from pydantic_settings import BaseSettings as PydanticBaseSettings


class Environment(StrEnum):
    LOCAL = "local"
    TEST = "test"


class BaseSettings(PydanticBaseSettings):
    app_name: str = "Neskuchka club"
    environment: Environment = "local"


class AppSettings(BaseSettings):
    API_V1_PREFIX: str = "/api/v1"


def get_env() -> Environment:
    return AppSettings().environment
