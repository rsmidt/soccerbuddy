import { RefreshControl, ScrollView } from "react-native";
import { Text } from "react-native-paper";
import { useLocalSearchParams } from "expo-router";
import { useCallback, useState } from "react";
import { accountApi } from "@/components/account/account-api";
import { useAppDispatch } from "@/store";

export default function TeamTab() {
  const params = useLocalSearchParams();
  const dispatch = useAppDispatch();
  const [isRefreshing, setIsRefreshing] = useState(false);

  const handleOnRefresh = useCallback(async () => {
    setIsRefreshing(true);
    // This refreshes the user data to refresh the list of teams.
    // Ideally, this would be more fine-grained.
    dispatch(accountApi.util.invalidateTags([{ type: "account", id: "me" }]));
    setIsRefreshing(false);
  }, [dispatch]);

  return (
    <ScrollView
      refreshControl={
        <RefreshControl refreshing={isRefreshing} onRefresh={handleOnRefresh} />
      }
    >
      <Text>{params.team}</Text>
    </ScrollView>
  );
}
