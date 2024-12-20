import { FAB, Portal } from "react-native-paper";
import { useEffect, useState } from "react";
import i18n from "@/components/i18n";
import { StyleSheet } from "react-native";
import { Navigator, useNavigation, useRouter } from "expo-router";
import { NavigationState } from "react-native-tab-view";

export type TeamActionsFabProps = {
  teamId: string;
};

function TeamActionsFab({ teamId }: TeamActionsFabProps) {
  const [fabState, setFabState] = useState({ open: false });
  const router = useRouter();
  const navigation = useNavigation();
  const [isFabVisible, setIsFabVisible] = useState(true);
  const isActiveTab = useIsActiveTab(teamId);

  useEffect(() => {
    const focusRef = navigation.addListener("focus", () =>
      setIsFabVisible(true),
    );
    const blurRef = navigation.addListener("blur", () =>
      setIsFabVisible(false),
    );

    return () => {
      navigation.removeListener("focus", focusRef);
      navigation.removeListener("blur", blurRef);
    };
  }, [navigation]);

  const onStateChange = ({ open }: { open: boolean }) => setFabState({ open });

  return (
    <Portal>
      <FAB.Group
        style={styles.fab}
        visible={isFabVisible && isActiveTab}
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
    right: 0,
    bottom: 80,
  },
});

export default TeamActionsFab;

function selectActiveTeamId(state: NavigationState<any>): string | null {
  return (
    state.routes.find((route, i) => state.index === i)?.params?.team ?? null
  );
}

function useIsActiveTab(teamId: string): boolean {
  const { state } = Navigator.useContext();

  return selectActiveTeamId(state) === teamId;
}
