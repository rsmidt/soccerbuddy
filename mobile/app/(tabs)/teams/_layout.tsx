import { useGetMeQuery } from "@/components/account/account-api";
import { View } from "react-native";
import { Text } from "react-native-paper";
import i18n from "@/components/i18n";
import { Navigator } from "expo-router";
import { CustomTabView } from "@/components/team/tabs";

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
