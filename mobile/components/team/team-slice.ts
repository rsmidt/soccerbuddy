import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import type { RootState } from "@/store";

/**
 * Represents the state for a single team.
 */
type TeamDetails = {
  parentHintRead: boolean;
};

type TeamState = {
  [teamId: string]: TeamDetails;
};

type MarkParentHintAsReadPayload = {
  teamId: string;
};

const initialState: TeamState = {};

const teamSlice = createSlice({
  name: "team",
  initialState,
  reducers: {
    /**
     * Marks the parent hint as read for a specific teamId.
     */
    markParentHintAsRead: (
      state,
      action: PayloadAction<MarkParentHintAsReadPayload>,
    ) => {
      const { teamId } = action.payload;
      if (!state[teamId]) {
        state[teamId] = { parentHintRead: true };
      } else {
        state[teamId].parentHintRead = true;
      }
    },
  },
});

/**
 * Selector to get the parentHintRead status for a given teamId.
 */
export const selectParentHintRead = (
  state: RootState,
  teamId: string,
): boolean => state.team[teamId]?.parentHintRead || false;

export const { actions, reducer } = teamSlice;

export const { markParentHintAsRead } = actions;
