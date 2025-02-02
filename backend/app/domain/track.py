from typing import Annotated
from uuid import UUID
from pydantic import ConfigDict
from pydantic.types import UuidVersion

from app.domain.entity import EntityModel
from app.domain.user import UserId


TrackId = Annotated[UUID, UuidVersion(5)]


class Track(EntityModel):
    """
    Модель программы тренировок (трека)
    """

    model_config = ConfigDict(frozen=True)

    id: TrackId
    name: str
    owner_id: UserId


class TrackRepository:
    def get_main_track(self) -> Track | None:
        """
        Временный метод, так как пока будет всего один трек
        """
        pass

    def add(self, track: Track) -> Track:
        pass

    def get_by_id(self, track_id: TrackId) -> Track | None:
        pass
