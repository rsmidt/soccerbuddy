import { MaterialCommunityIcons, MaterialIcons } from "@expo/vector-icons";
import { Tabs } from "expo-router";
import Header from "@/components/header";
import { Avatar, BottomNavigation, useTheme } from "react-native-paper";
import { CommonActions } from "@react-navigation/native";
import i18n from "@/components/i18n";
import { useLayoutEffect } from "react";
import * as NavigationBar from "expo-navigation-bar";
import { useAppSelector } from "@/store/custom";
import { selectAuthenticatedState } from "@/components/auth/auth-slice";

export default function Layout() {
  const theme = useTheme();
  const state = useAppSelector((state) => selectAuthenticatedState(state.auth));
  if (state === undefined) {
    throw Error("auth state not supposed to be undefined");
  }
  const { user } = state;

  // That's nasty.
  useLayoutEffect(() => {
    let cb = () => {};
    const run = async () => {
      const currentColor = await NavigationBar.getBackgroundColorAsync();
      await NavigationBar.setBackgroundColorAsync(
        theme.colors.elevation.level2,
      );
      cb = () => {
        NavigationBar.setBackgroundColorAsync(currentColor);
      };
    };
    run();
    return cb;
  }, [theme]);

  return (
    <Tabs
      screenOptions={{
        header: (props) => <Header {...props} />,
        sceneStyle: {
          backgroundColor: theme.colors.background,
        },
      }}
      tabBar={({ navigation, state, insets, descriptors }) => (
        <BottomNavigation.Bar
          navigationState={state}
          safeAreaInsets={insets}
          onTabPress={({ route, preventDefault }) => {
            const event = navigation.emit({
              type: "tabPress",
              target: route.key,
              canPreventDefault: true,
            });

            if (event.defaultPrevented) {
              preventDefault();
            } else {
              navigation.dispatch({
                ...CommonActions.navigate(route.name, route.params),
                target: state.key,
              });
            }
          }}
          renderIcon={({ route, focused, color }) => {
            const { options } = descriptors[route.key];
            if (options.tabBarIcon) {
              return options.tabBarIcon({ focused, color, size: 24 });
            }

            return null;
          }}
          getLabelText={({ route }) => {
            const { options } = descriptors[route.key];
            return options.tabBarLabel !== undefined
              ? options.tabBarLabel
              : route.name;
          }}
        />
      )}
    >
      <Tabs.Screen
        options={{
          title: i18n.t("app.nav.home.title"),
          tabBarIcon: ({ focused }) => (
            <MaterialCommunityIcons
              name={focused ? "home" : "home-outline"}
              size={24}
            />
          ),
          tabBarLabel: i18n.t("app.nav.home"),
        }}
        name="index"
      />
      <Tabs.Screen
        options={{
          title: i18n.t("app.nav.team.title"),
          tabBarIcon: ({ focused }) => (
            <MaterialIcons
              name={focused ? "people" : "people-outline"}
              size={24}
            />
          ),
          tabBarLabel: i18n.t("app.nav.team"),
        }}
        name="teams"
      />
      <Tabs.Screen
        options={{
          title: i18n.t("app.nav.profile.title"),
          tabBarIcon: () => (
            <Avatar.Text
              size={24}
              label={user.firstName[0] + user.lastName[0]}
            />
          ),
          tabBarLabel: i18n.t("app.nav.profile"),
          tabBarShowLabel: true,
        }}
        name="profile/index"
      />
    </Tabs>
  );
}
