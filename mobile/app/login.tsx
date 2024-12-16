import { StyleSheet, View } from "react-native";
import { Button, Text, TextInput } from "react-native-paper";
import { Controller, useForm } from "react-hook-form";
import i18n from "@/components/i18n";
import { useAppDispatch, useAppSelector } from "@/store";
import { loginUser } from "@/components/auth/auth-slice";
import { Redirect } from "expo-router";

type LoginForm = {
  email: string;
  password: string;
};

export default function Login() {
  const dispatch = useAppDispatch();
  const authState = useAppSelector((state) => state.auth);

  const {
    formState: { errors },
    handleSubmit,
    control,
  } = useForm<LoginForm>();

  const onSubmit = async (data: LoginForm) => {
    dispatch(loginUser({ email: data.email, password: data.password }));
  };

  if (authState.type === "authenticated") {
    return <Redirect href="/(tabs)" />;
  }

  return (
    <View style={styles.container}>
      <View>
        <Controller
          control={control}
          rules={{
            required: true,
          }}
          render={({ field: { onChange, onBlur, value } }) => (
            <TextInput
              mode="outlined"
              placeholder={i18n.t("app.login.form.email.placeholder")}
              label={i18n.t("app.login.form.email.label")}
              onBlur={onBlur}
              onChangeText={onChange}
              value={value}
            />
          )}
          name="email"
        />
        {errors.email?.type === "required" && (
          <Text>{i18n.t("app.validation.required")}</Text>
        )}
      </View>

      <View>
        <Controller
          control={control}
          rules={{
            maxLength: 100,
          }}
          render={({ field: { onChange, onBlur, value } }) => (
            <TextInput
              mode="outlined"
              placeholder={i18n.t("app.login.form.password.placeholder")}
              label={i18n.t("app.login.form.password.label")}
              onBlur={onBlur}
              onChangeText={onChange}
              value={value}
            />
          )}
          name="password"
        />
        {errors.password?.type === "required" && (
          <Text>{i18n.t("app.validation.required")}</Text>
        )}
      </View>

      <Button mode="contained" onPress={handleSubmit(onSubmit)}>
        {i18n.t("app.login.form.submit")}
      </Button>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    paddingHorizontal: 16,
    marginTop: "auto",
    marginBottom: "auto",
    justifyContent: "center",
    gap: 16,
  },
});
