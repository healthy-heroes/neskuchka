from app.domain.exercise import Exercise, ExerciseRepository
from app.domain.track import Track, TrackRepository
from app.domain.workout import Workout, WorkoutRepository
from app.tests.fixtures.domain import (
    create_exercise,
    create_track,
    create_workout,
    create_workout_exercise,
    create_workout_section,
)


def exercise_for_assert(exercise: Exercise) -> dict:
    exercise_assert = exercise.model_dump()
    exercise_assert["slug"] = str(exercise.slug)
    return exercise_assert


def track_for_assert(track: Track) -> dict:
    track_assert = track.model_dump()
    track_assert["id"] = str(track.id)
    track_assert["owner_id"] = str(track.owner_id)
    return track_assert


def workout_for_assert(workout: Workout) -> dict:
    workout_assert = workout.model_dump()
    workout_assert["id"] = str(workout.id)
    workout_assert["date"] = str(workout.date)
    workout_assert["track_id"] = str(workout.track_id)
    return workout_assert


def test_get_main_track(client, track_repo: TrackRepository):
    track = create_track()
    track_repo.add(track)

    response = client.get("/api/v1/tracks/main")

    assert response.status_code == 200
    assert response.json() == track_for_assert(track)


def test_get_non_existent_main_track(client, track_repo: TrackRepository):
    response = client.get("/api/v1/tracks/main")

    assert response.status_code == 404


def test_get_workouts_for_main_track(
    client,
    track_repo: TrackRepository,
    workout_repo: WorkoutRepository,
    exercise_repo: ExerciseRepository,
):
    track = create_track()
    track_repo.add(track)

    all_exercises = []

    for _ in range(5):
        exercise = create_exercise()
        exercise_repo.add(exercise)
        all_exercises.append(exercise)

    exercises = all_exercises[:3]

    workouts = []
    for _ in range(12):
        workout = create_workout(
            track=track,
            sections=[
                create_workout_section(
                    exercises=[
                        create_workout_exercise(exercise) for exercise in exercises
                    ]
                ),
            ],
        )
        workout_repo.add(workout)
        workouts.append(workout)

    response = client.get("/api/v1/tracks/main/last_workouts")

    assert response.status_code == 200
    assert response.json()["workouts"] == [
        workout_for_assert(workout) for workout in workouts[:10]
    ], "Должен возвращать последние 10 тренировок"
    assert response.json()["exercises"] == {
        str(ex.slug): exercise_for_assert(ex) for ex in exercises
    }, "Должен возвращать упражнения из воркаутов"


def test_get_workouts_for_non_existent_track(
    client,
    track_repo: TrackRepository,
    workout_repo: WorkoutRepository,
    exercise_repo: ExerciseRepository,
):
    response = client.get("/api/v1/tracks/main/last_workouts")

    assert response.status_code == 404
