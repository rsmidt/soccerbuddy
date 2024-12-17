import { createAsyncThunk } from "@reduxjs/toolkit";
import { createSelector } from "reselect";
export const createAppAsyncThunk = createAsyncThunk.withTypes();
export const createAppSelector = createSelector.withTypes();
