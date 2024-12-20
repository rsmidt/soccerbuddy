import { ScrollView, StyleSheet, View } from "react-native";
import { Stack, useLocalSearchParams } from "expo-router";
import i18n, { locale } from "@/components/i18n";
import React from "react";
import { MaterialCommunityIcons, MaterialIcons } from "@expo/vector-icons";
import { Text } from "react-native-paper";
import { useGetMyTeamHomeQuery } from "@/components/team/team-api";
import {
  GetMyTeamHomeResponse,
  GetMyTeamHomeResponse_Training,
} from "@/api/soccerbuddy/team/v1/team_service_pb";
import { FormRow } from "@/components/form/form-row";
import { pbToDateTime } from "@/components/proto";
import { DateTime } from "luxon";

function selectTraining(
  data: GetMyTeamHomeResponse | undefined,
  trainingId: string,
): GetMyTeamHomeResponse_Training | undefined {
  if (!data) return undefined;
  return data.trainings.find((training) => training.id === trainingId);
}

export default function TrainingDetail() {
  const { team, trainingId } = useLocalSearchParams<{
    team: string;
    trainingId: string;
  }>();
  const { training, isLoading, teamName } = useGetMyTeamHomeQuery(
    { teamId: team },
    {
      selectFromResult: ({ data, ...rest }) => ({
        ...rest,
        teamName: data?.teamName,
        training: selectTraining(data, trainingId),
      }),
    },
  );

  if (isLoading || !training) {
    return null;
  }

  return (
    <ScrollView style={styles.container}>
      <Stack.Screen
        options={{
          headerTitle: i18n.t("app.teams.training.detail.title"),
        }}
      />
      <View style={styles.list}>
        <FormRow style={{ paddingVertical: 8 }}>
          <FormRow.Icon>
            <MaterialIcons name="people-outline" size={24} />
          </FormRow.Icon>
          <FormRow.Controls>
            <Text variant="bodyLarge">{teamName}</Text>
          </FormRow.Controls>
        </FormRow>
        <FormRow style={{ paddingVertical: 8 }}>
          <FormRow.Icon>
            <MaterialIcons name="access-time" size={24} />
          </FormRow.Icon>
          <FormRow.Controls direction="row">
            <Text variant="bodyLarge">{selectTimeframe(training)}</Text>
          </FormRow.Controls>
        </FormRow>
        <FormRow style={{ paddingVertical: 8 }}>
          <FormRow.Icon>
            <MaterialCommunityIcons name="flag-outline" size={24} />
          </FormRow.Icon>
          <FormRow.Controls direction="row">
            <Text variant="bodyLarge">{selectDeadline(training)}</Text>
          </FormRow.Controls>
        </FormRow>
        <FormRow style={{ paddingVertical: 8 }}>
          <FormRow.Icon>
            <MaterialCommunityIcons name="text" size={24} />
          </FormRow.Icon>
          <FormRow.Controls>
            <Text variant="bodyLarge">{selectDescription(training)}</Text>
          </FormRow.Controls>
        </FormRow>
        <FormRow style={{ paddingVertical: 8 }}>
          <FormRow.Icon>
            <MaterialCommunityIcons name="map-marker-outline" size={24} />
          </FormRow.Icon>
          <FormRow.Controls>
            <Text variant="bodyLarge">{selectLocation(training)}</Text>
          </FormRow.Controls>
        </FormRow>
        <FormRow style={{ paddingVertical: 8 }}>
          <FormRow.Icon>
            <MaterialCommunityIcons name="soccer-field" size={24} />
          </FormRow.Icon>
          <FormRow.Controls>
            <Text variant="bodyLarge">{selectFieldType(training)}</Text>
          </FormRow.Controls>
        </FormRow>
        <FormRow style={{ paddingTop: 8 }}>
          <FormRow.Icon>
            <MaterialCommunityIcons name="map-marker-path" size={24} />
          </FormRow.Icon>
          <FormRow.Controls>
            <Text variant="bodyLarge">{selectGatherPoint(training)}</Text>
          </FormRow.Controls>
        </FormRow>
      </View>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  list: {
    marginTop: 16,
    flexDirection: "column",
    gap: 16,
  },
});

function selectTimeframe(training: GetMyTeamHomeResponse_Training): string {
  const start = pbToDateTime(training.scheduledAt);
  const end = pbToDateTime(training.endsAt);

  if (!start || !end) {
    return i18n.t("app.teams.training.detail.time.unknown");
  }

  const startDt = DateTime.fromJSDate(start).setLocale(locale);
  const endDt = DateTime.fromJSDate(end).setLocale(locale);

  if (startDt.hasSame(endDt, "day")) {
    return i18n.t("app.teams.training.detail.time.same-day", {
      date: startDt.toLocaleString(DateTime.DATE_FULL),
      startTime: startDt.toLocaleString(DateTime.TIME_SIMPLE),
      endTime: endDt.toLocaleString(DateTime.TIME_SIMPLE),
    });
  }

  return i18n.t("app.teams.training.detail.time.different-days", {
    startDate: startDt.toLocaleString(DateTime.DATE_FULL),
    startTime: startDt.toLocaleString(DateTime.TIME_SIMPLE),
    endDate: endDt.toLocaleString(DateTime.DATE_FULL),
    endTime: endDt.toLocaleString(DateTime.TIME_SIMPLE),
  });
}

function selectDeadline(training: GetMyTeamHomeResponse_Training): string {
  const deadline = pbToDateTime(training.acknowledgmentSettings?.deadline);
  if (!deadline) return i18n.t("app.teams.training.detail.acknowledgment.none");
  return i18n.t("app.teams.training.detail.acknowledgment.until", {
    time: DateTime.fromJSDate(deadline)?.setLocale(locale)?.toLocaleString(),
  });
}

function selectGatherPoint(training: GetMyTeamHomeResponse_Training): string {
  const gatherLocation = training.gatheringPoint?.location;
  const gatherUntil = (() => {
    const date = pbToDateTime(training.gatheringPoint?.gatheringUntil);
    if (!date) return undefined;
    return DateTime.fromJSDate(date)?.setLocale(locale)?.toLocaleString();
  })();

  if (!gatherLocation && !gatherUntil) {
    return i18n.t("app.teams.training.detail.gathering.none");
  }

  if (gatherLocation && gatherUntil) {
    return i18n.t("app.teams.training.detail.gathering.time-location", {
      time: gatherUntil,
      location: gatherLocation,
    });
  }

  if (gatherUntil) {
    return i18n.t("app.teams.training.detail.gathering.time", {
      time: gatherUntil,
    });
  }

  // Only gatherLocation is defined
  return i18n.t("app.teams.training.detail.gathering.location", {
    location: gatherLocation,
  });
}

function selectDescription(training: GetMyTeamHomeResponse_Training): string {
  if (!training.description) {
    return i18n.t("app.teams.training.detail.description.empty");
  }
  return training.description;
}

function selectFieldType(training: GetMyTeamHomeResponse_Training): string {
  if (!training.fieldType) {
    return i18n.t("app.teams.training.detail.field-type.empty");
  }
  return training.fieldType;
}

function selectLocation(training: GetMyTeamHomeResponse_Training): string {
  if (!training.location) {
    return i18n.t("app.teams.training.detail.location.empty");
  }
  return training.location;
}
