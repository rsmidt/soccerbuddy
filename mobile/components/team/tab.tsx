import { memo, useCallback, useState } from "react";
import { useAppDispatch, useAppSelector } from "@/store";
import { accountApi, useGetMeQuery } from "@/components/account/account-api";
import { Text } from "react-native-paper";
import { RefreshControl, ScrollView, StyleSheet, View } from "react-native";
import { AccountLink } from "@/api/soccerbuddy/shared_pb";
import TeamActionsFab from "@/components/team/team-actions-fab";
import {
  markParentHintAsRead,
  selectParentHintRead,
} from "@/components/team/team-slice";
import ParentHint from "@/components/team/parent-hint";

function TeamTab({ id }: { id: string }) {
  const dispatch = useAppDispatch();
  const isParentHintRead = useAppSelector((state) =>
    selectParentHintRead(state, id),
  );
  const [isRefreshing, setIsRefreshing] = useState(false);
  const { linkedPerson, isLoading } = useGetMeQuery(
    {},
    {
      selectFromResult: ({ data, isLoading }) => {
        const personsLinkedWithParent = data!.linkedPersons.filter(
          (person) => person.linkedAs === AccountLink.LINKED_AS_PARENT,
        );
        if (personsLinkedWithParent.length === 0) {
          return { linkedPerson: null, isLoading };
        }
        const firstPerson = personsLinkedWithParent[0];
        return { linkedPerson: firstPerson, isLoading };
      },
    },
  );

  const handleOnRefresh = useCallback(async () => {
    setIsRefreshing(true);
    // This refreshes the user data to refresh the list of teams.
    // Ideally, this would be more fine-grained.
    dispatch(accountApi.util.invalidateTags([{ type: "account", id: "me" }]));
    setIsRefreshing(false);
  }, [dispatch]);
  const onParentHintDismissed = useCallback(() => {
    dispatch(markParentHintAsRead({ teamId: id }));
  }, [dispatch, id]);

  if (isLoading) {
    return (
      <View>
        <Text>Loading</Text>
      </View>
    );
  }

  const isParentHintVisible = linkedPerson !== null && !isParentHintRead;

  return (
    <ScrollView
      refreshControl={
        <RefreshControl refreshing={isRefreshing} onRefresh={handleOnRefresh} />
      }
    >
      {isParentHintVisible && (
        <ParentHint
          linkedPerson={linkedPerson}
          onParentHintDismissed={onParentHintDismissed}
        />
      )}
      <Text>{id}</Text>
      <TeamActionsFab teamId={id} />
    </ScrollView>
  );
}

const styles = StyleSheet.create({});

export default memo(TeamTab);
