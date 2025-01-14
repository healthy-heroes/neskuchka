from typing import Annotated
from fastapi import Depends, FastAPI
from pydantic import BaseModel

from app.api.main import api_router
from app.infrastructure.db.session import create_db_and_tables
from app.config import settings, SettingsDependency

app = FastAPI()


# todo: убрать
@app.on_event("startup")
def on_startup():
    create_db_and_tables()


app.include_router(api_router, prefix=settings.API_V1_PREFIX)


class Info(BaseModel):
    name: str
    environment: str


@app.get("/")
def info(settings: SettingsDependency):
    return Info(name=settings.app_name, environment=settings.environment)
