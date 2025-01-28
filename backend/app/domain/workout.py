from dataclasses import dataclass
from typing import NewType
from datetime import date
from app.domain.entity import EntityModel
from app.domain.exercise import ExerciseId
from app.domain.track import TrackId


# todo: Нужно заменить на uuid
WorkoutId = NewType("WorkoutId", int)


class WorkoutExercise(EntityModel):
    """
    Упражнение в тренировке
    """

    exercise_id: ExerciseId

    # Числовые характеристики
    # Количество повторов
    repetitions: int

    # Вес (кг)
    weight: int


class WorkoutSection(EntityModel):
    """
    Часть тренировки
    """

    title: str
    # todo: должно превратиться в полноценный домен
    schema: str
    # todo: non empty list Field(min_length)
    exercises: list[WorkoutExercise]


class Workout(EntityModel):
    """
    Модель тренировки

    Может состоять из нескольких частей
    """

    id: WorkoutId

    date: date
    track_id: TrackId

    # todo: non empty list
    # Упорядоченный список частей
    sections: list[WorkoutSection]


@dataclass(frozen=True)
class WrokoutCriteria:
    limit: int


class WorkoutRepository:
    def add(self, workout: Workout) -> Workout:
        pass

    def get_by_id(self, workout_id: WorkoutId) -> Workout | None:
        pass
