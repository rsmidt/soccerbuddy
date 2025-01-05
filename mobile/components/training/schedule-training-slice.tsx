import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { createAppSelector } from "@/store/custom";

export type NominationClass = "player" | "staff";

type ScheduleTrainingState = {
  nominations: Record<
    NominationClass,
    {
      ids: string[];
      isDirty: boolean;
    }
  >;
};

const initialState = {
  nominations: {
    player: {
      ids: [],
      isDirty: false,
    },
    staff: {
      ids: [],
      isDirty: false,
    },
  },
  isDirty: false,
} as ScheduleTrainingState;

const scheduleTrainingSlice = createSlice({
  name: "scheduleTraining",
  initialState,
  reducers: {
    resetNominations: (state, action: PayloadAction<NominationClass>) => {
      state.nominations[action.payload].ids = [];
      state.nominations[action.payload].isDirty = false;
    },
    nominatePlayer: (
      state,
      action: PayloadAction<string | readonly string[]>,
    ) => {
      if (Array.isArray(action.payload)) {
        state.nominations.player.ids.push(...action.payload);
      } else {
        state.nominations.player.ids.push(action.payload as string);
      }
      state.nominations.player.isDirty = true;
    },
    nominateStaff: (
      state,
      action: PayloadAction<string | readonly string[]>,
    ) => {
      if (Array.isArray(action.payload)) {
        state.nominations.staff.ids.push(...action.payload);
      } else {
        state.nominations.staff.ids.push(action.payload as string);
      }
      state.nominations.staff.isDirty = true;
    },
    togglePersonNomination: (
      state,
      action: PayloadAction<{
        mode: NominationClass;
        personId: string;
      }>,
    ) => {
      const { mode, personId } = action.payload;
      const currentNominations = state.nominations[mode];
      const index = currentNominations.ids.indexOf(personId);

      if (index === -1) {
        currentNominations.ids.push(personId);
      } else {
        currentNominations.ids.splice(index, 1);
      }
      currentNominations.isDirty = true;
    },
  },
});

/**
 * Selects if the state is dirty for the given mode.
 */
export const selectIsSelectionDirtyByMode = createAppSelector(
  [
    (state) => state.scheduleTraining.nominations,
    (state, mode: NominationClass) => mode,
  ],
  (nominations, mode) => nominations[mode].isDirty,
);

const EMPTY_ARRAY: readonly any[] = [];

export const selectNominatedPersons = createAppSelector(
  [
    (state) => state.scheduleTraining.nominations,
    (state, mode: NominationClass) => mode,
  ],
  (nominationsByClass, mode) => nominationsByClass[mode].ids ?? EMPTY_ARRAY,
);

export const { actions, reducer } = scheduleTrainingSlice;

export const {
  nominatePlayer,
  nominateStaff,
  togglePersonNomination,
  resetNominations,
} = actions;
