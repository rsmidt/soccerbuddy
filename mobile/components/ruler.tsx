import { StyleSheet, View } from "react-native";
import { useTheme } from "react-native-paper";

export function Ruler() {
  const theme = useTheme();
  return (
    <View
      style={[styles.ruler, { backgroundColor: theme.colors.surfaceVariant }]}
    />
  );
}

const styles = StyleSheet.create({
  ruler: {
    height: 1,
    marginVertical: 12,
  },
});
