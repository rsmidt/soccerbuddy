import { View } from "react-native";
import { Text } from "react-native-paper";
import { Stack, useLocalSearchParams } from "expo-router";
import i18n from "@/components/i18n";

export default function ScheduleTraining() {
  const { team } = useLocalSearchParams<{ team: string }>();

  return (
    <View>
      <Stack.Screen
        options={{
          headerTitle: i18n.t("app.teams.schedule-training.title"),
        }}
      />
      <Text>Hello {team}</Text>
    </View>
  );
}
