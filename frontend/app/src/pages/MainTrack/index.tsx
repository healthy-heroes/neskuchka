import { ProgressCircle, Text } from "@adobe/react-spectrum";
import { useQuery } from "@tanstack/react-query";

import { getMainTrackWorkouts } from "../../services/api";

import { TrackHeader } from "./TrackHeader";
import { Workouts } from "./Workouts";

export function MainTrackPage() {
  const {
    data: trackWorkouts,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["main-track-workouts"],
    queryFn: getMainTrackWorkouts,
  });

  return (
    <main>
      <TrackHeader />

      {error && <Text>Error: {error.message}</Text>}

      {isLoading && (
        <ProgressCircle aria-label="Loading…" isIndeterminate={true} />
      )}

      {trackWorkouts && <Workouts trackWorkouts={trackWorkouts} />}
    </main>
  );
}
