import { useQuery } from "@tanstack/react-query";
import { Heading, ProgressCircle, Text, View } from "@adobe/react-spectrum";

import { pageProps } from "../../constants";
import { getMainTrackWorkouts } from "../../../services/api";

export function Workouts() {
  const {
    data: workouts,
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

      {workouts?.map((workout) => (
        <View
          key={workout.id}
          borderWidth="thin"
          borderColor="gray-200"
          borderRadius="medium"
          paddingX="size-200"
          paddingBottom="size-200"
          marginBottom="size-200"
        >
          <Heading level={2}>{workout.date}</Heading>

          {workout.sections.map((section, index) => {
            const sectionKey = `${workout.id}-${index}`;

            return (
              <View key={sectionKey}>
                <Heading level={3}>{section.title}</Heading>
                <Text>{section.protocol.title}</Text>
                <br />
                <Text>{section.protocol.description}</Text>
                <br />

                {section.exercises.map((exercise, index) => {
                  const exerciseKey = `${sectionKey}-${index}`;

                  return (
                    <Text key={exerciseKey}>
                      {exercise.exercise_slug}
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
