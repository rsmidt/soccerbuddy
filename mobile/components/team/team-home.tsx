import { useGetMyTeamHomeQuery } from "@/components/team/team-api";
import { StyleProp, StyleSheet, View, ViewStyle } from "react-native";
import { Text, TouchableRipple, useTheme } from "react-native-paper";
import {
  GetMyTeamHomeResponse,
  GetMyTeamHomeResponse_Training,
} from "@/api/soccerbuddy/team/v1/team_service_pb";
import { pbToDateTime } from "@/components/proto";
import { DateTime } from "@/api/google/type/datetime_pb";
import i18n from "@/components/i18n";
import { useRouter } from "expo-router";
import { useMemo } from "react";

type TeamHomeProps = { teamId: string; style: StyleProp<ViewStyle> };

export function TeamHome({ teamId, style }: TeamHomeProps) {
  const theme = useTheme();
  const { data, isLoading } = useGetMyTeamHomeQuery(
    { teamId },
    {
      selectFromResult: ({ data, ...rest }) => ({
        data: selectGroupedUpcomingTrainings(data),
        ...rest,
      }),
    },
  );
  const router = useRouter();

  const groupedUpcomingTrainings = useMemo(
    () =>
      data?.trainings?.reduce(
        (acc, training) => {
          const { scheduledAt } = training;
          const dateKey = generateDateKey(scheduledAt!);
          if (!acc[dateKey]) {
            acc[dateKey] = [];
          }
          acc[dateKey].push(training);
          return acc;
        },
        {} as Record<string, GetMyTeamHomeResponse_Training[]>,
      ) ?? {},
    [data?.trainings],
  );

  const dates = useMemo(
    () => Object.keys(groupedUpcomingTrainings).sort(),
    [groupedUpcomingTrainings],
  );

  const allDates = useMemo(() => {
    // Get the range of dates
    const firstDate = new Date(dates[0]);
    const lastDate = new Date(dates[dates.length - 1]);

    // Generate all months between first and last date
    const allDates: { type: "year" | "month" | "day"; date: string }[] = [];
    const currentDate = new Date(firstDate);

    while (currentDate <= lastDate) {
      const year = currentDate.getFullYear();
      const month = currentDate.getMonth();
      const monthKey = `${year}-${padNumber(month + 1)}`;

      // Add month divider
      allDates.push({ type: "month", date: monthKey });

      // Add all days that have trainings for this month
      dates
        .filter((date) => date.startsWith(monthKey))
        .forEach((date) => {
          allDates.push({ type: "day", date });
        });

      // Move to next month
      currentDate.setMonth(currentDate.getMonth() + 1);
    }
    return allDates;
  }, [dates]);

  if (isLoading || !data) {
    return null;
  }

  const now = new Date();
  const todayDateKey = `${now.getFullYear()}-${padNumber(now.getMonth() + 1)}-${padNumber(now.getDate())}`;

  return (
    <View style={[styles.container, style]}>
      {allDates.map(({ type, date }) => {
        if (type === "month") {
          return (
            <Text
              key={`month-${date}`}
              variant="bodyLarge"
              style={styles.monthDivider}
            >
              {new Date(date + "-01").toLocaleDateString(undefined, {
                month: "long",
                year: "numeric",
              })}
            </Text>
          );
        }

        // Existing day rendering code
        const trainings = groupedUpcomingTrainings[date];
        const isToday = todayDateKey === date;

        return (
          <View style={styles.dayContainer} key={`day-${date}`}>
            <View style={styles.dayIconContainer}>
              <Text
                variant="labelSmall"
                style={[styles.weekDay, isToday && { fontWeight: "bold" }]}
              >
                {getWeekDay(date)}
              </Text>
              <View
                style={[
                  styles.dayIcon,
                  isToday && { backgroundColor: theme.colors.primary },
                ]}
              >
                <Text
                  variant="bodyMedium"
                  style={[isToday && { color: theme.colors.onPrimary }]}
                >
                  {date.split("-")[2]}
                </Text>
              </View>
            </View>
            <View style={styles.trainingsContainer}>
              {trainings.map((training) => (
                <View
                  key={training.id}
                  style={[
                    styles.datePill,
                    { backgroundColor: theme.colors.primaryContainer },
                  ]}
                >
                  <TouchableRipple
                    style={styles.datePillRipple}
                    borderless
                    onPress={() =>
                      router.navigate({
                        pathname: "/teams/[team]/training/[trainingId]/detail",
                        params: {
                          team: teamId,
                          trainingId: training.id,
                        },
                      })
                    }
                  >
                    <View>
                      <Text
                        variant="labelMedium"
                        style={{ color: theme.colors.onPrimaryContainer }}
                      >
                        {i18n.t("app.teams.home.training")}
                      </Text>
                      <Text
                        variant="labelMedium"
                        style={{ color: theme.colors.onPrimaryContainer }}
                      >
                        {getHourRange(training)}
                      </Text>
                    </View>
                  </TouchableRipple>
                </View>
              ))}
            </View>
          </View>
        );
      })}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    marginHorizontal: 16,
    gap: 16,
  },
  dayContainer: {
    flexDirection: "row",
    gap: 12,
  },
  dayIconContainer: {
    alignItems: "center",
  },
  datePill: {
    borderRadius: 12,
  },
  datePillRipple: {
    padding: 10,
    borderRadius: 12,
  },
  dayIcon: {
    width: 24,
    height: 24,
    justifyContent: "center",
    alignItems: "center",
    borderRadius: 999,
  },
  weekDay: {
    textAlign: "center",
  },
  yearDivider: {
    marginTop: 16,
    fontWeight: "bold",
  },
  monthDivider: {
    marginLeft: 32,
  },
  trainingsContainer: {
    flex: 1,
    gap: 8,
  },
});

function selectGroupedUpcomingTrainings(
  data: GetMyTeamHomeResponse | undefined,
): GetMyTeamHomeResponse | undefined {
  if (!data) return undefined;
  const now = new Date();

  return {
    ...data,
    trainings: data.trainings
      .filter((training) => {
        const date = pbToDateTime(training.scheduledAt);
        return date ? date > now : false;
      })
      .sort(
        (a, b) =>
          (pbToDateTime(a.scheduledAt)?.getTime() ?? 0) -
          (pbToDateTime(b.scheduledAt)?.getTime() ?? 0),
      ),
  };
}

function padNumber(num: number, length: number = 2): string {
  return num.toString().padStart(length, "0");
}

/**
 * Generates a standardized date key in the format 'YYYY-MM-DD'.
 */
function generateDateKey(scheduledAt: DateTime): string {
  const year =
    scheduledAt.year && scheduledAt.year !== 0 ? scheduledAt.year : "0000";
  const month = padNumber(scheduledAt.month);
  const day = padNumber(scheduledAt.day);
  return `${year}-${month}-${day}`;
}

function getWeekDay(dateKey: string): string {
  const date = new Date(dateKey);
  return date.toLocaleDateString(undefined, {
    weekday: "short",
  });
}

function getHourRange(training: GetMyTeamHomeResponse_Training): string {
  const { scheduledAt, endsAt } = training;
  return `${padNumber(scheduledAt?.hours!)}:${padNumber(scheduledAt?.minutes!)} – ${padNumber(endsAt?.hours!)}:${padNumber(endsAt?.minutes!)}`.toLowerCase();
}
