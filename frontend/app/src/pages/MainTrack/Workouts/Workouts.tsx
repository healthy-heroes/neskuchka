import { useQuery } from "@tanstack/react-query";
import { Heading, ProgressCircle, Text, View } from "@adobe/react-spectrum";

import { pageProps } from "../../constants";
import { getMainTrackWorkouts } from "../../../services/api";
import { WorkoutDate } from "./WorkoutDate";

export function Workouts() {
  const {
    data: trackWorkouts,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["main-track-workouts"],
    queryFn: getMainTrackWorkouts,
  });

  return (
    <View
      paddingX="size-400"
      paddingY="size-400"
      maxWidth={pageProps.maxWidth}
      marginX="auto"
    >
      {isLoading && (
        <ProgressCircle aria-label="Loading…" isIndeterminate={true} />
      )}
      {error && <Text>Error: {error.message}</Text>}

      {trackWorkouts?.Workouts.map((workout) => (
        <View
          key={workout.ID}
          borderWidth="thin"
          borderColor="gray-200"
          borderRadius="medium"
          paddingX="size-200"
          paddingBottom="size-200"
          marginBottom="size-200"
        >
          <Heading level={2}>{WorkoutDate(workout.Date)}</Heading>

          {workout.Sections.map((section, index) => {
            const sectionKey = `${workout.ID}-${index}`;

            return (
              <View key={sectionKey}>
                <Heading level={3}>{section.Title}</Heading>
                <Text>{section.Protocol.Title}</Text>
                <br />
                <Text>{section.Protocol.Description}</Text>
                <br />

                {section.Exercises.map((exercise, index) => {
                  const exerciseKey = `${sectionKey}-${index}`;

                  return (
                    <Text key={exerciseKey}>
                      {trackWorkouts.Exercises[exercise.ExerciseSlug]?.Name ||
                        exercise.ExerciseSlug}
                      <br />
                    </Text>
                  );
                })}
              </View>
            );
          })}
        </View>
      ))}
    </View>
  );
}
