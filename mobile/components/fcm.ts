import messaging, {
  FirebaseMessagingTypes,
} from "@react-native-firebase/messaging";
import notifee from "@notifee/react-native";

type RemoteMessage = FirebaseMessagingTypes.RemoteMessage;

async function onMessageReceived(message: RemoteMessage) {
  const channelId = await notifee.createChannel({
    id: "default",
    name: "Default Channel",
  });

  await notifee.displayNotification({
    title: message.notification?.title,
    body: message.notification?.body,
    android: {
      channelId: channelId,
    },
  });
}

messaging().onMessage(onMessageReceived);
messaging().setBackgroundMessageHandler(onMessageReceived);
