import { RefreshControl, ScrollView, View } from "react-native";
import { Text } from "react-native-paper";
import { useLocalSearchParams } from "expo-router";
import { useCallback, useState } from "react";
import { accountApi, useGetMeQuery } from "@/components/account/account-api";
import { useAppDispatch } from "@/store";
import { AccountLink } from "@/api/soccerbuddy/shared_pb";
import { Trans } from "react-i18next";
export default function TeamTab() {
    const params = useLocalSearchParams();
    const dispatch = useAppDispatch();
    const [isRefreshing, setIsRefreshing] = useState(false);
    const { data, isLoading } = useGetMeQuery({});
    const handleOnRefresh = useCallback(async () => {
        setIsRefreshing(true);
        // This refreshes the user data to refresh the list of teams.
        // Ideally, this would be more fine-grained.
        dispatch(accountApi.util.invalidateTags([{ type: "account", id: "me" }]));
        setIsRefreshing(false);
    }, [dispatch]);
    if (isLoading) {
        return (<View>
        <Text>Loading</Text>
      </View>);
    }
    const parents = data.linkedPersons.filter((person) => person.linkedAs === AccountLink.LINKED_AS_PARENT);
    const isParentHintVisible = parents.length !== 0;
    const name = isParentHintVisible
        ? `${parents[0].firstName} ${parents[0].lastName}`
        : "";
    return (<ScrollView refreshControl={<RefreshControl refreshing={isRefreshing} onRefresh={handleOnRefresh}/>}>
      {isParentHintVisible && (<View>
          <Text>
            <Trans i18nKey="app.teams.parent_hint" default="You are seeing this hint because you are linked as parent to <bold>{{name}}</bold>" values={{
                name,
            }} components={{
                bold: <Text style={{ fontWeight: "bold" }}/>,
            }}/>
          </Text>
        </View>)}
      <Text>{params.team}</Text>
    </ScrollView>);
}
