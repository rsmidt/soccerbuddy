import { configureStore } from "@reduxjs/toolkit";
import { useDispatch, useSelector } from "react-redux";
import { reducer as authReducer } from "@/components/auth/auth-slice";
import devToolsEnhancer from "redux-devtools-expo-dev-plugin";
import { accountApi } from "@/components/account/account-api";
import { setupListeners } from "@reduxjs/toolkit/query";
// Required for Redux DevTools to serialize BigInts from ConnectRPC.
// @ts-ignore
// eslint-disable-next-line no-extend-native
BigInt.prototype.toJSON = function () {
    return this.toString();
};
export const store = configureStore({
    reducer: {
        auth: authReducer,
        [accountApi.reducerPath]: accountApi.reducer,
    },
    devTools: false,
    enhancers: (getDefaultEnhancers) => getDefaultEnhancers().concat(devToolsEnhancer()),
    // @ts-ignore Some redux internal clash with thunk signatures...
    middleware: (getDefaultMiddleware) => getDefaultMiddleware({
        serializableCheck: {
            // ConnectRPC sends BigInts for these values...
            ignoredActionPaths: [/\.seconds/, /\.nanos/],
            ignoredPaths: [/\.seconds/, /\.nanos/],
        },
    }).concat(accountApi.middleware),
});
setupListeners(store.dispatch);
// Use throughout your app instead of plain `useDispatch` and `useSelector`.
export const useAppDispatch = useDispatch.withTypes();
export const useAppSelector = useSelector.withTypes();
