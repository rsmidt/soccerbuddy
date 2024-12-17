import { Text, View } from "react-native";
import i18n from "@/components/i18n";
import { useGetMeQuery } from "@/components/account/account-api";
import { timestampDate } from "@bufbuild/protobuf/wkt";

export default function Index() {
  const { data } = useGetMeQuery({});

  if (!data) {
    return null;
  }

  console.log(timestampDate(data.linkedPersons[0].linkedAt!));

  return (
    <View>
      <Text>{i18n.t("test2", { what: "World" })}</Text>
    </View>
  );
}
