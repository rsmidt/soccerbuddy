import { createAppSelector } from "@/store/custom";
import { useAppSelector } from "@/store";
import { Redirect } from "expo-router";
const selectInitialRouteName = createAppSelector((state) => state.auth.type, (type) => {
    switch (type) {
        case "authenticated":
            return "(tabs)";
        case "unauthenticated":
            return "login";
        default:
            return "(tabs)";
    }
});
export default function IndexRedirect() {
    const initialRouteName = useAppSelector(selectInitialRouteName);
    if (initialRouteName === "login") {
        return <Redirect href="/login"/>;
    }
    return null;
}
