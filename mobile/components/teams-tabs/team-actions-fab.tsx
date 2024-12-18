import { FAB, Portal } from "react-native-paper";
import { useEffect, useState } from "react";
import i18n from "@/components/i18n";
import { StyleSheet } from "react-native";
import { useNavigation, useRouter } from "expo-router";

export type TeamActionsFabProps = {
  teamId: string;
};

function TeamActionsFab({ teamId }: TeamActionsFabProps) {
  const [fabState, setFabState] = useState({ open: false });
  const router = useRouter();
  const navigation = useNavigation();
  const [isFabVisible, setIsFabVisible] = useState(true);

  useEffect(() => {
    navigation.addListener("focus", () => setIsFabVisible(true));
    navigation.addListener("blur", () => setIsFabVisible(false));
  }, [navigation]);

  const onStateChange = ({ open }: { open: boolean }) => setFabState({ open });

  return (
    <Portal>
      <FAB.Group
        style={styles.fab}
        visible={isFabVisible}
        open={fabState.open}
        icon="plus"
        actions={[
          {
            icon: "calendar",
            label: i18n.t("app.teams.actions.add_event"),
            onPress: () => {
              router.navigate({
                pathname: "/teams/[team]/schedule-training",
                params: { team: teamId },
              });
            },
          },
        ]}
        onStateChange={onStateChange}
      />
    </Portal>
  );
}

const styles = StyleSheet.create({
  fab: {
    position: "absolute",
    margin: 16,
    right: 0,
    bottom: 80,
  },
});

export default TeamActionsFab;
