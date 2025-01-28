from collections.abc import Generator
from typing import Annotated

from fastapi import Depends
from sqlmodel import Session

from app.infrastructure.db.database import db
from app.infrastructure.db.exercise import ExerciseDbRepository
from app.domain.exercise import ExerciseRepository
from app.domain.track import RepositoryTrack
from app.domain.user import RepositoryUser
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
) -> Generator[RepositoryUser, None, None]:
    try:
        yield UserDbRepository(session)
    finally:
        session.close()


UserRepoDependency = Annotated[RepositoryUser, Depends(get_user_repository)]


# Track repository dependency per request
def get_track_repository(
    session: SessionDependency,
) -> Generator[RepositoryTrack, None, None]:
    try:
        yield TrackDbRepository(session)
    finally:
        session.close()


TrackRepoDependency = Annotated[RepositoryTrack, Depends(get_track_repository)]
