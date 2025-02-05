import { Children, cloneElement, ReactElement, ReactNode } from "react";
import { StyleProp, View, ViewStyle } from "react-native";
import { Divider, Surface, Text, TouchableRipple } from "react-native-paper";

type ButtonListItemProps = {
  title: string;
  icon?: () => ReactNode;
  supportingText?: string;
  style?: StyleProp<ViewStyle>;
};

export function ButtonListItem({
  style,
  title,
  icon,
  supportingText,
}: ButtonListItemProps) {
  return (
    <Surface mode="flat" elevation={1} style={style}>
      <TouchableRipple style={style} borderless onPress={() => {}}>
        <View
          style={{ flexDirection: "row", padding: 16, alignItems: "center" }}
        >
          <View style={{ width: 24, height: 24 }}>{icon?.()}</View>
          <View style={{ marginLeft: 16, flex: 1 }}>
            <Text variant="bodyLarge">{title}</Text>
            {supportingText && (
              <Text variant="bodySmall">{supportingText}</Text>
            )}
          </View>
        </View>
      </TouchableRipple>
    </Surface>
  );
}

type ButtonListProps = {
  children:
    | ReactElement<{ style?: StyleProp<ViewStyle> }>
    | ReactElement<{ style?: StyleProp<ViewStyle> }>[];
};

export function ButtonList({ children }: ButtonListProps) {
  const totalChildren = Children.count(children);
  const childrenWithStyles = Children.map(children, (child, index) => {
    const toReturn = [];
    const isFirstElement = index === 0;
    const isLastElement = index === totalChildren - 1;

    let style: StyleProp<ViewStyle> = {};
    if (isFirstElement && isLastElement) {
      style = { borderRadius: 12 };
    } else if (isLastElement) {
      style = { borderBottomLeftRadius: 12, borderBottomRightRadius: 12 };
    } else if (index === 0) {
      style = { borderTopLeftRadius: 12, borderTopRightRadius: 12 };
    }

    toReturn.push(cloneElement(child, { style }));
    if (index < totalChildren - 1) {
      toReturn.push(<Divider />);
    }
    return toReturn;
  });
  return <View>{childrenWithStyles}</View>;
}
