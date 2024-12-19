import { MessageInitShape } from "@bufbuild/protobuf";
import { DateTime, DateTimeSchema } from "@/api/google/type/datetime_pb";

export function dateTimeToPb(
  date: Date,
): MessageInitShape<typeof DateTimeSchema> {
  return {
    day: date.getDate(),
    month: date.getMonth() + 1,
    year: date.getFullYear(),
    hours: date.getHours(),
    minutes: date.getMinutes(),
    seconds: date.getSeconds(),
  };
}

export function pbToDateTime(dateTime: DateTime | undefined): Date | undefined {
  if (!dateTime) return undefined;

  return new Date(
    dateTime.year,
    dateTime.month - 1,
    dateTime.day,
    dateTime.hours,
    dateTime.minutes,
    dateTime.seconds,
  );
}
