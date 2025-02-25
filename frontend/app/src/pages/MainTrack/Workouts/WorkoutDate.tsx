import { DateFormatter } from "@internationalized/date";

// todo: move locale to config
const dateFormatter = new DateFormatter("ru-RU", {
  day: "numeric",
  month: "long",
});

export function WorkoutDate(isoDate: string) {
  const date = new Date(isoDate);

  let formattedDate = dateFormatter.format(date);

  if (date.getFullYear() !== new Date().getFullYear()) {
    formattedDate += ` ${date.getFullYear()}`;
  }

  return formattedDate;
}
