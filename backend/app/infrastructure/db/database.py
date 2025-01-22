from typing import Generator
from sqlmodel import Session, create_engine

from app.config.default import DatabaseSettings
from app.config.main import settings


class Database:
    def __init__(self, settings: DatabaseSettings):
        self.engine = create_engine(settings.dsn, **settings.engine_args.model_dump())

    def dispose(self):
        self.engine.dispose()

    def session_getter(self) -> Generator[Session, None]:
        with Session(self.engine) as session:
            yield session


db = Database(settings=settings.database)
