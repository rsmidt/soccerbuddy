import { MaterialBottomTabs } from "@/components/material-bottom-tabs";
import { MaterialCommunityIcons } from "@expo/vector-icons";

export default function Layout() {
  return (
    <MaterialBottomTabs>
      <MaterialBottomTabs.Screen
        options={{
          tabBarIcon: ({ color }) => <MaterialCommunityIcons name="home" />,
          tabBarLabel: "Home",
        }}
        name="index"
      />
    </MaterialBottomTabs>
  );
}
