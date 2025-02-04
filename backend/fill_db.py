from sqlmodel import SQLModel

from app.domain.track import Track
from app.domain.user import User

from app.infrastructure.db.database import db
from app.infrastructure.db.user import UserDbRepository
from app.infrastructure.db.track import TrackDbRepository
from app.domain.exercise import Exercise, ExerciseSlug
from app.infrastructure.db.exercise import ExerciseDbRepository

session = next(db.session_getter())

# Clear database
SQLModel.metadata.drop_all(session.get_bind())
SQLModel.metadata.create_all(session.get_bind())

session.begin()

# create user
user_repo = UserDbRepository(session)

user_login = "first_user"
user_id = User.create_id(user_login)
user = User(
    id=user_id, name="First User", login=user_login, email="first_user@example.com"
)

print(f"create user: {user}")
user_repo.add(user)

# create track
track_repo = TrackDbRepository(session)

track = Track(name="Neskuchka", owner_id=user_id)

print(f"create track: {track}")
track_repo.add(track)

# create exercises
exercise_repo = ExerciseDbRepository(session)

exercises = [
    Exercise(
        name="Push-up",
        slug=ExerciseSlug("push-up"),
        description="Push-up is a bodyweight exercise that targets the chest, shoulders, and triceps. It is performed by lying on your back with your hands and feet on the ground, and then pushing yourself up and down using your arms.",
    ),
    Exercise(
        name="Pull-up",
        slug=ExerciseSlug("pull-up"),
        description="Pull-up is a bodyweight exercise that targets the back, shoulders, and biceps. It is performed by hanging from a bar with your hands and feet on the ground, and then pulling yourself up using your arms.",
    ),
    Exercise(
        name="Air Squat",
        slug=ExerciseSlug("air-squat"),
        description="Air Squat is a bodyweight exercise that targets the legs, glutes, and core. It is performed by standing with your feet shoulder-width apart, and then squatting down and up using your legs.",
    ),
    Exercise(
        name="Deadlift",
        slug=ExerciseSlug("deadlift"),
        description="Deadlift is a bodyweight exercise that targets the back, shoulders, and biceps. It is performed by standing with your feet shoulder-width apart, and then squatting down and up using your legs.",
    ),
    Exercise(
        name="Bench Press",
        slug=ExerciseSlug("bench-press"),
        description="Bench Press is a bodyweight exercise that targets the chest, shoulders, and triceps. It is performed by lying on your back with your hands and feet on the ground, and then pushing yourself up and down using your arms.",
    ),
]

for exercise in exercises:
    print(f"create exercise: {exercise}")
    exercise_repo.add(exercise)

# create workout

# todo

session.commit()
