import { NativeStackHeaderProps } from "@react-navigation/native-stack";
import { getHeaderTitle } from "@react-navigation/elements";
import { Appbar } from "react-native-paper";
import * as SecureStore from "expo-secure-store";
import { SESSION_TOKEN_KEY } from "@/components/auth/constants";
import React from "react";

function Header({
  navigation,
  route,
  options,
  ...rest
}: NativeStackHeaderProps) {
  const title = getHeaderTitle(options, route.name);

  const backEnabled =
    !!rest.back && options.headerBackButtonMenuEnabled !== false;

  return (
    <Appbar.Header>
      {backEnabled && <Appbar.BackAction onPress={navigation.goBack} />}
      <Appbar.Content title={title} />
      {options.headerRight?.({ canGoBack: backEnabled })}
      {/* TODO: Remove... */}
      <Appbar.Action
        icon="magnify"
        onPress={() => {
          SecureStore.deleteItemAsync(SESSION_TOKEN_KEY);
        }}
      />
    </Appbar.Header>
  );
}

export default Header;
