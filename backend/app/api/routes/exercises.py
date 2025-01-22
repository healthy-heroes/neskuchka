from fastapi import APIRouter, HTTPException
import sqlalchemy

from app.api.deps import ExerciseRepoDependency
from app.domain.exercise import Exercise, ExerciseSlug

router = APIRouter(prefix="/exercises", tags=["exercises"])


@router.post("/", status_code=201)
async def create_exercise(
    item: Exercise, exercise_repo: ExerciseRepoDependency
) -> Exercise:
    exercise = Exercise(**item.model_dump())
    try:
        return exercise_repo.add(exercise)
    except sqlalchemy.exc.IntegrityError:
        raise HTTPException(
            status_code=400, detail=f"Упражнение с slug={exercise.slug} уже существует"
        )


@router.get("/")
async def get_exercises(exercise_repo: ExerciseRepoDependency) -> list[Exercise]:
    return exercise_repo.get_all()


@router.get("/{slug}")
async def get_exersice(slug: str, exercise_repo: ExerciseRepoDependency) -> Exercise:
    exercise = exercise_repo.get_by_slug(ExerciseSlug(slug))
    if not exercise:
        raise HTTPException(
            status_code=404, detail=f"Упражнение с slug={slug} не найдено"
        )
    return exercise
