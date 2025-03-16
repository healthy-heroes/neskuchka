from typing import NewType

from pydantic import ConfigDict
from app.domain.entity import EntityModel, CriteriaModel

ExerciseSlug = NewType("ExerciseSlug", str)


class Exercise(EntityModel):
    """
    Модель базового упражнения

    Например, push-up, pull-up и подобное
    """

    model_config = ConfigDict(frozen=True)

    slug: ExerciseSlug
    name: str
    description: str


class ExerciseCriteria(CriteriaModel):
    slugs: list[ExerciseSlug] | None = None


class ExerciseRepository:
    def add(self, exercise: Exercise) -> Exercise:
        pass

    def get_all(self) -> list[Exercise]:
        pass

    def get_by_slug(self, slug: ExerciseSlug) -> Exercise | None:
        pass

    def find(self, criteria: ExerciseCriteria) -> list[Exercise]:
        pass
