import {
  createMaterialBottomTabNavigator,
  MaterialBottomTabNavigationOptions,
} from "react-native-paper/react-navigation";

import { withLayoutContext } from "expo-router";
import { ParamListBase, StackNavigationState } from "@react-navigation/native";
import { NativeStackNavigationEventMap } from "@react-navigation/native-stack";

const { Navigator } = createMaterialBottomTabNavigator();

export const MaterialBottomTabs = withLayoutContext<
  MaterialBottomTabNavigationOptions,
  typeof Navigator,
  StackNavigationState<ParamListBase>,
  NativeStackNavigationEventMap
>(Navigator);
