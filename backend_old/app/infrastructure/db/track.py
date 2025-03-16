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

    @staticmethod
    def from_domain(track: Track) -> "TrackModel":
        return TrackModel(
            id=str(track.id),
            owner_id=str(track.owner_id),
            **track.model_dump(exclude={"id", "owner_id"}),
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

    def add(self, track: Track) -> Track:
        db_track = TrackModel.from_domain(track)

        self.session.add(db_track)
        self.session.commit()

    def get_by_id(self, track_id: TrackId) -> Track | None:
        query = select(TrackModel).where(TrackModel.id == str(track_id)).limit(1)
        result = self.session.exec(query).first()

        if not result:
            return None

        return result.to_domain()
