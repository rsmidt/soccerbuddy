import { ScrollView, StyleSheet, View } from "react-native";
import { Stack, useLocalSearchParams } from "expo-router";
import i18n from "@/components/i18n";
import { Controller, FormProvider, useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import React from "react";
import { DatePickerButton } from "@/components/form/date-picker-button";
import { MaterialCommunityIcons, MaterialIcons } from "@expo/vector-icons";
import { Ruler } from "@/components/ruler";
import { Text } from "react-native-paper";
import { InlineMaterialInput } from "@/components/form/inline-material-input";
import { DateTime } from "luxon";
import { useScheduleTrainingMutation } from "@/components/team/team-api";
import {
  ScheduleTrainingRequest,
  ScheduleTrainingRequest_AcknowledgementSettingsSchema,
  ScheduleTrainingRequest_GatheringPointSchema,
} from "@/api/soccerbuddy/team/v1/team_service_pb";
import { DateTimeSchema } from "@/api/google/type/datetime_pb";
import { MessageInitShape } from "@bufbuild/protobuf";
import { extractBadRequestDetail } from "@/components/connect-base-query";
import Toast from "react-native-toast-message";

const gatheringPointSchema = z
  .object({
    location: z.string().optional(),
    gatheringUntil: z
      .date({
        required_error: "Gathering time is required",
        invalid_type_error: "Invalid gathering time",
      })
      .optional(),
  })
  .refine(
    (data) => {
      const isAnyFieldFilled = !data.location || !data.gatheringUntil;
      if (isAnyFieldFilled) {
        return !data.location && !data.gatheringUntil;
      }
      return true;
    },
    {
      message:
        "Both location and gatheringUntil are required if gatheringPoint is provided",
      path: ["gatheringPoint"],
    },
  );

const acknowledgementSettingsSchema = z.object({
  deadline: z
    .date({
      required_error: "Deadline is required",
      invalid_type_error: "Invalid deadline",
    })
    .optional(),
});

const ratingSettingsSchema = z.object({
  policy: z
    .enum(["UNSPECIFIED", "FORBIDDEN", "ALLOWED", "REQUIRED"], {
      required_error: "Rating policy is required",
    })
    .optional(),
});

export const scheduleTrainingSchema = z
  .object({
    scheduledAt: z.date({
      required_error: "Start time is required",
      invalid_type_error: "Invalid start time",
    }),
    endsAt: z.date({
      required_error: "End time is required",
      invalid_type_error: "Invalid end time",
    }),
    location: z.string().optional(),
    fieldType: z.string().optional(),
    description: z.string().optional(),
    gatheringPoint: gatheringPointSchema.optional(),
    acknowledgmentSettings: acknowledgementSettingsSchema.optional(),
    ratingSettings: ratingSettingsSchema.optional(),
  })
  .refine((data) => data.endsAt > data.scheduledAt, {
    message: "End time must be after start time",
    path: ["endsAt"],
  });

type ScheduleTrainingForm = z.infer<typeof scheduleTrainingSchema>;

export default function ScheduleTraining() {
  const { team } = useLocalSearchParams<{ team: string }>();
  const [scheduleTraining, { isLoading }] = useScheduleTrainingMutation();

  const now = DateTime.now();
  const form = useForm<ScheduleTrainingForm>({
    resolver: zodResolver(scheduleTrainingSchema),
    defaultValues: {
      scheduledAt: now.toJSDate(),
      endsAt: now.plus({ hour: 1 }).toJSDate(),
    },
  });
  const { handleSubmit, control } = form;

  const onSubmit = handleSubmit(async (data) => {
    const cleanedData = {
      ...data,
      gatheringPoint:
        data.gatheringPoint?.location || data.gatheringPoint?.gatheringUntil
          ? data.gatheringPoint
          : undefined,
      acknowledgmentSettings: data.acknowledgmentSettings?.deadline
        ? data.acknowledgmentSettings
        : undefined,
      ratingSettings: data.ratingSettings?.policy
        ? data.ratingSettings
        : undefined,
    };

    try {
      await scheduleTraining({
        teamId: team,
        scheduledAt: dateTimeToPb(cleanedData.scheduledAt),
        endsAt: dateTimeToPb(cleanedData.endsAt),
        location: cleanedData.location,
        fieldType: cleanedData.fieldType,
        description: cleanedData.description,
        gatheringPoint: maybeGatheringPointToPb(
          cleanedData.gatheringPoint as any,
        ),
        acknowledgmentSettings: maybeAcknowledgmentSettingsToPb(
          cleanedData.acknowledgmentSettings as any,
        ),
      } as ScheduleTrainingRequest).unwrap();

      Toast.show({
        type: "success",
        text1: i18n.t("app.team.schedule-training.success"),
        position: "bottom",
      });
    } catch (error) {
      // TODO: proper error handling...
      const badRequestDetail = extractBadRequestDetail(error);
      if (badRequestDetail) {
        const errors = badRequestDetail.fieldViolations
          .map((violation) => `${violation.field} => ${violation.description}`)
          .join("\n");
        Toast.show({
          type: "error",
          text1: "Input is wrong: " + errors,
          position: "bottom",
        });
      }
    }
  });

  return (
    <FormProvider {...form}>
      <ScrollView style={styles.container}>
        <Stack.Screen
          options={{
            headerTitle: i18n.t("app.teams.schedule-training.title"),
            headerRight: () => (
              <MaterialCommunityIcons
                name={isLoading ? "loading" : "send"}
                size={24}
                onPress={onSubmit}
              />
            ),
          }}
        />
        <View style={styles.form}>
          <View style={styles.formRow}>
            <View style={styles.iconContainer}>
              <MaterialIcons name="people-outline" size={24} />
            </View>
            <Text>{team}</Text>
          </View>
          <Ruler />
          <View style={styles.formRow}>
            <View style={styles.iconContainer}>
              <MaterialIcons name="access-time" size={24} />
            </View>
            <View style={styles.dateRow}>
              <DatePickerButton
                control={control}
                type="date"
                name="scheduledAt"
                label={i18n.t("app.teams.schedule-training.scheduled-at.label")}
              />
              <DatePickerButton
                style={{ maxWidth: 50 }}
                control={control}
                type="time"
                name="scheduledAt"
                label={i18n.t("app.teams.schedule-training.scheduled-at.label")}
              />
            </View>
          </View>
          <View style={styles.formRow}>
            <View style={styles.iconContainer} />
            <View style={styles.dateRow}>
              <DatePickerButton
                control={control}
                type="date"
                name="endsAt"
                label={i18n.t("app.teams.schedule-training.ends-at.label")}
              />
              <DatePickerButton
                style={{ maxWidth: 50 }}
                control={control}
                type="time"
                name="endsAt"
                label={i18n.t("app.teams.schedule-training.ends-at.label")}
              />
            </View>
          </View>
          <Ruler />
          <View style={[styles.formRow, { paddingVertical: 8 }]}>
            <View style={styles.iconContainer}>
              <MaterialCommunityIcons name="text" size={24} />
            </View>
            <Controller
              control={control}
              render={({ field: { onChange, onBlur, value } }) => (
                <InlineMaterialInput
                  multiline
                  onChangeText={onChange}
                  onBlur={onBlur}
                  value={value}
                  placeholder={i18n.t(
                    "app.teams.schedule-training.description.label",
                  )}
                />
              )}
              name="description"
            />
          </View>
          <Ruler />
          <View style={[styles.formRow, { paddingVertical: 8 }]}>
            <View style={styles.iconContainer}>
              <MaterialCommunityIcons name="map-marker-outline" size={24} />
            </View>
            <Controller
              control={control}
              render={({ field: { onChange, onBlur, value } }) => (
                <InlineMaterialInput
                  onChangeText={onChange}
                  onBlur={onBlur}
                  value={value}
                  placeholder={i18n.t(
                    "app.teams.schedule-training.location.label",
                  )}
                />
              )}
              name="location"
            />
          </View>
          <Ruler />
          <View style={[styles.formRow, { paddingVertical: 8 }]}>
            <View style={styles.iconContainer}>
              <MaterialCommunityIcons name="soccer-field" size={24} />
            </View>
            <Controller
              control={control}
              render={({ field: { onChange, onBlur, value } }) => (
                <InlineMaterialInput
                  onChangeText={onChange}
                  onBlur={onBlur}
                  value={value}
                  placeholder={i18n.t(
                    "app.teams.schedule-training.field-type.label",
                  )}
                />
              )}
              name="fieldType"
            />
          </View>
          <Ruler />
          <View style={[styles.formRow, { paddingTop: 8, paddingBottom: 4 }]}>
            <View style={styles.iconContainer}>
              <MaterialCommunityIcons name="map-marker-path" size={24} />
            </View>
            <Controller
              control={control}
              render={({ field: { onChange, onBlur, value } }) => (
                <InlineMaterialInput
                  onChangeText={onChange}
                  onBlur={onBlur}
                  value={value}
                  placeholder={i18n.t(
                    "app.teams.schedule-training.gathering-point.location.label",
                  )}
                />
              )}
              name="gatheringPoint.location"
            />
          </View>
          <View style={styles.formRow}>
            <View style={styles.iconContainer} />
            <View style={styles.dateRow}>
              <DatePickerButton
                unsetText={i18n.t(
                  "app.teams.schedule-training.gathering-point.gathering-until.label",
                )}
                control={control}
                type="date"
                name="gatheringPoint.gatheringUntil"
                label={i18n.t(
                  "app.teams.schedule-training.gathering-point.gathering-until.label",
                )}
              />
              <DatePickerButton
                allowRemoval
                unsetText=""
                style={{ maxWidth: 100 }}
                control={control}
                type="time"
                name="gatheringPoint.gatheringUntil"
                label={i18n.t(
                  "app.teams.schedule-training.gathering-point.gathering-until.label",
                )}
              />
            </View>
          </View>
          <Ruler />
          <View style={styles.formRow}>
            <View style={styles.iconContainer}>
              <MaterialCommunityIcons name="flag-outline" size={24} />
            </View>
            <View style={styles.dateRow}>
              <DatePickerButton
                unsetText={i18n.t(
                  "app.teams.schedule-training.acknowledgment.deadline.label",
                )}
                control={control}
                type="date"
                name="acknowledgmentSettings.deadline"
                label={i18n.t(
                  "app.teams.schedule-training.acknowledgment.deadline.label",
                )}
              />
              <DatePickerButton
                allowRemoval
                style={{ maxWidth: 100 }}
                unsetText=""
                control={control}
                type="time"
                name="acknowledgmentSettings.deadline"
                label={i18n.t(
                  "app.teams.schedule-training.acknowledgment.deadline.label",
                )}
              />
            </View>
          </View>
          <View style={styles.formRow}>
            <View style={styles.iconContainer} />
            <Text style={{ flex: 1, paddingBottom: 12 }}>
              {i18n.t(
                "app.teams.schedule-training.acknowledgment.deadline.hint",
              )}
            </Text>
          </View>
          <Ruler />
        </View>
      </ScrollView>
    </FormProvider>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  form: {
    marginTop: 16,
    flexDirection: "column",
  },
  formRow: {
    paddingHorizontal: 16,
    flexDirection: "row",
    alignItems: "center",
    width: "100%",
  },
  iconContainer: {
    width: 40,
    justifyContent: "flex-start",
  },
  dateRow: {
    flexDirection: "row",
    alignItems: "center",
    flex: 1,
  },
});

function dateTimeToPb(date: Date): MessageInitShape<typeof DateTimeSchema> {
  return {
    day: date.getDate(),
    month: date.getMonth() + 1,
    year: date.getFullYear(),
    hours: date.getHours(),
    minutes: date.getMinutes(),
    seconds: date.getSeconds(),
  };
}

function maybeGatheringPointToPb(
  gatheringPoint: ScheduleTrainingForm["gatheringPoint"],
): MessageInitShape<
  typeof ScheduleTrainingRequest_GatheringPointSchema
> | null {
  if (!gatheringPoint) {
    return null;
  }
  return {
    location: gatheringPoint.location,
    gatheringUntil: dateTimeToPb(gatheringPoint.gatheringUntil!),
  };
}

function maybeAcknowledgmentSettingsToPb(
  acknowledgmentSettings: ScheduleTrainingForm["acknowledgmentSettings"],
): MessageInitShape<
  typeof ScheduleTrainingRequest_AcknowledgementSettingsSchema
> | null {
  if (!acknowledgmentSettings) {
    return null;
  }
  return {
    deadline: dateTimeToPb(acknowledgmentSettings.deadline!),
  };
}
