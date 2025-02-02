from sqlmodel import Field, SQLModel, Session, select

from app.domain.user import UserRepository, User, UserId


class UserModel(SQLModel, table=True):
    __tablename__ = "user"

    id: str = Field(default=None, primary_key=True)
    name: str
    login: str
    email: str

    def to_domain(self) -> User:
        return User(
            id=UserId(self.id),
            name=self.name,
            login=self.login,
            email=self.email,
        )

    @staticmethod
    def from_domain(user: User) -> "UserModel":
        return UserModel(
            id=str(user.id),
            **user.model_dump(exclude={"id"}),
        )


class UserDbRepository(UserRepository):
    def __init__(self, session: Session):
        self.session = session

    def add(self, user: User) -> User:
        db_user = UserModel.from_domain(user)

        self.session.add(db_user)
        self.session.commit()

    def get_by_id(self, user_id: UserId) -> User | None:
        query = select(UserModel).where(UserModel.id == str(user_id)).limit(1)
        result = self.session.exec(query).first()

        if not result:
            return None

        return result.to_domain()
