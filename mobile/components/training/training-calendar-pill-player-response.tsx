import { MaterialCommunityIcons } from "@expo/vector-icons";
import { Text, useTheme } from "react-native-paper";
import { StyleSheet, View } from "react-native";
import React, { ComponentProps } from "react";

type TrainingCalendarPillPlayerResponseProps = {
  accepted: number;
  tentative: number;
  declined: number;
  unknown: number;
};

export function TrainingCalendarPillPlayerResponse(
  props: TrainingCalendarPillPlayerResponseProps,
) {
  const { accepted, tentative, declined, unknown } = props;
  return (
    <View style={styles.container}>
      <IconWithText icon="check-circle-outline" text={accepted} />
      <IconWithText icon="minus-circle-outline" text={tentative} />
      <IconWithText icon="close-circle-outline" text={declined} />
      <IconWithText icon="help-circle-outline" text={unknown} />
    </View>
  );
}

function IconWithText(props: {
  icon: ComponentProps<typeof MaterialCommunityIcons>["name"];
  text: number;
}) {
  const { icon, text } = props;
  const theme = useTheme();

  return (
    <View style={styles.iconWithText}>
      <MaterialCommunityIcons name={icon} size={16} />
      <Text
        variant="labelMedium"
        style={{ color: theme.colors.onPrimaryContainer }}
      >
        {text}
      </Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flexDirection: "row",
    gap: 4,
  },
  iconWithText: {
    flexDirection: "row",
    alignItems: "center",
    gap: 2,
  },
});
