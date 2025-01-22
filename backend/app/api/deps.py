from collections.abc import Generator
from typing import Annotated

from fastapi import Depends
from sqlmodel import Session

from app.infrastructure.db.database import db
from app.infrastructure.db.exercise import ExerciseDbRepository
from app.domain.exercise import ExerciseRepository

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
