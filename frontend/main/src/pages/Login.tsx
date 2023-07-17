import {
  Flex,
  Text,
  Group,
  TextInput,
  PasswordInput,
  Button,
  Box,
  Title,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { useTranslation } from "react-i18next";
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
//@ts-ignore
import { ReactComponent as Logo } from "./../assets/logo.svg";
import { IconArrowRight } from "@tabler/icons-react";
import { useNavigate } from "react-router-dom";
import ThemeSwitcher from "../components/commons/ThemeSwitcher";

const Login = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const form = useForm({
    initialValues: {
      email: "",
      password: "",
    },

    validate: {
      email: (value) => (/^\S+@\S+$/.test(value) ? null : t("invalidEmail")),
      password: (value) => (value.length > 6 ? null : t("invalidPassword")),
    },
  });

  const onLoginClick = (values: { email: string; password: string }) => {};
  return (
    <Flex h={"100vh"}>
      <Box w={"50%"} pt="xs" px="lg" pos="relative">
        <Logo />
        <ThemeSwitcher
          style={{
            width: "2.2em",
            height: "2.2em",
            position: "absolute",
            bottom: "20px",
            right: "20px",
          }}
        />
        <Group w={"100%"} position="center" mt="xl">
          <form onSubmit={form.onSubmit((values) => onLoginClick(values))}>
            <Title order={3} ta="center">
              {t("signIn")}
            </Title>
            <TextInput
              mb={"sm"}
              withAsterisk
              label="Email"
              placeholder="your@email.com"
              {...form.getInputProps("email")}
            />
            <PasswordInput
              mb={"sm"}
              withAsterisk
              label="Password"
              placeholder="your@email.com"
              {...form.getInputProps("password")}
            />
            <Group mb={"lg"}>
              <Button type="submit" rightIcon={<IconArrowRight />}>
                {t("signIn")}
              </Button>
              <Text onClick={() => navigate("reset")}>
                {t("forgotPassword")}
              </Text>
            </Group>
            <Button color="white" bg="black">
              {t("createNewAccount").toLocaleUpperCase()}
            </Button>
          </form>
        </Group>
        <Group></Group>
      </Box>
      <Flex w={"50%"} bg="#FAFAFB">
        123
      </Flex>
    </Flex>
  );
};

export default Login;
