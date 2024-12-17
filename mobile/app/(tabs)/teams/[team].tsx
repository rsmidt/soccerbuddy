import { RefreshControl, ScrollView, StyleSheet, View } from "react-native";
import { Text, useTheme } from "react-native-paper";
import { useLocalSearchParams } from "expo-router";
import { useCallback, useState } from "react";
import { accountApi, useGetMeQuery } from "@/components/account/account-api";
import { useAppDispatch } from "@/store";
import { AccountLink } from "@/api/soccerbuddy/shared_pb";
import { Trans } from "react-i18next";
import { MaterialCommunityIcons } from "@expo/vector-icons";

export default function TeamTab() {
  const params = useLocalSearchParams();
  const dispatch = useAppDispatch();
  const [isRefreshing, setIsRefreshing] = useState(false);
  const { data, isLoading } = useGetMeQuery({});
  const theme = useTheme();

  const handleOnRefresh = useCallback(async () => {
    setIsRefreshing(true);
    // This refreshes the user data to refresh the list of teams.
    // Ideally, this would be more fine-grained.
    dispatch(accountApi.util.invalidateTags([{ type: "account", id: "me" }]));
    setIsRefreshing(false);
  }, [dispatch]);

  if (isLoading) {
    return (
      <View>
        <Text>Loading</Text>
      </View>
    );
  }

  // TODO: Extract to selector and extract hint as component.
  const parents = data!.linkedPersons.filter(
    (person) => person.linkedAs === AccountLink.LINKED_AS_PARENT,
  );
  const isParentHintVisible = parents.length !== 0;
  const name = isParentHintVisible
    ? `${parents[0].firstName}\xa0${parents[0].lastName}`
    : "";

  return (
    <ScrollView
      refreshControl={
        <RefreshControl refreshing={isRefreshing} onRefresh={handleOnRefresh} />
      }
    >
      {isParentHintVisible && (
        <View
          style={[
            styles.parentHint,
            { backgroundColor: theme.colors.surfaceVariant },
          ]}
        >
          <MaterialCommunityIcons name="shield-account-outline" size={36} />
          <Text
            style={[
              { color: theme.colors.onSurfaceVariant },
              styles.parentHintText,
            ]}
          >
            <Trans
              i18nKey="app.teams.parent_hint"
              default="You are seeing this hint because you are linked as parent to <bold>{{name}}</bold>"
              values={{ name }}
              components={{
                bold: <Text style={{ fontWeight: "bold" }}>{""}</Text>,
              }}
            />
          </Text>
        </View>
      )}
      <Text>{params.team}</Text>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  parentHint: {
    backgroundColor: "#fff",
    paddingVertical: 16,
    paddingHorizontal: 40,
    flexDirection: "row",
    gap: 16,
  },
  parentHintText: {
    flex: 1,
  },
});
