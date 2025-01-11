from typing import NewType
from app.domain.entity import EntityModel

ExerciseId = NewType("ExerciseId", int)
ExerciseSlug = NewType("ExerciseSlug", str)

class Exercise(EntityModel):
    """
    Модель базового упражнения
    
    Например, push-up, pull-up и подобное
    """
    id: ExerciseId
    slug: ExerciseSlug
    name: str
    description: str


class ExerciseRepository():
    def add(exercise: Exercise) -> Exercise:
        pass

    def get_all() -> list[Exercise]:
        pass

    def get_by_slug(slug: str) -> Exercise:
        pass
