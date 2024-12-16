import { createAsyncThunk } from "@reduxjs/toolkit";
import { AppDispatch, RootState } from "./index";
import { createSelector } from "reselect";

export const createAppAsyncThunk = createAsyncThunk.withTypes<{
  state: RootState;
  dispatch: AppDispatch;
}>();
export const createAppSelector = createSelector.withTypes<RootState>();
