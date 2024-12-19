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
// @generated from file google/type/timeofday.proto (package google.type, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc } from "@bufbuild/protobuf/codegenv1";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file google/type/timeofday.proto.
 */
export const file_google_type_timeofday: GenFile = /*@__PURE__*/
  fileDesc("Chtnb29nbGUvdHlwZS90aW1lb2ZkYXkucHJvdG8SC2dvb2dsZS50eXBlIksKCVRpbWVPZkRheRINCgVob3VycxgBIAEoBRIPCgdtaW51dGVzGAIgASgFEg8KB3NlY29uZHMYAyABKAUSDQoFbmFub3MYBCABKAVCowEKD2NvbS5nb29nbGUudHlwZUIOVGltZW9mZGF5UHJvdG9QAVowZ2l0aHViLmNvbS9yc21pZHQvc29jY2VyYnVkZHkvZ2VuL2dvL2dvb2dsZS90eXBl+AEBogIDR1RYqgILR29vZ2xlLlR5cGXKAgtHb29nbGVcVHlwZeICF0dvb2dsZVxUeXBlXEdQQk1ldGFkYXRh6gIMR29vZ2xlOjpUeXBlYgZwcm90bzM");

/**
 * Represents a time of day. The date and time zone are either not significant
 * or are specified elsewhere. An API may choose to allow leap seconds. Related
 * types are [google.type.Date][google.type.Date] and
 * `google.protobuf.Timestamp`.
 *
 * @generated from message google.type.TimeOfDay
 */
export type TimeOfDay = Message<"google.type.TimeOfDay"> & {
  /**
   * Hours of day in 24 hour format. Should be from 0 to 23. An API may choose
   * to allow the value "24:00:00" for scenarios like business closing time.
   *
   * @generated from field: int32 hours = 1;
   */
  hours: number;

  /**
   * Minutes of hour of day. Must be from 0 to 59.
   *
   * @generated from field: int32 minutes = 2;
   */
  minutes: number;

  /**
   * Seconds of minutes of the time. Must normally be from 0 to 59. An API may
   * allow the value 60 if it allows leap-seconds.
   *
   * @generated from field: int32 seconds = 3;
   */
  seconds: number;

  /**
   * Fractions of seconds in nanoseconds. Must be from 0 to 999,999,999.
   *
   * @generated from field: int32 nanos = 4;
   */
  nanos: number;
};

/**
 * Describes the message google.type.TimeOfDay.
 * Use `create(TimeOfDaySchema)` to create a new message.
 */
export const TimeOfDaySchema: GenMessage<TimeOfDay> = /*@__PURE__*/
  messageDesc(file_google_type_timeofday, 0);

