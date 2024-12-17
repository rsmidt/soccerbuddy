import { Text } from "react-native-paper";
import { View } from "react-native";
import { Navigator } from "expo-router";
import { useGetMeQuery } from "@/components/account/account-api";
import i18n from "@/components/i18n";
import { CustomTabView } from "@/components/teams-tabs/tabs";

export default function Layout() {
  const { data, isLoading } = useGetMeQuery({});

  if (isLoading) {
    return (
      <View>
        <Text>Loading</Text>
      </View>
    );
  }

  const teamMemberships = data!.linkedPersons.flatMap(
    (person) => person.teamMemberships,
  );
  if (teamMemberships.length === 0) {
    return (
      <View>
        <Text>{i18n.t("app.teams.not_member")}</Text>
      </View>
    );
  }

  return (
    <Navigator>
      <CustomTabView teams={teamMemberships} />
    </Navigator>
  );
}
