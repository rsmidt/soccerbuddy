import { configureStore } from "@reduxjs/toolkit";
import { useDispatch, useSelector } from "react-redux";
import { reducer as authReducer } from "@/components/auth/auth-slice";
import { reducer as teamReducer } from "@/components/team/team-slice";
import devToolsEnhancer from "redux-devtools-expo-dev-plugin";
import { accountApi } from "@/components/account/account-api";
import { setupListeners } from "@reduxjs/toolkit/query";
import { teamApi } from "@/components/team/team-api";
import {
  FLUSH,
  PAUSE,
  PERSIST,
  persistReducer,
  persistStore,
  PURGE,
  REGISTER,
  REHYDRATE,
} from "redux-persist";
import AsyncStorage from "@react-native-async-storage/async-storage";

// Required for Redux DevTools to serialize BigInts from ConnectRPC.
// @ts-ignore
// eslint-disable-next-line no-extend-native
BigInt.prototype.toJSON = function () {
  return this.toString();
};

// For now, persist only the team config.
const teamPersistConfig = {
  key: "team",
  storage: AsyncStorage,
};

const persistedTeamReducer = persistReducer(teamPersistConfig, teamReducer);

export const store = configureStore({
  reducer: {
    auth: authReducer,
    team: persistedTeamReducer,
    [accountApi.reducerPath]: accountApi.reducer,
    [teamApi.reducerPath]: teamApi.reducer,
  },
  devTools: false,
  enhancers: (getDefaultEnhancers) =>
    getDefaultEnhancers().concat(devToolsEnhancer()),
  // @ts-ignore Some redux internal clash with thunk signatures...
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        // ConnectRPC sends BigInts for these values...
        ignoredActionPaths: [/\.seconds/, /\.nanos/],
        ignoredPaths: [/\.seconds/, /\.nanos/],
        ignoredActions: [FLUSH, REHYDRATE, PAUSE, PERSIST, PURGE, REGISTER],
      },
    }).concat(accountApi.middleware, teamApi.middleware),
});

setupListeners(store.dispatch);

export const persistor = persistStore(store);

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

// Use throughout your app instead of plain `useDispatch` and `useSelector`.
export const useAppDispatch = useDispatch.withTypes<AppDispatch>();
export const useAppSelector = useSelector.withTypes<RootState>();
