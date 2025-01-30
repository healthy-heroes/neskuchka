from sqlmodel import Field, SQLModel, Session, select

from app.domain.track import TrackRepository, Track, TrackId
from app.domain.user import UserId


class TrackModel(SQLModel, table=True):
    __tablename__ = "track"

    id: str = Field(default=None, primary_key=True)
    name: str
    owner_id: str

    def to_domain(self) -> Track:
        return Track(
            id=TrackId(self.id),
            name=self.name,
            owner_id=UserId(self.owner_id),
        )


class TrackDbRepository(TrackRepository):
    def __init__(self, session: Session):
        self.session = session

    def get_main_track(self) -> Track | None:
        query = select(TrackModel).limit(1)
        result = self.session.exec(query).first()

        if not result:
            return None

        return result.to_domain()

    def get_by_id(self, track_id: TrackId) -> Track | None:
        query = select(TrackModel).where(TrackModel.id == track_id).limit(1)
        result = self.session.exec(query).first()

        if not result:
            return None

        return result.to_domain()
