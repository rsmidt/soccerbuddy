import { AccountLink } from "$lib/gen/soccerbuddy/shared_pb";

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

export function pbToAccountLink(pb: AccountLink): string {
  switch (pb) {
    case AccountLink.LINKED_AS_PARENT:
      return "Eltern";
    case AccountLink.LINKED_AS_SELF:
      return "Selbst";
    default:
      return "Unbekannt";
  }
}
