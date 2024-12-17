import { Text, View } from "react-native";
import { useGetMeQuery } from "@/components/account/account-api";
import i18n from "@/components/i18n";

export default function Index() {
  const { data } = useGetMeQuery({});

  return (
    <View>
      <Text>{i18n.t("test2", { what: "World" })}</Text>
    </View>
  );
}
