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
import { ReactComponent as Logo } from "./../assets/logo.svg";
import { IconArrowRight } from "@tabler/icons-react";
import { useNavigate } from "react-router-dom";

const Login = () => {
  const navigate = useNavigate();

  const form = useForm({
    initialValues: {
      email: "",
      password: "",
    },

    validate: {
      email: (value) => (/^\S+@\S+$/.test(value) ? null : "Invalid email"),
      password: (value) =>
        value.length > 6 ? null : "Password must be more than 6 symbols length",
    },
  });

  const onLoginClick = (values: { email: string; password: string }) => {};

  return (
    <Flex h={"100vh"}>
      <Box w={"50%"} pt="xs" px="lg">
        <Logo />
        <Group w={"100%"} position="center" mt="xl">
          <form onSubmit={form.onSubmit((values) => onLoginClick(values))}>
            <Title order={3} ta="center">
              Sign in
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
                Login
              </Button>
              <Text onClick={() => navigate("reset")}>
                Forgot your password?
              </Text>
            </Group>
            <Button color="white" bg="black">
              CREATE NEW ACCOUNT
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
