import { getContext, setContext } from "svelte";
import { afterNavigate, beforeNavigate } from "$app/navigation";

type ScreenConfigureOpts = {
  backUrl?: URL | string;
  fallbackBackUrl?: URL | string;
  replaceBack?: boolean;
};

export function configureScreen({
  backUrl,
  fallbackBackUrl,
  replaceBack = true,
}: ScreenConfigureOpts) {
  const { state } = getContext("screen") as ScreenContext;
  state.backUrl = backUrl;
  state.fallbackBackUrl = fallbackBackUrl;
  state.replaceBack = replaceBack;
  beforeNavigate(() => {
    state.fallbackBackUrl = undefined;
  });
}

type ScreenConfig = {
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
    replaceBack: true,
    fallbackBackUrl: undefined,
    backUrl: undefined,
    stack: [],
    popping: false,
  });

  const backUrlHref = $derived(
    typeof state.backUrl === "string" ? state.backUrl : state.backUrl?.href,
  );

  afterNavigate((navigation) => {
    if (!state.popping && navigation.from?.url !== undefined) {
      state.stack.push(navigation.from.url);
    }
    const adjustedStackLength = state.popping
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
