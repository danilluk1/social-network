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
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
//@ts-ignore
import { ReactComponent as About } from "./../assets/about.svg";
import { IconArrowRight } from "@tabler/icons-react";
import { useNavigate } from "react-router-dom";
import ThemeSwitcher from "../components/commons/ThemeSwitcher";
import { useMediaQuery } from "@mantine/hooks";
import LanguageSwitcher from "../components/commons/LanguageSwitcher";
import { gql, useMutation } from "@apollo/client";

const REGISTER_USER = gql`
  mutation createUser(input: {

  })
`;
const Register = () => {
  const matches = useMediaQuery("(min-width: 640px)");
  const { t } = useTranslation("register");
  const navigate = useNavigate();
  // const;

  const form = useForm({
    initialValues: {
      username: "",
      full_name: "",
      email: "",
      password: "",
    },

    validate: {
      email: (value) => (/^\S+@\S+$/.test(value) ? null : t("invalidEmail")),
      password: (value) => (value.length >= 6 ? null : t("invalidPassword")),
    },
  });

  const onRegisterClick = (values: {
    email: string;
    password: string;
    username: string;
    full_name: string;
  }) => {};

  return (
    <Flex h={"100vh"} display={"flex"}>
      <Box w={matches ? "50%" : "100%"} pt="lg" px="lg" pos="relative">
        <Logo />
        <LanguageSwitcher
          style={{
            width: "2.3em",
            height: "2.3em",
            position: "absolute",
            bottom: "60px",
            right: "20px",
          }}
        />
        <ThemeSwitcher
          style={{
            width: "2.2em",
            height: "2.2em",
            position: "absolute",
            bottom: "20px",
            right: "20px",
          }}
        />
        <Flex w={"100%"} h={"70vh"} align="center" justify="center">
          <form onSubmit={form.onSubmit((values) => onRegisterClick(values))}>
            <Title order={3} ta="center">
              {t("signIn")}
            </Title>
            <TextInput
              mb={"sm"}
              withAsterisk
              label="Email"
              {...form.getInputProps("email")}
            />
            <PasswordInput
              mb={"sm"}
              withAsterisk
              label="Password"
              {...form.getInputProps("password")}
            />
            <TextInput
              mb={"sm"}
              withAsterisk
              label="Username"
              {...form.getInputProps("username")}
            />
            <TextInput
              mb={"sm"}
              withAsterisk
              label="Full name"
              {...form.getInputProps("full_name")}
            />
            <Group mb={"lg"}>
              <Button type="submit" rightIcon={<IconArrowRight />}>
                {t("signUp")}
              </Button>
            </Group>
            <Button color="white" bg="black" onClick={() => navigate("/login")}>
              {t("signIn").toLocaleUpperCase()}
            </Button>
          </form>
        </Flex>
        <Group></Group>
      </Box>
      {matches ? (
        <Flex bg="#FAFAFB" w={"50%"} align="center" justify="center">
          <About />
        </Flex>
      ) : (
        <></>
      )}
    </Flex>
  );
};

export default Register;
