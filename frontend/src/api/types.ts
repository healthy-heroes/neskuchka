import { Workout } from '@/types/domain';

export interface TrackWorkouts {
	Workouts: Array<Workout>;
}

export interface TrackWorkout {
	Workout: Workout;
}
