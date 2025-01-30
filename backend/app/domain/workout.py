from dataclasses import dataclass
from enum import StrEnum
from typing import NewType
from datetime import date
from pydantic import Field
from ulid import ULID

from app.domain.entity import EntityModel
from app.domain.exercise import ExerciseSlug
from app.domain.track import TrackId


# todo: Нужно заменить на uuid
WorkoutId = NewType("WorkoutId", ULID)


class WorkoutProtocolType(StrEnum):
    DEFAULT = "DEFAULT"


class WorkoutProtocol(EntityModel):
    """
    Протокол тренировки
    """

    type: WorkoutProtocolType = Field(default=WorkoutProtocolType.DEFAULT)
    title: str
    description: str = ""


class WorkoutExercise(EntityModel):
    """
    Упражнение в тренировке
    """

    exercise_slug: ExerciseSlug

    # Количество повторов
    repetitions: int | None = None
    # Текстовое описание повторов
    repetitions_text: str = ""

    # Вес (кг)
    weight: int | None = None
    # Текстовое описание веса
    weight_text: str = ""


class WorkoutSection(EntityModel):
    """
    Часть тренировки
    """

    title: str
    protocol: WorkoutProtocol
    exercises: list[WorkoutExercise] = Field(min_length=1)


class Workout(EntityModel):
    """
    Модель тренировки

    Может состоять из нескольких частей
    """

    id: WorkoutId = Field(default_factory=ULID)

    date: date
    track_id: TrackId

    # Упорядоченный список частей
    sections: list[WorkoutSection] = Field(min_length=1)


@dataclass(frozen=True)
class WrokoutCriteria:
    limit: int


class WorkoutRepository:
    def add(self, workout: Workout) -> Workout:
        pass

    def get_by_id(self, workout_id: WorkoutId) -> Workout | None:
        pass
