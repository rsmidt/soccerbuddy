import { memo, useCallback, useEffect, useState } from "react";
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
import TeamActionsFab from "@/components/team/team-actions-fab";
import {
  markParentHintAsRead,
  selectParentHintRead,
} from "@/components/team/team-slice";
import ParentHint from "@/components/team/parent-hint";
import { TeamHome } from "@/components/team/team-home";
import {
  selectHasEditAllowance,
  selectPersonsWithParentLink,
  teamApi,
} from "@/components/team/team-api";

function TeamTab({ id }: { id: string }) {
  const dispatch = useAppDispatch();
  const isParentHintRead = useAppSelector((state) =>
    selectParentHintRead(state, id),
  );
  const { personWithLinkedParent, hasEditAllowance, isLoading } = useGetMeQuery(
    {},
    {
      selectFromResult: ({ data, ...rest }) => ({
        ...rest,
        personWithLinkedParent: selectPersonsWithParentLink(data, id),
        hasEditAllowance: selectHasEditAllowance(data, id),
      }),
    },
  );

  // Prefetch the team home data.
  const prefetch = teamApi.usePrefetch("getMyTeamHome");
  useEffect(() => {
    prefetch({ teamId: id });
  }, [id, prefetch]);

  const [isRefreshing, setIsRefreshing] = useState(false);
  const handleOnRefresh = useCallback(async () => {
    setIsRefreshing(true);
    // This refreshes the user data to refresh the list of teams.
    // Ideally, this would be more fine-grained.
    dispatch(accountApi.util.invalidateTags([{ type: "account", id: "me" }]));
    dispatch(teamApi.util.invalidateTags([{ type: "team", id }]));
    setIsRefreshing(false);
  }, [dispatch, id]);

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

  const isParentHintVisible =
    personWithLinkedParent !== undefined && !isParentHintRead;
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
          linkedPerson={personWithLinkedParent}
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
      {hasEditAllowance && <TeamActionsFab teamId={id} />}
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
