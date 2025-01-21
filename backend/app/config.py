from typing import Annotated
from fastapi import Depends
from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    app_name: str = "Neskuchka club"

    environment: str = "local"

    API_V1_PREFIX: str = "/api/v1"


settings = Settings()


def get_settings():
    return settings


SettingsDependency = Annotated[Settings, Depends(get_settings)]
