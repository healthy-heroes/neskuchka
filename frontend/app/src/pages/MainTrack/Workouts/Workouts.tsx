import { Heading, Text, View } from "@adobe/react-spectrum";

import { pageProps } from "../../constants";
import { TrackWorkouts } from "../../../services/api";
import { formatIsoDate } from "../../../utils/dates";

type WorkoutsProps = {
  // Track workouts data containing workout information and dictionary of exercises.
  trackWorkouts: TrackWorkouts;
};

export function Workouts({ trackWorkouts }: WorkoutsProps) {
  return (
    <View
      paddingX="size-400"
      paddingY="size-400"
      maxWidth={pageProps.maxWidth}
      marginX="auto"
    >
      {trackWorkouts.Workouts.map((workout) => (
        <View
          key={workout.ID}
          borderWidth="thin"
          borderColor="gray-200"
          borderRadius="medium"
          paddingX="size-200"
          paddingBottom="size-200"
          marginBottom="size-200"
        >
          <Heading level={2}>{formatIsoDate(workout.Date)}</Heading>

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
