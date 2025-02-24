from fastapi import APIRouter, Depends, HTTPException
from pydantic import BaseModel

from app.api.deps import (
    ExerciseRepoDependency,
    TrackRepoDependency,
    WorkoutRepoDependency,
)
from app.domain.track import Track
from app.domain.workout import Workout, WrokoutCriteria
from app.domain.exercise import Exercise, ExerciseCriteria, ExerciseSlug

router = APIRouter(prefix="/tracks", tags=["tracks"])


class TrackWorkoutsSchema(BaseModel):
    workouts: list[Workout]
    exercises: dict[ExerciseSlug, Exercise]


def get_track_or_404(track_repo: TrackRepoDependency) -> Track:
    track = track_repo.get_main_track()

    if not track:
        raise HTTPException(status_code=404, detail="Main track not found")

    return track


@router.get("/main")
async def get_main_track(track: Track = Depends(get_track_or_404)) -> Track:
    return track


@router.get("/main/last_workouts")
async def get_main_track_workouts(
    workout_repo: WorkoutRepoDependency,
    exercise_repo: ExerciseRepoDependency,
    track: Track = Depends(get_track_or_404),
) -> TrackWorkoutsSchema:
    workouts = workout_repo.get_list(track.id, WrokoutCriteria(limit=10))

    exercise_slugs = {
        workoout_exercise.exercise_slug
        for workout in workouts
        for section in workout.sections
        for workoout_exercise in section.exercises
    }

    exercises = exercise_repo.find(ExerciseCriteria(slugs=list(exercise_slugs)))

    return TrackWorkoutsSchema(
        workouts=workouts, exercises={ex.slug: ex for ex in exercises}
    )
