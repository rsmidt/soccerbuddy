import { FormRow } from "@/components/form/form-row";
import { MaterialCommunityIcons } from "@expo/vector-icons";
import {
  getFilterCriteria,
  useListTeamMembersQuery,
} from "@/components/team/team-api";
import { Text, TouchableRipple } from "react-native-paper";
import { useFocusEffect, useRouter } from "expo-router";
import { useCallback } from "react";
import { useAppDispatch, useAppSelector } from "@/store/custom";
import {
  nominatePlayer,
  nominateStaff,
  NominationClass,
  resetNominations,
  selectIsSelectionDirtyByMode,
  selectNominatedPersons,
} from "@/components/training/schedule-training-slice";
import { PayloadAction } from "@reduxjs/toolkit";
import i18n from "@/components/i18n";
import { StyleSheet } from "react-native";

const EMPTY_ARRAY: any[] = [];

type NominationsProps = {
  teamId: string;
  nominatedPersonIds: readonly string[];
  onNominationsChanged?: (nominations: readonly string[]) => void;
  mode: NominationClass;
};

export function Nominations({
  teamId,
  mode,
  nominatedPersonIds,
  onNominationsChanged,
}: NominationsProps) {
  const router = useRouter();
  const { data, isLoading } = useListTeamMembersQuery({ teamId });
  const dispatch = useAppDispatch();

  const defaultPersonSelection =
    data?.members
      ?.filter(getFilterCriteria(mode))
      ?.map((person) => person.personId) ?? EMPTY_ARRAY;

  const onOpenPlayerSelector = useCallback(() => {
    // Only set a default selection if the currently nominated persons are empty.
    if (nominatedPersonIds.length === 0) {
      dispatch(nominatePersons(mode, defaultPersonSelection));
    } else {
      dispatch(nominatePersons(mode, nominatedPersonIds));
    }

    router.navigate({
      pathname: "/teams/[team]/person-selector",
      params: {
        mode,
        team: teamId,
      },
    });
  }, [
    nominatedPersonIds,
    router,
    teamId,
    dispatch,
    mode,
    defaultPersonSelection,
  ]);

  const isSelectionDirty = useAppSelector((state) =>
    selectIsSelectionDirtyByMode(state, mode),
  );
  const newlyNominatedPersonIds = useAppSelector((state) =>
    selectNominatedPersons(state, mode),
  );
  // When returning from the selector, trigger the callback.
  useFocusEffect(
    useCallback(() => {
      if (isSelectionDirty) {
        onNominationsChanged?.(newlyNominatedPersonIds);
        dispatch(resetNominations(mode));
      }
    }, [
      isSelectionDirty,
      onNominationsChanged,
      newlyNominatedPersonIds,
      dispatch,
      mode,
    ]),
  );

  if (isLoading || !data) {
    return <Text>Loading...</Text>;
  }

  const buttonText = getButtonText(mode, nominatedPersonIds, data);

  return (
    <FormRow>
      <FormRow.Icon>
        <MaterialCommunityIcons
          name={mode === "player" ? "run" : "checkerboard"}
          size={24}
        />
      </FormRow.Icon>
      <FormRow.Controls>
        <TouchableRipple style={styles.ripple} onPress={onOpenPlayerSelector}>
          <Text variant="bodyLarge">{buttonText}</Text>
        </TouchableRipple>
      </FormRow.Controls>
    </FormRow>
  );
}

const styles = StyleSheet.create({
  ripple: {
    paddingVertical: 8,
  },
});

function nominatePersons(
  mode: NominationClass,
  personIds: readonly string[],
): PayloadAction<string | readonly string[]> {
  if (mode === "player") {
    return nominatePlayer(personIds);
  } else {
    return nominateStaff(personIds);
  }
}

function getButtonText(
  mode: NominationClass,
  nominatedPersonIds: readonly string[],
  data: ReturnType<typeof useListTeamMembersQuery>["data"],
) {
  const filteredMembers =
    data?.members?.filter(getFilterCriteria(mode)).length ?? 0;
  const nominatedCount = nominatedPersonIds.length;

  if (nominatedCount === 0) {
    return i18n.t("app.teams.schedule-training.no-nominations", { mode });
  }

  if (nominatedCount === filteredMembers) {
    return i18n.t("app.teams.schedule-training.all-nominated", {
      mode,
      nominatedCount,
    });
  }

  return i18n.t("app.teams.schedule-training.some-nominated", {
    mode,
    nominatedCount,
    totalCount: filteredMembers,
  });
}
