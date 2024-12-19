import { memo, useCallback, useState } from "react";
import { useAppDispatch, useAppSelector } from "@/store";
import { accountApi, useGetMeQuery } from "@/components/account/account-api";
import { Text } from "react-native-paper";
import {
  Image,
  RefreshControl,
  ScrollView,
  StyleSheet,
  View,
} from "react-native";
import { AccountLink } from "@/api/soccerbuddy/shared_pb";
import TeamActionsFab from "@/components/team/team-actions-fab";
import {
  markParentHintAsRead,
  selectParentHintRead,
} from "@/components/team/team-slice";
import ParentHint from "@/components/team/parent-hint";
import {
  GetMeResponse,
  GetMeResponse_LinkedPerson,
} from "@/api/soccerbuddy/account/v1/account_service_pb";
import { TeamHome } from "@/components/team/team-home";

/**
 * Selects any person wih a parent link ONLY when there's no person with a self link.
 * We do this because we assume that parents often do not really know about the names of their children teams.
 */
function selectPersonsWithParentLink(
  data?: GetMeResponse,
): GetMeResponse_LinkedPerson | null {
  if (!data) return null;

  const personsLinkedWithParent = data.linkedPersons.filter(
    (person) => person.linkedAs === AccountLink.LINKED_AS_PARENT,
  );
  const hasPersonsWithSelfLink = data.linkedPersons.some(
    (person) => person.linkedAs === AccountLink.LINKED_AS_SELF,
  );
  if (hasPersonsWithSelfLink || personsLinkedWithParent.length === 0) {
    return null;
  }
  return personsLinkedWithParent[0];
}

function TeamTab({ id }: { id: string }) {
  const dispatch = useAppDispatch();
  const isParentHintRead = useAppSelector((state) =>
    selectParentHintRead(state, id),
  );
  const [isRefreshing, setIsRefreshing] = useState(false);
  const { linkedPerson, isLoading } = useGetMeQuery(
    {},
    {
      selectFromResult: ({ data, ...rest }) => ({
        ...rest,
        linkedPerson: selectPersonsWithParentLink(data),
      }),
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
  const teamImageUrl = "https://p.rsmidt.dev/500x500?bg=e8e7e9";

  return (
    <ScrollView
      style={styles.container}
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
      <View style={styles.teamBannerFrame}>
        <Image
          style={styles.teamBanner}
          source={{
            uri: teamImageUrl,
          }}
        />
      </View>
      <TeamHome teamId={id} style={styles.home} />
      <TeamActionsFab teamId={id} />
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  home: {
    marginTop: 16,
  },
  teamBannerFrame: {
    flexDirection: "row",
    marginTop: 16,
    marginHorizontal: 16,
    alignSelf: "center",
    height: 300,
    padding: 16,
    backgroundColor: "white",
  },
  teamBanner: {
    flex: 1,
    width: "100%",
    resizeMode: "cover",
  },
});

export default memo(TeamTab);
