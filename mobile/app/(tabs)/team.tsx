import { Text, View } from "react-native";
import { useTranslation } from "react-i18next";

export default function Index() {
  const { t } = useTranslation();

  return (
    <View>
      <Text>{t("test2", { what: "World" })}</Text>
    </View>
  );
}
