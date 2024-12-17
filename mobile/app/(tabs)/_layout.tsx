import { MaterialCommunityIcons, MaterialIcons } from "@expo/vector-icons";
import { useTranslation } from "react-i18next";
import { Tabs } from "expo-router";
import Header from "@/components/header";

export default function Layout() {
  const { t } = useTranslation();

  return (
    <Tabs
      screenOptions={{
        header: (props) => <Header {...props} />,
      }}
    >
      <Tabs.Screen
        options={{
          title: t("app.nav.home.title"),
          tabBarIcon: ({ focused }) => (
            <MaterialCommunityIcons
              name={focused ? "home" : "home-outline"}
              size={24}
            />
          ),
          tabBarLabel: t("app.nav.home"),
        }}
        name="index"
      />
      <Tabs.Screen
        options={{
          title: t("app.nav.team.title"),
          tabBarIcon: ({ focused }) => (
            <MaterialIcons
              name={focused ? "people" : "people-outline"}
              size={24}
            />
          ),
          tabBarLabel: t("app.nav.team"),
        }}
        name="teams"
      />
    </Tabs>
  );
}
