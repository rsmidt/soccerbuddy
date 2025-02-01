import { getContext, setContext } from "svelte";
import { afterNavigate, beforeNavigate } from "$app/navigation";
import { contentTypeUnaryRegExp } from "@connectrpc/connect/protocol-connect";

type ScreenConfigureOpts = {
  backUrl?: URL | string;
  fallbackBackUrl?: URL | string;
  replaceBack?: boolean;
  backButtonShown?: boolean;
};

/**
 * Allows configuring the currently active screen (aka route).
 *
 * @param backUrl - Hard overwrite of the back button url.
 * @param fallbackBackUrl - Fallback url for the back button.
 * @param backButtonHidden - If the back button should be hidden.
 * @param replaceBack - If the back button should replace the current history entry.
 */
export function configureScreen({
  backUrl,
  fallbackBackUrl,
  backButtonShown,
  replaceBack = true,
}: ScreenConfigureOpts) {
  const { state } = getContext("screen") as ScreenContext;
  const currentFallbackBackUrl = state.fallbackBackUrl;
  const currentBackButtonShown = state.backButtonShown;
  state.backUrl = backUrl;
  state.fallbackBackUrl = fallbackBackUrl;
  state.replaceBack = replaceBack;
  if (backButtonShown !== undefined) {
    state.backButtonShown = backButtonShown;
  }
  beforeNavigate(() => {
    state.fallbackBackUrl = currentFallbackBackUrl;
    state.backButtonShown = currentBackButtonShown;
  });
}

type ScreenConfig = {
  backButtonShown: boolean;
  backUrl?: URL | string;
  fallbackBackUrl?: URL | string;
  replaceBack: boolean;
  stack: URL[];
  popping: boolean;
};

type ScreenContext = {
  state: ScreenConfig;
  backUrlHref: string | undefined;
};

export function initializeScreenContext(): ScreenContext {
  const state = $state<ScreenConfig>({
    backButtonShown: true,
    replaceBack: true,
    fallbackBackUrl: ".",
    backUrl: undefined,
    stack: [],
    popping: false,
  });

  const backUrlHref = $derived(
    typeof state.backUrl === "string" ? state.backUrl : state.backUrl?.href,
  );

  afterNavigate((navigation) => {
    if (navigation.type === "popstate") {
      state.stack.pop();
    }
    const isPopping = state.popping || navigation.type === "popstate";
    if (!isPopping && navigation.from?.url !== undefined) {
      state.stack.push(navigation.from.url);
    }
    const adjustedStackLength = isPopping
      ? Math.max(state.stack.length - 1, 0)
      : state.stack.length;
    if (adjustedStackLength === 0) {
      state.backUrl = state.fallbackBackUrl;
    } else {
      state.backUrl = state.stack[adjustedStackLength - 1];
    }
    return;
  });

  const context: ScreenContext = {
    get state() {
      return state;
    },
    get backUrlHref() {
      return backUrlHref;
    },
  };
  setContext("screen", context);
  return context;
}
