import { DateFormatter } from "@internationalized/date";

// todo: move locale to config
const isoDateFormatter = new DateFormatter("ru-RU", {
  day: "numeric",
  month: "long",
});

export function formatIsoDate(isoDate: string) {
  const date = new Date(isoDate);

  let formattedDate = isoDateFormatter.format(date);

  if (date.getFullYear() !== new Date().getFullYear()) {
    formattedDate += ` ${date.getFullYear()}`;
  }

  return formattedDate;
}
