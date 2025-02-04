from json import loads
from sqlalchemy import Date
from sqlmodel import Field, SQLModel, Session, select
from app.domain.workout import WorkoutRepository, Workout, WorkoutId, WrokoutCriteria
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

    @staticmethod
    def from_domain(workout: Workout) -> "WorkoutModel":
        print(f"workout: {workout.model_dump(exclude={"id", "track_id", "sections"})}")

        return WorkoutModel(
            id=str(workout.id),
            track_id=str(workout.track_id),
            sections=workout.model_dump_json(include={"sections"}),
            **workout.model_dump(exclude={"id", "track_id", "sections"}),
        )


class WorkoutDbRepository(WorkoutRepository):
    def __init__(self, session: Session):
        self.session = session

    def add(self, workout: Workout) -> Workout:
        db_workout = WorkoutModel.from_domain(workout)

        self.session.add(db_workout)
        self.session.commit()
        self.session.refresh(db_workout)

    def get_by_id(self, workout_id: WorkoutId) -> Workout | None:
        query = select(WorkoutModel).where(WorkoutModel.id == str(workout_id)).limit(1)
        result = self.session.exec(query).first()

        if not result:
            return None

        return result.to_domain()

    def get_list(self, track_id: TrackId, criteria: WrokoutCriteria) -> list[Workout]:
        query = (
            select(WorkoutModel)
            .where(WorkoutModel.track_id == str(track_id))
            .limit(criteria.limit)
        )
        result = self.session.exec(query).all()

        return [result.to_domain() for result in result]
