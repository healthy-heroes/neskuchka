from fastapi import APIRouter, HTTPException

from app.api.deps import TrackRepoDependency
from app.domain.track import Track


router = APIRouter(prefix="/tracks", tags=["tracks"])


@router.get("/main")
async def get_main_track(track_repo: TrackRepoDependency) -> Track:
    track = track_repo.get_main_track()

    if not track:
        raise HTTPException(status_code=404, detail="Main track not found")

    return track
