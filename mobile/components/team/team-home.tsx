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

  if (isLoading || !data) {
    return null;
  }

  const now = new Date();
  const todayDateKey = `${now.getFullYear()}-${padNumber(now.getMonth() + 1)}-${padNumber(now.getDate())}`;

  const groupedUpcomingTrainings = data.trainings.reduce(
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
  );

  return (
    <View style={[styles.container, style]}>
      {Object.entries(groupedUpcomingTrainings).map(([dateKey, trainings]) => {
        const isToday = todayDateKey === dateKey;

        return (
          <View style={styles.dayContainer} key={dateKey}>
            <View style={styles.dayIconContainer}>
              <Text
                variant="labelSmall"
                style={[styles.weekDay, isToday && { fontWeight: "bold" }]}
              >
                {getWeekDay(dateKey)}
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
                  {dateKey.split("-")[2]}
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
