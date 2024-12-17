import { MaterialCommunityIcons, MaterialIcons } from "@expo/vector-icons";
import { Tabs } from "expo-router";
import Header from "@/components/header";
import { BottomNavigation, useTheme } from "react-native-paper";
import { CommonActions } from "@react-navigation/native";
import i18n from "@/components/i18n";
export default function Layout() {
    const theme = useTheme();
    return (<Tabs screenOptions={{
            header: (props) => <Header {...props}/>,
            sceneStyle: {
                backgroundColor: theme.colors.background,
            },
        }} tabBar={({ navigation, state, insets, descriptors }) => (<BottomNavigation.Bar navigationState={state} safeAreaInsets={insets} onTabPress={({ route, preventDefault }) => {
                const event = navigation.emit({
                    type: "tabPress",
                    target: route.key,
                    canPreventDefault: true,
                });
                if (event.defaultPrevented) {
                    preventDefault();
                }
                else {
                    navigation.dispatch({
                        ...CommonActions.navigate(route.name, route.params),
                        target: state.key,
                    });
                }
            }} renderIcon={({ route, focused, color }) => {
                const { options } = descriptors[route.key];
                if (options.tabBarIcon) {
                    return options.tabBarIcon({ focused, color, size: 24 });
                }
                return null;
            }} getLabelText={({ route }) => {
                const { options } = descriptors[route.key];
                return options.title !== undefined ? options.title : route.name;
            }}/>)}>
      <Tabs.Screen options={{
            title: i18n.t("app.nav.home.title"),
            tabBarIcon: ({ focused }) => (<MaterialCommunityIcons name={focused ? "home" : "home-outline"} size={24}/>),
            tabBarLabel: i18n.t("app.nav.home"),
        }} name="index"/>
      <Tabs.Screen options={{
            title: i18n.t("app.nav.team.title"),
            tabBarIcon: ({ focused }) => (<MaterialIcons name={focused ? "people" : "people-outline"} size={24}/>),
            tabBarLabel: i18n.t("app.nav.team"),
        }} name="teams"/>
    </Tabs>);
}
