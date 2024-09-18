import { LinkedAs } from "$lib/gen/soccerbuddy/person/v1/person_service_pb";

export function pbToRole(pb: string): string {
  switch (pb) {
    case "COACH":
      return "Trainer:in";
    case "PLAYER":
      return "Spieler:in";
    default:
      return "Unbekannt";
  }
}

export function pbToLinkedAs(pb: LinkedAs): string {
  switch (pb) {
    case LinkedAs.PARENT:
      return "Eltern";
    case LinkedAs.SELF:
      return "Selbst";
    default:
      return "Unbekannt";
  }
}
