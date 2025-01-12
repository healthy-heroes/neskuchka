from fastapi import APIRouter

from app.api.deps import ExerciseRepoDependency
from app.domain.exercise import Exercise, ExerciseSlug

router = APIRouter(prefix="/exercises", tags=["exercises"])


@router.post("/", status_code=201)
async def create_exercise(
    item: Exercise, exercise_repo: ExerciseRepoDependency
) -> Exercise:
    exercise = Exercise(**item.model_dump())
    return exercise_repo.add(exercise)


@router.get("/")
async def get_exercises(exercise_repo: ExerciseRepoDependency) -> list[Exercise]:
    return exercise_repo.get_all()


@router.get("/{slug}")
async def get_exersice(slug: str, exercise_repo: ExerciseRepoDependency) -> Exercise:
    return exercise_repo.get_by_slug(ExerciseSlug(slug))
