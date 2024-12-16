import { configureStore } from "@reduxjs/toolkit";
import { useDispatch, useSelector } from "react-redux";
import { reducer as authReducer } from "@/components/auth/auth-slice";
import devToolsEnhancer from "redux-devtools-expo-dev-plugin";
import { accountApi } from "@/components/account/account-api";

export const store = configureStore({
  reducer: {
    auth: authReducer,
    [accountApi.reducerPath]: accountApi.reducer,
  },
  devTools: false,
  enhancers: (getDefaultEnhancers) =>
    getDefaultEnhancers().concat(devToolsEnhancer()),
  // @ts-ignore Some redux internal clash with thunk signatures...
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(accountApi.middleware),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

// Use throughout your app instead of plain `useDispatch` and `useSelector`.
export const useAppDispatch = useDispatch.withTypes<AppDispatch>();
export const useAppSelector = useSelector.withTypes<RootState>();
