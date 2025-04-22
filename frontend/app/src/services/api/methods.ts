import { Exercise, Workout } from "#types/domain";
import { HttpError } from "./httpErrors";

const API_URL = "http://localhost:8080/api/v1";

export interface TrackWorkouts {
  Workouts: Array<Workout>;
  Exercises: Record<string, Exercise>;
}

export async function getMainTrackWorkouts(): Promise<TrackWorkouts> {
  const response = await fetch(`${API_URL}/tracks/main/last_workouts`);

  if (!response.ok) {
    throw new HttpError(response.status, await response.json());
  }

  return response.json();
}
