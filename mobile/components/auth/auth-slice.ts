import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { ConnectError, createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import {
  AccountService,
  GetMeResponse,
} from "@/api/soccerbuddy/account/v1/account_service_pb";
import { createAppAsyncThunk } from "@/store/custom";
import { SESSION_TOKEN_KEY } from "./constants";
import * as SecureStore from "expo-secure-store";

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
    builder.addCase(fetchMe.fulfilled, (state, action) => {
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

export const fetchMe = createAppAsyncThunk(
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

    return (await client.getMe({})) as GetMeResponse;
  },
);

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
      console.log(cErr);
      return thunkAPI.rejectWithValue(cErr.code);
    }

    thunkAPI.dispatch(setPending({ token: sessionId }));

    const result = await client.getMe(
      {},
      { headers: { Authorization: `Bearer ${sessionId}` } },
    );
    await SecureStore.setItemAsync(SESSION_TOKEN_KEY, sessionId, {
      requireAuthentication: false,
    });

    return result as GetMeResponse;
  },
);

export const { actions, reducer } = authSlice;

export const { setPending, logout, setUnauthenticated } = actions;
