import { PropsWithChildren, ReactNode } from "react";
import {
  GestureResponderEvent,
  StyleProp,
  StyleSheet,
  View,
  ViewStyle,
} from "react-native";
import { Text, TouchableRipple, useTheme } from "react-native-paper";

type ListItemProps = {
  left?: () => ReactNode;
  right?: () => ReactNode;
  onPress?: (event: GestureResponderEvent) => void;
  style?: StyleProp<ViewStyle>;
  label: string;
};

export function ListItem({
  left,
  right,
  style,
  label,
  onPress,
}: ListItemProps) {
  const theme = useTheme();

  const content = (
    <View
      style={[
        styles.container,
        { backgroundColor: theme.colors.surface },
        style,
      ]}
    >
      {left !== undefined && <View style={styles.left}>{left?.()}</View>}
      <View style={styles.content}>
        <Text
          style={{ color: theme.colors.onSurfaceVariant }}
          variant="bodyLarge"
        >
          {label}
        </Text>
      </View>
      {right !== undefined && <View style={styles.right}>{right?.()}</View>}
    </View>
  );

  if (onPress !== undefined) {
    return (
      <TouchableRipple borderless style={styles.ripple} onPress={onPress}>
        {content}
      </TouchableRipple>
    );
  }

  return content;
}

type ListItemAvatarProps = PropsWithChildren<{
  style?: StyleProp<ViewStyle>;
}>;

export function ListItemAvatar({ children, style }: ListItemAvatarProps) {
  return <View style={[styles.avatar, style]}>{children}</View>;
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    flexDirection: "row",
    paddingHorizontal: 16,
    gap: 16,
  },
  content: {
    alignSelf: "center",
  },
  right: {
    marginLeft: "auto",
    alignSelf: "center",
  },
  left: {},
  avatar: {
    height: 40,
  },
  ripple: {
    paddingVertical: 8,
  },
});
