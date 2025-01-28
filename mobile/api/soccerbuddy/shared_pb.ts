// @generated by protoc-gen-es v2.2.3 with parameter "target=ts"
// @generated from file soccerbuddy/shared.proto (package soccerbuddy.shared, syntax proto3)
/* eslint-disable */

import type { GenEnum, GenFile } from "@bufbuild/protobuf/codegenv1";
import { enumDesc, fileDesc } from "@bufbuild/protobuf/codegenv1";

/**
 * Describes the file soccerbuddy/shared.proto.
 */
export const file_soccerbuddy_shared: GenFile = /*@__PURE__*/
  fileDesc("Chhzb2NjZXJidWRkeS9zaGFyZWQucHJvdG8SEnNvY2NlcmJ1ZGR5LnNoYXJlZCpSCgtBY2NvdW50TGluaxIZChVMSU5LRURfQVNfVU5TUEVDSUZJRUQQABISCg5MSU5LRURfQVNfU0VMRhABEhQKEExJTktFRF9BU19QQVJFTlQQAiqBAQoMUmF0aW5nUG9saWN5Eh0KGVJBVElOR19QT0xJQ1lfVU5TUEVDSUZJRUQQABIbChdSQVRJTkdfUE9MSUNZX0ZPUkJJRERFThABEhkKFVJBVElOR19QT0xJQ1lfQUxMT1dFRBACEhoKFlJBVElOR19QT0xJQ1lfUkVRVUlSRUQQA0LAAQoWY29tLnNvY2NlcmJ1ZGR5LnNoYXJlZEILU2hhcmVkUHJvdG9QAVowZ2l0aHViLmNvbS9yc21pZHQvc29jY2VyYnVkZHkvZ2VuL2dvL3NvY2NlcmJ1ZGR5ogIDU1NYqgISU29jY2VyYnVkZHkuU2hhcmVkygISU29jY2VyYnVkZHlcU2hhcmVk4gIeU29jY2VyYnVkZHlcU2hhcmVkXEdQQk1ldGFkYXRh6gITU29jY2VyYnVkZHk6OlNoYXJlZGIGcHJvdG8z");

/**
 * @generated from enum soccerbuddy.shared.AccountLink
 */
export enum AccountLink {
  /**
   * @generated from enum value: LINKED_AS_UNSPECIFIED = 0;
   */
  LINKED_AS_UNSPECIFIED = 0,

  /**
   * @generated from enum value: LINKED_AS_SELF = 1;
   */
  LINKED_AS_SELF = 1,

  /**
   * @generated from enum value: LINKED_AS_PARENT = 2;
   */
  LINKED_AS_PARENT = 2,
}

/**
 * Describes the enum soccerbuddy.shared.AccountLink.
 */
export const AccountLinkSchema: GenEnum<AccountLink> = /*@__PURE__*/
  enumDesc(file_soccerbuddy_shared, 0);

/**
 * @generated from enum soccerbuddy.shared.RatingPolicy
 */
export enum RatingPolicy {
  /**
   * @generated from enum value: RATING_POLICY_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * @generated from enum value: RATING_POLICY_FORBIDDEN = 1;
   */
  FORBIDDEN = 1,

  /**
   * @generated from enum value: RATING_POLICY_ALLOWED = 2;
   */
  ALLOWED = 2,

  /**
   * @generated from enum value: RATING_POLICY_REQUIRED = 3;
   */
  REQUIRED = 3,
}

/**
 * Describes the enum soccerbuddy.shared.RatingPolicy.
 */
export const RatingPolicySchema: GenEnum<RatingPolicy> = /*@__PURE__*/
  enumDesc(file_soccerbuddy_shared, 1);

