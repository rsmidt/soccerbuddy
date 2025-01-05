import React, { PropsWithChildren } from "react";
import { StyleProp, StyleSheet, View, ViewStyle } from "react-native";

type FormRowProps = PropsWithChildren<{
  style?: StyleProp<ViewStyle>;
}>;

export function FormRow({ children, style }: FormRowProps) {
  return <View style={[styles.formRow, style]}>{children}</View>;
}

type FormRowIconProps = PropsWithChildren<{
  style?: StyleProp<ViewStyle>;
}>;

export function FormRowIcon({ children, style }: FormRowIconProps) {
  return <View style={[styles.iconContainer, style]}>{children}</View>;
}

type FormRowControlsProps = PropsWithChildren<{
  direction?: "row" | "column";
  style?: StyleProp<ViewStyle>;
}>;

export function FormRowControls({
  children,
  style,
  direction = "column",
}: FormRowControlsProps) {
  return (
    <View style={[styles.controls, { flexDirection: direction }, style]}>
      {children}
    </View>
  );
}

FormRow.Icon = FormRowIcon;
FormRow.Controls = FormRowControls;

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
  controls: {
    flex: 1,
  },
});
