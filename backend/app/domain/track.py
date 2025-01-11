from typing import NewType
from app.domain.entity import EntityModel
from app.domain.user import UserId


# todo: Нужно заменить на uuid
TrackId = NewType("TrackId", int)


class Track(EntityModel):
    """
    Модель программы тренировок (трека)
    """

    # todo: uuid?
    id: TrackId
    name: str
    # todo:
    owner_id: UserId


class RepositoryTrack:
    def add(track: Track) -> Track:
        pass

    def get_by_id(track_id: TrackId) -> Track:
        pass
