import { type ActionFailure, error, fail, redirect } from "@sveltejs/kit";
import { setError, type SuperValidated } from "sveltekit-superforms";
import { Code, ConnectError } from "@connectrpc/connect";
import { BadRequest } from "$lib/gen/google/rpc/error_details_pb";

type ValidationHandler = (fieldErrors: Record<string, string>) => ActionFailure<any>;
type UnauthenticatedHandler<TResponse> = () => TResponse;

export class GrpcMutationHandler<TResponse, TUnauthResponse> {
  private _validationHandler: ValidationHandler | undefined;
  private _unauthenticatedHandler: UnauthenticatedHandler<TUnauthResponse> | undefined = () => {
    redirect(302, "/login");
  };

  static from<TResponse>(fn: () => Promise<() => TResponse>) {
    return new GrpcMutationHandler(fn);
  }

  static overwriteFormHandler(form: SuperValidated<any>): ValidationHandler {
    return (fieldErrors) => {
      Object.entries(fieldErrors).forEach(([field, error]) => {
        setError(form, field, error);
      });
      return fail(400, { form });
    };
  }

  constructor(private readonly fn: () => Promise<() => TResponse>) {}

  onFailedValidation(cb: ValidationHandler) {
    this._validationHandler = cb;
    return this;
  }

  onUnauthenticated(cb: UnauthenticatedHandler<TUnauthResponse>) {
    this._unauthenticatedHandler = cb;
    return this;
  }

  async run(): Promise<TResponse | ActionFailure<any> | TUnauthResponse> {
    let finalizer;
    try {
      finalizer = await this.fn();
    } catch (err) {
      const cErr = ConnectError.from(err);

      switch (cErr.code) {
        case Code.InvalidArgument: {
          const violations = cErr.findDetails(BadRequest).reduce(
            (prev, br) => {
              for (const fv of br.fieldViolations) {
                prev[fv.field] = fv.description;
              }
              return prev;
            },
            {} as Record<string, string>,
          );
          if (Object.keys(violations).length === 0) {
            return error(400, { message: cErr.rawMessage });
          }
          if (this._validationHandler) {
            return this._validationHandler(violations);
          }
          return fail(400, { form: { valid: false, errors: violations } });
        }
        case Code.Unauthenticated:
          if (this._unauthenticatedHandler) {
            return this._unauthenticatedHandler();
          }
          return error(401, { message: cErr.rawMessage });
        default:
          return error(500, { message: cErr.rawMessage });
      }
    }
    return finalizer();
  }
}
