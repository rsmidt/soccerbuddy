import { ScrollView, View } from "react-native";
import { Avatar } from "react-native-paper";
import { Image } from "expo-image";
import { useAuthenticatedState } from "@/components/auth/auth-slice";
import i18n from "@/components/i18n";
import { ButtonList, ButtonListItem } from "@/components/list/button-list";

export default function ProfileIndex() {
  const { user } = useAuthenticatedState();
  return (
    <ScrollView>
      <View
        style={{
          paddingHorizontal: 16,
          gap: 24,
        }}
      >
        <ButtonList>
          <ButtonListItem
            title={`${user.firstName} ${user.lastName}`}
            supportingText={i18n.t("app.profile.account.supporting-text")}
            icon={() => (
              <Avatar.Text
                size={24}
                label={user.firstName[0] + user.lastName[0]}
              />
            )}
          />
        </ButtonList>
        <ButtonList>
          <ButtonListItem
            title={i18n.t("app.profile.clubs.headline")}
            supportingText={i18n.t("app.profile.clubs.supporting-text")}
            icon={() => (
              <Image
                source={require("@/assets/images/club-icon.png")}
                style={{ flex: 1, width: "100%" }}
                contentFit="cover"
              />
            )}
          />
        </ButtonList>
      </View>
    </ScrollView>
  );
}
