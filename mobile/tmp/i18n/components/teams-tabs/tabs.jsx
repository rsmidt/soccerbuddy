import { TabBar, TabView, } from "react-native-tab-view";
import { Text, useTheme } from "react-native-paper";
import { Navigator, Slot, useRouter } from "expo-router";
import { CommonActions } from "@react-navigation/native";
import { useEffect } from "react";
export function CustomTabBarLabel({ style, color, labelText, focused, }) {
    return (<Text style={[
            { color },
            style,
            {
                textAlign: "center",
                fontSize: 14,
                color: "black",
                backgroundColor: "transparent",
                // Dirty hack because otherwise the label text would disappear once focussed (and getting bold).
                width: "101%",
            },
            focused && {
                fontWeight: "bold",
            },
        ]}>
      {labelText}
    </Text>);
}
export function CustomTabBar({ state, ...props }) {
    const theme = useTheme();
    const { navigation } = Navigator.useContext();
    const options = Object.fromEntries(state.routes.map((route) => {
        const tab = {
            labelText: route.name,
            label: (props) => <CustomTabBarLabel {...props}/>,
            labelAllowFontScaling: true,
        };
        return [route.key, tab];
    }));
    return (<TabBar {...props} scrollEnabled onTabPress={({ route, preventDefault }) => {
            const event = navigation.emit({
                type: "tabPress",
                target: route.key,
                canPreventDefault: true,
            });
            if ("defaultPrevented" in event && event.defaultPrevented) {
                preventDefault();
            }
        }} onTabLongPress={({ route }) => {
            navigation.emit({
                type: "tabLongPress",
                target: route.key,
            });
        }} indicatorStyle={{
            backgroundColor: theme.colors.primary,
        }} options={options} style={{
            backgroundColor: theme.colors.background,
        }} navigationState={state}/>);
}
// This is heavily inspired by https://github.com/EvanBacon/evanbacon.dev/tree/master
// and the implementation of Material Tab Bar.
export function CustomTabView({ teams }) {
    const { state: state1, navigation } = Navigator.useContext();
    const router = useRouter();
    const catchAllRoute = state1.routes.find((route, i) => state1.index === i);
    const params = catchAllRoute?.params;
    // TODO: make this setting persistent.
    // This is also a very dirty workaround. Hear me out on this rant.
    // So I've been using Expo now because so many folks on the internet recommend it.
    // However, the simple use case CANNOT BE ACHIEVED from the docs alone:
    // rendering a tab list dynamically based on a list of something fetched from somewhere.
    // This is totally a use case even mentioned in the freaking Material docs
    // https://m3.material.io/components/tabs/guidelines#2691a4ac-ea12-467e-8568-de0c024e89e6
    // Why is this so hard with expo???????????????????????????
    const firstTeamId = teams[0].id;
    useEffect(() => {
        if (!params?.team) {
            router.navigate({
                pathname: "/teams/[team]",
                params: { team: firstTeamId },
            });
        }
    }, [firstTeamId, params?.team, router]);
    if (!params?.team) {
        return null;
    }
    const index = teams.findIndex((team) => team.id === params.team);
    const state = {
        index,
        routes: teams.map((team) => ({
            key: team.id,
            name: team.name,
        })),
    };
    return (<TabView onIndexChange={(index) => {
            const route = state.routes[index];
            navigation.dispatch({
                ...CommonActions.navigate({
                    name: "[team]",
                    merge: true,
                    params: { team: route.key },
                }),
                target: state1.key,
            });
        }} renderTabBar={(props) => <CustomTabBar state={state} {...props}/>} navigationState={state} renderScene={() => <Slot />}/>);
}
