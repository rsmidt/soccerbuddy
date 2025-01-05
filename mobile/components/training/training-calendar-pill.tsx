import {
  GetMyTeamHomeResponse_Nomination,
  GetMyTeamHomeResponse_Training,
} from "@/api/soccerbuddy/team/v1/team_service_pb";
import { StyleProp, StyleSheet, View, ViewStyle } from "react-native";
import { Text, TouchableRipple, useTheme } from "react-native-paper";
import i18n from "@/components/i18n";
import { useRouter } from "expo-router";
import React, { memo } from "react";
import { padNumber } from "../date";
import { useGetMeQuery } from "@/components/account/account-api";
import { selectHasEditAllowance } from "@/components/team/team-api";
import { TrainingCalendarPillPlayerResponse } from "@/components/training/training-calendar-pill-player-response";

type TrainingCalendarPillProps = {
  teamId: string;
  training: GetMyTeamHomeResponse_Training;
  style?: StyleProp<ViewStyle>;
};

const TrainingCalendarPillImpl = ({
  teamId,
  training,
  style,
}: TrainingCalendarPillProps) => {
  const theme = useTheme();
  const router = useRouter();
  const { hasEditAllowance, isLoading: isGetMeLoading } = useGetMeQuery(
    {},
    {
      selectFromResult: ({ data, ...rest }) => ({
        ...rest,
        hasEditAllowance: selectHasEditAllowance(data, teamId),
      }),
    },
  );
  if (isGetMeLoading) {
    return null;
  }
  const playerAcknowledgments = selectPlayerResponse(
    training.nominations?.players,
  );
  const staffAcknowledgments = selectStaffResponse(training.nominations?.staff);

  return (
    <View
      key={training.id}
      style={[
        styles.datePill,
        { backgroundColor: theme.colors.primaryContainer },
        style,
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
            {selectGatheringPoint(training)}
            {getHourRange(training)}
          </Text>
          {hasEditAllowance && (
            <>
              {playerAcknowledgments && (
                <TrainingCalendarPillPlayerResponse
                  {...playerAcknowledgments}
                />
              )}
              {staffAcknowledgments && (
                <Text
                  variant="labelMedium"
                  style={{ color: theme.colors.onPrimaryContainer }}
                >
                  {staffAcknowledgments}
                </Text>
              )}
            </>
          )}
        </View>
      </TouchableRipple>
    </View>
  );
};

const styles = StyleSheet.create({
  datePill: {
    borderRadius: 12,
  },
  datePillRipple: {
    padding: 10,
    borderRadius: 12,
  },
});

function getHourRange(training: GetMyTeamHomeResponse_Training): string {
  const { scheduledAt, endsAt } = training;
  return `${padNumber(scheduledAt?.hours!)}:${padNumber(scheduledAt?.minutes!)} – ${padNumber(endsAt?.hours!)}:${padNumber(endsAt?.minutes!)}`.toLowerCase();
}

function selectGatheringPoint(
  training: GetMyTeamHomeResponse_Training,
): string | undefined {
  if (training.gatheringPoint === undefined) return undefined;
  const { gatheringUntil } = training.gatheringPoint;
  return `${padNumber(gatheringUntil?.hours!)}:${padNumber(gatheringUntil?.minutes!)} | `;
}

function selectStaffResponse(
  nominations: GetMyTeamHomeResponse_Nomination[] | undefined,
): string {
  if (nominations === undefined) return "";
  return nominations
    .filter((nomination) => nomination.response.case === "accepted")
    .map((nomination) => nomination.personName.split(" ")[0])
    .join(", ");
}

function selectPlayerResponse(
  nominations: GetMyTeamHomeResponse_Nomination[] | undefined,
) {
  if (nominations === undefined) return undefined;
  return nominations.reduce(
    (previousValue, currentValue) => {
      if (currentValue.response.case === "accepted") {
        previousValue.accepted++;
      } else if (currentValue.response.case === "declined") {
        previousValue.declined++;
      } else if (currentValue.response.case === "tentative") {
        previousValue.tentative++;
      } else {
        previousValue.unknown++;
      }
      return previousValue;
    },
    {
      accepted: 0,
      declined: 0,
      tentative: 0,
      unknown: 0,
    },
  );
}

export const TrainingCalendarPill = memo(TrainingCalendarPillImpl);
