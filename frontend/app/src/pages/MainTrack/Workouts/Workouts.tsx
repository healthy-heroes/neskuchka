import { Heading, Text, View } from "@adobe/react-spectrum";

import { TrackWorkouts } from "#api/methods";
import { formatIsoDate } from "../../../utils/dates";

type WorkoutsProps = {
  /** TrackWorkouts stores a list of workouts and a dictionary with exercise details */
  trackWorkouts: TrackWorkouts;
};

export function Workouts({ trackWorkouts }: WorkoutsProps) {
  return trackWorkouts.Workouts.map((workout) => (
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
  ));
}
