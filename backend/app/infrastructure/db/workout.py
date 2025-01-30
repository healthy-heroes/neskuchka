from json import loads
from sqlalchemy import Date
from sqlmodel import Field, SQLModel, Session, select
from app.domain.workout import WorkoutRepository, Workout, WorkoutId
from app.domain.track import TrackId


class WorkoutModel(SQLModel, table=True):
    __tablename__ = "workout"

    id: str = Field(
        default=None,
        primary_key=True,
    )
    date: str = Field(sa_type=Date)
    track_id: str
    sections: str

    def to_domain(self) -> Workout:
        return Workout(
            id=WorkoutId(self.id),
            date=self.date,
            track_id=TrackId(self.track_id),
            sections=loads(self.sections)["sections"],
        )


class WorkoutDbRepository(WorkoutRepository):
    def __init__(self, session: Session):
        self.session = session

    def add(self, workout: Workout) -> Workout:
        db_workout = WorkoutModel(
            id=str(workout.id),
            date=workout.date,
            track_id=str(workout.track_id),
            sections=workout.model_dump_json(include={"sections"}),
        )

        self.session.add(db_workout)
        self.session.commit()
        self.session.refresh(db_workout)

    def get_by_id(self, workout_id: WorkoutId) -> Workout | None:
        query = select(WorkoutModel).where(WorkoutModel.id == str(workout_id)).limit(1)
        result = self.session.exec(query).first()

        if not result:
            return None

        return result.to_domain()
