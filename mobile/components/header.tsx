import { NativeStackHeaderProps } from "@react-navigation/native-stack";
import { getHeaderTitle } from "@react-navigation/elements";
import { Appbar } from "react-native-paper";
import { BottomTabHeaderProps } from "@react-navigation/bottom-tabs";
import * as SecureStore from "expo-secure-store";
import { SESSION_TOKEN_KEY } from "@/components/auth/constants";

function Header({
  navigation,
  route,
  options,
  ...rest
}: NativeStackHeaderProps | BottomTabHeaderProps) {
  const title = getHeaderTitle(options, route.name);

  return (
    <Appbar.Header>
      {"back" in rest && rest.back ? (
        <Appbar.BackAction onPress={navigation.goBack} />
      ) : null}
      <Appbar.Content title={title} />
      {/* TODO: Remove... */}
      <Appbar.Action
        icon="mangify"
        onPress={() => {
          SecureStore.deleteItemAsync(SESSION_TOKEN_KEY);
        }}
      />
    </Appbar.Header>
  );
}

export default Header;
