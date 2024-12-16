import { NativeStackHeaderProps } from "@react-navigation/native-stack";
import { getHeaderTitle } from "@react-navigation/elements";
import { Appbar } from "react-native-paper";
import { BottomTabHeaderProps } from "@react-navigation/bottom-tabs";

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
    </Appbar.Header>
  );
}

export default Header;
