from fastapi import FastAPI

from app.api.main import api_router
from app.infrastructure.db.session import create_db_and_tables

app = FastAPI()


# todo: убрать
@app.on_event("startup")
def on_startup():
    create_db_and_tables()


app.include_router(api_router, prefix="/api/v1")
