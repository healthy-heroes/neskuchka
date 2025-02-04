from datetime import date
from sqlmodel import SQLModel

from app.domain.track import Track
from app.domain.user import User
from app.domain.exercise import Exercise, ExerciseSlug
from app.domain.workout import Workout, WorkoutSection, WorkoutProtocol, WorkoutExercise

from app.infrastructure.db.database import db
from app.infrastructure.db.user import UserDbRepository
from app.infrastructure.db.track import TrackDbRepository
from app.infrastructure.db.exercise import ExerciseDbRepository
from app.infrastructure.db.workout import WorkoutDbRepository

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

track = Track(name="Нескучный спорт", owner_id=user_id)

print(f"create track: {track}")
track_repo.add(track)

# create exercises
exercise_repo = ExerciseDbRepository(session)
    
exercises = {
    'Раскрытия в планке': Exercise(name="Раскрытия в планке", slug=ExerciseSlug("plank-hip-opening"), description=""),
    'Ягодичные марши': Exercise(name="Ягодичные марши", slug=ExerciseSlug("glute-march"), description=""),
    'Джампинг джек': Exercise(name="Джампинг джек", slug=ExerciseSlug("jumping-jack"), description=""),
    'C пола на грудь + 2 выпада назад': Exercise(name="C пола на грудь + 2 выпада назад", slug=ExerciseSlug("push-up-with-back-drop"), description=""),
    'Становая на одной': Exercise(name="Становая на одной", slug=ExerciseSlug("deadlift-on-one-leg"), description=""),
    'Пресс на прямых руки над головой': Exercise(name="Пресс на прямых руки над головой", slug=ExerciseSlug("situps-with-hands-over-head"), description=""),
    
    'Стол': Exercise(name="Стол", slug=ExerciseSlug("table"), description=""),
    'Наклоны вперед': Exercise(name="Наклоны вперед", slug=ExerciseSlug("forward-bend"), description=""),
    'Качающиеся планки': Exercise(name="Качающиеся планки", slug=ExerciseSlug("plank-with-jumping-jack"), description=""),
    'Становая тяга': Exercise(name="Становая тяга", slug=ExerciseSlug("deadlift"), description=""),
    'Приседания': Exercise(name="Приседания", slug=ExerciseSlug("squats"), description=""),
    'С груди над головой': Exercise(name="С груди над головой", slug=ExerciseSlug("push-up-with-hands-over-head"), description=""),
}

for exercise in exercises.values():
    print(f"create exercise: {exercise}")
    exercise_repo.add(exercise)

# create workout
workout_repo = WorkoutDbRepository(session)

"""
31 Jan
Разминка
3 раунда:
- 20 раскрытий в планке
- 20 ягодичных маршей
- 30 джампинг Джеков

Комплекс
5 раундов не на время:
- 10 с пола на грудь + 2 выпада назад*
- 10+10 становая на одной*
- 10 пресс на прямых руки над головой*
"""

workout = Workout(
    date=date(2025, 1, 31),
    track_id=track.id,

    sections=[
        WorkoutSection(
            title="Разминка", 
            protocol=WorkoutProtocol(title="3 раунда"),
            exercises=[
                WorkoutExercise(exercise_slug=exercises["Раскрытия в планке"].slug, repetitions=20),
                WorkoutExercise(exercise_slug=exercises["Ягодичные марши"].slug, repetitions=20),
                WorkoutExercise(exercise_slug=exercises["Джампинг джек"].slug, repetitions=30),
            ],
        ),
        WorkoutSection(
            title="Комплекс", 
            protocol=WorkoutProtocol(title="5 раундов"),
            description="*можно использовать спортивные снаряды",
            exercises=[
                WorkoutExercise(exercise_slug=exercises["C пола на грудь + 2 выпада назад"].slug, repetitions=10),
                WorkoutExercise(exercise_slug=exercises["Становая на одной"].slug, repetitions=20, repetitions_text="10+10"),
                WorkoutExercise(exercise_slug=exercises["Пресс на прямых руки над головой"].slug, repetitions=10),
            ],
        ),
    ]
)

print(f"create workout: {workout}")
workout_repo.add(workout)


"""
3 feb 
Разминка
3 раунда:
- 10 столов
- 10 наклонов вперед
- 20 качающихся планок на локтях

Комплекс
5 раундов на время:
- 24 становых*
- 18 приседаний*
- 12 с груди над головой*
"""

workout = Workout(
    date=date(2025, 2, 3),
    track_id=track.id,
    sections=[
        WorkoutSection(
            title="Разминка", 
            protocol=WorkoutProtocol(title="3 раунда"),
            exercises=[
                WorkoutExercise(exercise_slug=exercises["Стол"].slug, repetitions=10),
                WorkoutExercise(exercise_slug=exercises["Наклоны вперед"].slug, repetitions=10),
                WorkoutExercise(exercise_slug=exercises["Качающиеся планки"].slug, repetitions=20),
            ],
        ),
        WorkoutSection(
            title="Комплекс", 
            protocol=WorkoutProtocol(title="5 раундов"),
            exercises=[
                WorkoutExercise(exercise_slug=exercises["Становая на одной"].slug, repetitions=24),
                WorkoutExercise(exercise_slug=exercises["Приседания"].slug, repetitions=18),
                WorkoutExercise(exercise_slug=exercises["С груди над головой"].slug, repetitions=12),
            ],
        ),
    ]
)

print(f"create workout: {workout}")
workout_repo.add(workout)


# commit all changes
session.commit()
