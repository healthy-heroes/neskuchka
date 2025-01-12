from sqlmodel import Field, SQLModel, Session, select

from app.domain.exercise import Exercise, ExerciseRepository, ExerciseSlug


class ExerciseModel(SQLModel, table=True):
    __tablename__ = "exercise"

    slug: str = Field(default=None, primary_key=True)
    name: str
    description: str

    def to_domain(self) -> Exercise:
        return Exercise(
            name=self.name, slug=ExerciseSlug(self.slug), description=self.description
        )


class ExerciseDbRepository(ExerciseRepository):
    def __init__(self, session: Session):
        self.session = session

    def add(self, exercise: Exercise) -> Exercise:
        db_exercise = ExerciseModel(
            name=exercise.name, slug=exercise.slug, description=exercise.description
        )
        self.session.add(db_exercise)
        self.session.commit()
        self.session.refresh(db_exercise)

        return db_exercise.to_domain()

    def get_all(self) -> list[Exercise]:
        exercises = self.session.exec(select(ExerciseModel)).all()

        return [ex.to_domain() for ex in exercises]

    def get_by_slug(self, slug: ExerciseSlug) -> Exercise | None:
        exercise = self.session.exec(
            select(ExerciseModel).where(ExerciseModel.slug == slug)
        ).first()

        if not exercise:
            return None

        return exercise.to_domain()
