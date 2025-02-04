from fastapi import APIRouter, Depends, HTTPException

from app.api.deps import TrackRepoDependency, WorkoutRepoDependency
from app.domain.track import Track
from app.domain.workout import Workout, WrokoutCriteria


router = APIRouter(prefix="/tracks", tags=["tracks"])


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
    workout_repo: WorkoutRepoDependency, track: Track = Depends(get_track_or_404)
) -> list[Workout]:
    return workout_repo.get_list(track.id, WrokoutCriteria(limit=10))
