import { ProgressCircle, Text, View } from "@adobe/react-spectrum";
import { useQuery } from "@tanstack/react-query";

import { getMainTrackWorkouts } from "#api/methods";
import { pageProps } from "../constants";

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

      <View
        paddingX="size-400"
        paddingY="size-400"
        maxWidth={pageProps.maxWidth}
        marginX="auto"
      >
        {error && <Text>Error: {error.message}</Text>}

        {isLoading && (
          <ProgressCircle aria-label="Loading…" isIndeterminate={true} />
        )}

        {trackWorkouts && <Workouts trackWorkouts={trackWorkouts} />}
      </View>
    </main>
  );
}
