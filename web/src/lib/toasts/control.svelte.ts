import type { Component } from "svelte";
import { getContext, setContext } from "svelte";

export type Toast<TComp extends Component<TProps>, TProps extends Record<string, any>> = {
  id?: string;
  component: TComp;
  props: TProps;
  timeout?: number;
};

export type IdedToast<TComp extends Component<TProps>, TProps extends Record<string, any>> = Toast<
  TComp,
  TProps
> & { id: string };

export class ToastControl {
  private _toasts = $state<IdedToast<any, any>[]>([]);

  private constructor() {}

  static createGlobal(): ToastControl {
    const control = new ToastControl();
    setContext("toastControl", control);
    return control;
  }

  static getGlobal(): ToastControl {
    return getContext("toastControl") as ToastControl;
  }

  add<TComp extends Component<TProps>, TProps extends Record<string, any>>(
    toast: Toast<TComp, TProps>,
  ) {
    if (!toast.id) {
      toast.id = Math.random().toString(36);
    }
    this._toasts.push(toast as IdedToast<TComp, TProps>);
  }

  remove(toastId: string) {
    this._toasts = this._toasts.filter((t) => t.id !== toastId);
  }

  clear() {
    this._toasts = [];
  }

  get toasts(): ReadonlyArray<IdedToast<any, any>> {
    return this._toasts;
  }
}
