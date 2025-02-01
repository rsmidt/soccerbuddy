import {
  ControllerProps,
  FieldPath,
  FieldValues,
} from "react-hook-form/dist/types";
import React, { useCallback, useState } from "react";
import { Controller, useFormContext } from "react-hook-form";
import { IconButton, Text, TouchableRipple } from "react-native-paper";
import { DateTime } from "luxon";
import { StyleProp, StyleSheet, View, ViewStyle } from "react-native";
import { MaterialIcons } from "@expo/vector-icons";
import { DatePickerModal, TimePickerModal } from "react-native-paper-dates";

export type DatePickerButtonProps<
  TFieldValues extends FieldValues = FieldValues,
  // @ts-ignore
  TName extends FieldPath<TFieldValues>,
> = Omit<ControllerProps<TFieldValues, TName>, "render"> & {
  label: string;
  type: "date" | "time";
  style?: StyleProp<ViewStyle>;
  unsetText?: string;
  allowRemoval?: boolean;
  minRenderYear?: number;
  maxRenderYear?: number;
  minDate?: Date;
  maxDate?: Date;
};

export function DatePickerButton<
  TFieldValues extends FieldValues = FieldValues,
  // @ts-ignore
  TName extends FieldPath<TFieldValues>,
>({
  control,
  name,
  rules,
  disabled,
  defaultValue,
  label,
  type,
  style,
  unsetText,
  minRenderYear,
  maxRenderYear,
  minDate,
  maxDate,
  allowRemoval = false,
}: DatePickerButtonProps<TFieldValues, TName>) {
  const [isPickerShown, setIsPickerShown] = useState(false);
  const { watch, resetField } = useFormContext<TFieldValues>();
  const onIconButtonPress = useCallback(() => {
    resetField(name);
  }, [resetField, name]);

  const date = watch(name);

  const isUnset = date === undefined;
  const formattedValue = isUnset
    ? unsetText
    : type === "date"
      ? DateTime.fromJSDate(date).toLocaleString({
          year: "numeric",
          month: "short",
          day: "numeric",
          weekday: "short",
        })
      : DateTime.fromJSDate(date).toLocaleString({
          hour: "2-digit",
          minute: "2-digit",
          hour12: false,
        });
  const isRemovalButtonShown = allowRemoval && date !== undefined;

  const defaultDateForTime = (() => {
    const jsDate = toDate(date);
    if (!jsDate) {
      return undefined;
    }
    return { hours: jsDate.getHours(), minutes: jsDate.getMinutes() };
  })();

  return (
    <View
      style={[
        styles.container,
        isRemovalButtonShown && { marginVertical: -8 },
        style,
      ]}
    >
      <TouchableRipple
        style={styles.button}
        onPress={() => setIsPickerShown(true)}
      >
        <Text style={[isUnset && { color: "gray" }]} variant="bodyLarge">
          {formattedValue}
        </Text>
      </TouchableRipple>
      <Controller
        control={control}
        rules={rules}
        disabled={disabled}
        defaultValue={defaultValue}
        render={({ field: { onChange, onBlur, value } }) =>
          type === "date" ? (
            <DatePickerModal
              mode="single"
              label={label}
              date={toDate(date)}
              onConfirm={({ date }) => {
                setIsPickerShown(false);
                if (date === undefined || Number.isNaN(date)) {
                  onChange(new Date());
                } else {
                  onChange(date);
                }
                onBlur();
              }}
              locale="en"
              startYear={minRenderYear}
              endYear={maxRenderYear}
              validRange={{
                startDate: minDate,
                endDate: maxDate,
              }}
              presentationStyle="pageSheet"
              onDismiss={() => {
                setIsPickerShown(false);
                onBlur();
              }}
              visible={isPickerShown}
            />
          ) : (
            <TimePickerModal
              use24HourClock
              visible={isPickerShown}
              minutes={defaultDateForTime?.minutes}
              hours={defaultDateForTime?.hours}
              onDismiss={() => {
                setIsPickerShown(false);
                onBlur();
              }}
              onConfirm={({ hours, minutes }) => {
                const jsDate = toDate(date);
                if (jsDate !== undefined) {
                  jsDate.setHours(hours);
                  jsDate.setMinutes(minutes);
                  onChange(jsDate);
                }
                setIsPickerShown(false);
              }}
            />
          )
        }
        name={name}
      />
      {isRemovalButtonShown && (
        <IconButton
          icon={(props) => <MaterialIcons {...props} name="highlight-remove" />}
          onPress={onIconButtonPress}
        />
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flexDirection: "row",
    alignItems: "center",
    flex: 1,
  },
  button: {
    flex: 1,
    paddingVertical: 8,
  },
});

function toDate(
  input: Date | undefined | null | number | string,
): Date | undefined {
  if (
    input === undefined ||
    input === null ||
    (typeof input === "number" && isNaN(input))
  ) {
    return undefined;
  }

  // If input is already a Date
  if (input instanceof Date) {
    return isNaN(input.getTime()) ? undefined : input;
  }

  // If input is a number, convert to Date
  if (typeof input === "number") {
    const date = new Date(input);
    return isNaN(date.getTime()) ? undefined : date;
  }

  // If input is a string, try to parse it as a Date
  if (typeof input === "string") {
    const date = new Date(input);
    return isNaN(date.getTime()) ? undefined : date;
  }

  // Fallback to undefined
  return undefined;
}
