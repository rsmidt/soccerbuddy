import { NativeStackHeaderProps } from "@react-navigation/native-stack";
import { getHeaderTitle } from "@react-navigation/elements";
import { Appbar, useTheme } from "react-native-paper";
import { BottomTabHeaderProps } from "@react-navigation/bottom-tabs";
import * as SecureStore from "expo-secure-store";
import { SESSION_TOKEN_KEY } from "@/components/auth/constants";
import { StatusBar } from "expo-status-bar";
import React from "react";

function Header({
  navigation,
  route,
  options,
  ...rest
}: NativeStackHeaderProps | BottomTabHeaderProps) {
  const title = getHeaderTitle(options, route.name);
  const theme = useTheme();

  return (
    <>
      <StatusBar backgroundColor={theme.colors.surface} />
      <Appbar.Header>
        {"back" in rest && rest.back ? (
          <Appbar.BackAction onPress={navigation.goBack} />
        ) : null}
        <Appbar.Content title={title} />
        {/* TODO: Remove... */}
        <Appbar.Action
          icon="magnify"
          onPress={() => {
            SecureStore.deleteItemAsync(SESSION_TOKEN_KEY);
          }}
        />
      </Appbar.Header>
    </>
  );
}

export default Header;
