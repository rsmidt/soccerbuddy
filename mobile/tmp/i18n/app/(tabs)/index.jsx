import { Text, View } from "react-native";
import i18n from "@/components/i18n";
import { useGetMeQuery } from "@/components/account/account-api";
export default function Index() {
    const { data } = useGetMeQuery({});
    if (!data) {
        return null;
    }
    return (<View>
      <Text>{i18n.t("test2", { what: "World" })}</Text>
    </View>);
}
