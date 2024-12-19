// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// @generated by protoc-gen-es v2.2.3 with parameter "target=ts"
// @generated from file google/type/money.proto (package google.type, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc } from "@bufbuild/protobuf/codegenv1";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file google/type/money.proto.
 */
export const file_google_type_money: GenFile =
  /*@__PURE__*/
  fileDesc(
    "Chdnb29nbGUvdHlwZS9tb25leS5wcm90bxILZ29vZ2xlLnR5cGUiPAoFTW9uZXkSFQoNY3VycmVuY3lfY29kZRgBIAEoCRINCgV1bml0cxgCIAEoAxINCgVuYW5vcxgDIAEoBUKfAQoPY29tLmdvb2dsZS50eXBlQgpNb25leVByb3RvUAFaMGdpdGh1Yi5jb20vcnNtaWR0L3NvY2NlcmJ1ZGR5L2dlbi9nby9nb29nbGUvdHlwZfgBAaICA0dUWKoCC0dvb2dsZS5UeXBlygILR29vZ2xlXFR5cGXiAhdHb29nbGVcVHlwZVxHUEJNZXRhZGF0YeoCDEdvb2dsZTo6VHlwZWIGcHJvdG8z",
  );

/**
 * Represents an amount of money with its currency type.
 *
 * @generated from message google.type.Money
 */
export type Money = Message<"google.type.Money"> & {
  /**
   * The three-letter currency code defined in ISO 4217.
   *
   * @generated from field: string currency_code = 1;
   */
  currencyCode: string;

  /**
   * The whole units of the amount.
   * For example if `currencyCode` is `"USD"`, then 1 unit is one US dollar.
   *
   * @generated from field: int64 units = 2;
   */
  units: bigint;

  /**
   * Number of nano (10^-9) units of the amount.
   * The value must be between -999,999,999 and +999,999,999 inclusive.
   * If `units` is positive, `nanos` must be positive or zero.
   * If `units` is zero, `nanos` can be positive, zero, or negative.
   * If `units` is negative, `nanos` must be negative or zero.
   * For example $-1.75 is represented as `units`=-1 and `nanos`=-750,000,000.
   *
   * @generated from field: int32 nanos = 3;
   */
  nanos: number;
};

/**
 * Describes the message google.type.Money.
 * Use `create(MoneySchema)` to create a new message.
 */
export const MoneySchema: GenMessage<Money> = /*@__PURE__*/
  messageDesc(file_google_type_money, 0);

