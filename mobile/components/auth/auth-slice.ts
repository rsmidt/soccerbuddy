import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { ConnectError, createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import {
  AccountService,
  GetMeResponse,
} from "@/api/soccerbuddy/account/v1/account_service_pb";
import { createAppAsyncThunk } from "@/store/custom";
import { INSTALLATION_ID_KEY, SESSION_TOKEN_KEY } from "./constants";
import * as SecureStore from "expo-secure-store";
import messaging from "@react-native-firebase/messaging";
import { PermissionsAndroid, Platform } from "react-native";

export type AuthState =
  | {
      type: "unresolved";
    }
  | {
      type: "unauthenticated";
    }
  | {
      type: "authenticated";
      token: string;
      user: {
        id: string;
      };
    }
  | {
      type: "pending";
      token: string;
    };

const initialState: AuthState = {
  type: "unauthenticated",
};

const authSlice = createSlice({
  name: "auth",
  initialState: initialState as AuthState,
  reducers: {
    setPending: (state, action: PayloadAction<{ token: string }>) => {
      switch (state.type) {
        case "unauthenticated":
          return {
            type: "pending",
            token: action.payload.token,
          };
        default:
          return state;
      }
    },
    setUnauthenticated: (state) => {
      return {
        type: "unauthenticated",
      };
    },
    logout: (state) => {
      return {
        type: "unauthenticated",
      };
    },
  },
  extraReducers: (builder) => {
    builder.addCase(validateStoredToken.fulfilled, (state, action) => {
      if (state.type !== "pending") {
        return state;
      }
      return {
        type: "authenticated",
        token: state.token,
        user: {
          id: action.payload.id,
        },
      };
    });
    builder.addCase(loginUser.fulfilled, (state, action) => {
      if (state.type !== "pending") {
        return state;
      }
      return {
        type: "authenticated",
        token: state.token,
        user: {
          id: action.payload.id,
        },
      };
    });
  },
});

export const validateStoredToken = createAppAsyncThunk(
  "auth/fetchMe",
  async (token: string, thunkAPI) => {
    thunkAPI.dispatch(setPending({ token }));

    const client = createClient(
      AccountService,
      createConnectTransport({
        baseUrl: process.env.EXPO_PUBLIC_API_URL!,
        interceptors: [
          (next) => async (req) => {
            req.header.set("Authorization", `Bearer ${token}`);
            return next(req);
          },
        ],
      }),
    );
    const me = (await client.getMe({})) as GetMeResponse;

    const info = await getDeviceInfo();
    if (info) {
      await client.attachMobileDevice({
        deviceNotificationToken: info.deviceToken,
        installationId: info.installationId,
      });
    }

    return me;
  },
);

// TODO: Bad, bad, bad. Refactor this.
async function getDeviceInfo(): Promise<{
  deviceToken: string;
  installationId: string;
} | null> {
  const notificationEnabled = await requestUserPermissions();
  if (!notificationEnabled) {
    // TODO: Probably, we need to store this information somewhere.
    console.log("Notifications are not enabled.");
  } else {
    const installationId = await fetchOrCreateInstallationId();
    await messaging().registerDeviceForRemoteMessages();
    const deviceToken = await messaging().getToken();
    return { deviceToken, installationId };
  }
  return null;
}

async function requestUserPermissions(): Promise<boolean> {
  if (Platform.OS === "ios") {
    const authStatus = await messaging().requestPermission();
    return (
      authStatus === messaging.AuthorizationStatus.AUTHORIZED ||
      authStatus === messaging.AuthorizationStatus.PROVISIONAL
    );
  } else if (Platform.OS === "android") {
    return (
      (await PermissionsAndroid.request(
        PermissionsAndroid.PERMISSIONS.POST_NOTIFICATIONS,
      )) === "granted"
    );
  }
  return false;
}

async function fetchOrCreateInstallationId(): Promise<string> {
  let installationId = await SecureStore.getItemAsync(INSTALLATION_ID_KEY);
  if (!installationId) {
    installationId = Math.random().toString();
    await SecureStore.setItemAsync(INSTALLATION_ID_KEY, installationId);
  }
  return installationId;
}

export const loginUser = createAppAsyncThunk(
  "auth/loginUser",
  async (credentials: { email: string; password: string }, thunkAPI) => {
    const client = createClient(
      AccountService,
      createConnectTransport({
        baseUrl: process.env.EXPO_PUBLIC_API_URL!,
      }),
    );

    let sessionId: string;
    try {
      const { sessionId: sId } = await client.login({
        email: credentials.email,
        password: credentials.password,
      });
      sessionId = sId;
    } catch (e) {
      const cErr = ConnectError.from(e);
      console.error("failed to login", cErr);
      return thunkAPI.rejectWithValue(cErr.code);
    }

    thunkAPI.dispatch(setPending({ token: sessionId }));

    const headers = { headers: { Authorization: `Bearer ${sessionId}` } };
    const result = await client.getMe({}, headers);
    await SecureStore.setItemAsync(SESSION_TOKEN_KEY, sessionId, {
      requireAuthentication: false,
    });

    const info = await getDeviceInfo();
    if (info) {
      try {
        await client.attachMobileDevice(
          {
            deviceNotificationToken: info.deviceToken,
            installationId: info.installationId,
          },
          headers,
        );
      } catch (e) {
        const cErr = ConnectError.from(e);
        console.error("failed to attach device", cErr);
        return thunkAPI.rejectWithValue(cErr.code);
      }
    }

    return result as GetMeResponse;
  },
);

export const { actions, reducer } = authSlice;

export const { setPending, logout, setUnauthenticated } = actions;
