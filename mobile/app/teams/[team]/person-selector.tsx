import {
  selectTeamMemberInitials,
  selectTeamMemberName,
  selectTeamMembersByMode,
  useListTeamMembersQuery,
} from "@/components/team/team-api";
import { Stack, useLocalSearchParams } from "expo-router";
import { FlatList, View } from "react-native";
import { useAppDispatch, useAppSelector } from "@/store/custom";
import { Avatar, Checkbox, Text } from "react-native-paper";
import { ListTeamMembersResponse_Member } from "@/api/soccerbuddy/team/v1/team_service_pb";
import { ListItem, ListItemAvatar } from "@/components/list/list";
import {
  NominationClass,
  selectNominatedPersons,
  togglePersonNomination,
} from "@/components/training/schedule-training-slice";
import i18n from "@/components/i18n";

export default function PersonSelector() {
  const { team, mode } = useLocalSearchParams<{
    team: string;
    mode: NominationClass;
  }>();
  const { data, isLoading } = useListTeamMembersQuery(
    { teamId: team },
    {
      selectFromResult: ({ data, ...rest }) => ({
        ...rest,
        data: selectTeamMembersByMode(data, mode),
      }),
    },
  );
  const selectedPersons = useAppSelector((state) =>
    selectNominatedPersons(state, mode),
  );
  const dispatch = useAppDispatch();

  if (isLoading || !data) {
    return <Text>Loading...</Text>;
  }

  function isSelected(member: ListTeamMembersResponse_Member): boolean {
    return selectedPersons.includes(member.personId);
  }

  function handleListItemClick(item: ListTeamMembersResponse_Member) {
    dispatch(togglePersonNomination({ mode, personId: item.personId }));
  }

  return (
    <View>
      <Stack.Screen
        options={{
          title: i18n.t("app.teams.schedule-training.person-selector.title", {
            mode,
          }),
        }}
      />
      <FlatList
        data={data}
        renderItem={({ item }) => (
          <ListItem
            key={item.id}
            label={selectTeamMemberName(item)}
            onPress={() => handleListItemClick(item)}
            left={() => (
              <ListItemAvatar>
                <Avatar.Text size={40} label={selectTeamMemberInitials(item)} />
              </ListItemAvatar>
            )}
            right={() => (
              <Checkbox
                status={isSelected(item) ? "checked" : "unchecked"}
                onPress={() => handleListItemClick(item)}
              />
            )}
          />
        )}
      />
    </View>
  );
}
