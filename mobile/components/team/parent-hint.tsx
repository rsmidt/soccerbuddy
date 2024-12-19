import { GetMeResponse_LinkedPerson } from "@/api/soccerbuddy/account/v1/account_service_pb";
import { StyleSheet, View } from "react-native";
import { MaterialCommunityIcons } from "@expo/vector-icons";
import { Button, Text, useTheme } from "react-native-paper";
import { Trans } from "react-i18next";
import i18n from "i18next";

type ParentHintProps = {
  linkedPerson: GetMeResponse_LinkedPerson;
  onParentHintDismissed: () => void;
};

function ParentHint({ linkedPerson, onParentHintDismissed }: ParentHintProps) {
  const theme = useTheme();

  const name = `${linkedPerson.firstName}\xa0${linkedPerson.lastName}`;

  return (
    <View
      style={[
        styles.parentHint,
        { backgroundColor: theme.colors.surfaceVariant },
      ]}
    >
      <MaterialCommunityIcons name="shield-account-outline" size={36} />
      <View style={styles.textContainer}>
        <Text
          style={[
            { color: theme.colors.onSurfaceVariant },
            styles.parentHintText,
          ]}
        >
          <Trans
            i18nKey="app.teams.parent_hint"
            default="You are seeing this hint because you are linked as parent to <bold>{{name}}</bold>"
            values={{ name }}
            components={{
              bold: <Text style={{ fontWeight: "bold" }}>{""}</Text>,
            }}
          />
        </Text>
        <Button style={styles.dismissButton} onPress={onParentHintDismissed}>
          {i18n.t("app.teams.parent_hint.dismiss")}
        </Button>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  parentHint: {
    backgroundColor: "#fff",
    paddingTop: 16,
    paddingBottom: 8,
    paddingHorizontal: 40,
    flexDirection: "row",
    gap: 16,
  },
  parentHintText: {
    flex: 1,
  },
  textContainer: {
    flex: 1,
  },
  dismissButton: {
    alignSelf: "flex-end",
  },
});

export default ParentHint;
