import { createAsyncThunk } from "@reduxjs/toolkit";
import { AppDispatch, RootState } from "./index";
import { createSelector } from "reselect";
import { useDispatch, useSelector } from "react-redux";

export const createAppAsyncThunk = createAsyncThunk.withTypes<{
  state: RootState;
  dispatch: AppDispatch;
}>();
export const createAppSelector = createSelector.withTypes<RootState>();

// Use throughout your app instead of plain `useDispatch` and `useSelector`.
export const useAppDispatch = useDispatch.withTypes<AppDispatch>();
export const useAppSelector = useSelector.withTypes<RootState>();
