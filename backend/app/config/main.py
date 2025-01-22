import os
from typing import Annotated

from fastapi import Depends

from app.config.default import AppSettings, Environment, get_env


def init_settings():
    env = os.getenv("ENVIRONMENT")

    if env == Environment.TEST:
        from app.config.test import TestAppSettings

        return TestAppSettings()

    return AppSettings()


env = get_env()
settings = init_settings()


def get_settings():
    return settings


SettingsDependency = Annotated[AppSettings, Depends(get_settings)]
