from sqlmodel import Session

from app.infrastructure.db.track import TrackDbRepository
from app.tests.fixtures.domain import create_track


def test_add_get_track(session: Session):
    track = create_track()

    track_repository = TrackDbRepository(session)

    track_from_db = track_repository.get_main_track()
    assert track_from_db is None, "Should be no main track if not added"

    track_repository.add(track)

    track_from_db = track_repository.get_by_id(track.id)
    assert track_from_db == track, "Should be able to get track by id"

    track_from_db = track_repository.get_main_track()
    assert track_from_db == track, "Should be able to get main track"
