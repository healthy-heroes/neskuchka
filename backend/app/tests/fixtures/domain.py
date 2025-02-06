import uuid

from app.domain.exercise import Exercise, ExerciseSlug


def create_exercise() -> Exercise:
    return Exercise(
        slug=ExerciseSlug(str(uuid.uuid4())),
        name="Some exercise",
        description="Some exercise description",
    )
