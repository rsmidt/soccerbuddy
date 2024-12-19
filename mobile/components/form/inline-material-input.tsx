import React, { ComponentProps } from "react";
import { TextInput } from "react-native";
import { useTheme } from "react-native-paper";

export function InlineMaterialInput(props: ComponentProps<typeof TextInput>) {
  const theme = useTheme();

  return (
    <TextInput
      {...props}
      style={{
        fontSize: theme.fonts.bodyLarge.fontSize,
        fontWeight: theme.fonts.bodyLarge.fontWeight,
        fontFamily: theme.fonts.bodyLarge.fontFamily,
        flex: 1,
      }}
      placeholderTextColor="gray"
    />
  );
}
