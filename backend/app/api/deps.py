from collections.abc import Generator
from typing import Annotated

from fastapi import Depends
from sqlmodel import Session

from app.infrastructure.db.database import db
from app.infrastructure.db.exercise import ExerciseDbRepository
from app.domain.exercise import ExerciseRepository
from app.domain.track import TrackRepository
from app.domain.user import UserRepository
from app.infrastructure.db.track import TrackDbRepository
from app.infrastructure.db.user import UserDbRepository

# Session dependency per request
SessionDependency = Annotated[Session, Depends(db.session_getter)]


# Exercise repository dependency per request
def get_exercise_repository(
    session: SessionDependency,
) -> Generator[ExerciseRepository, None, None]:
    try:
        yield ExerciseDbRepository(session)
    finally:
        session.close()


ExerciseRepoDependency = Annotated[ExerciseRepository, Depends(get_exercise_repository)]


# User repository dependency per request
def get_user_repository(
    session: SessionDependency,
) -> Generator[UserRepository, None, None]:
    try:
        yield UserDbRepository(session)
    finally:
        session.close()


UserRepoDependency = Annotated[UserRepository, Depends(get_user_repository)]


# Track repository dependency per request
def get_track_repository(
    session: SessionDependency,
) -> Generator[TrackRepository, None, None]:
    try:
        yield TrackDbRepository(session)
    finally:
        session.close()


TrackRepoDependency = Annotated[TrackRepository, Depends(get_track_repository)]
