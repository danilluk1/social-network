import LanguageSwitcher from "../components/commons/LanguageSwitcher";
import { Box, Button, Flex, Title, Text, TextInput } from "@mantine/core";
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
//@ts-ignore
import { ReactComponent as Logo } from "./../assets/logo.svg";
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
//@ts-ignore
import { ReactComponent as About } from "./../assets/about.svg";
import ThemeSwitcher from "../components/commons/ThemeSwitcher";
import { useTranslation } from "react-i18next";
import { useForm } from "@mantine/form";

const Confirm = () => {
  const { t } = useTranslation("confirm");

  const form = useForm({
    initialValues: {
      code: "",
    },

    validate: {
      code: (value) => (value.length != 6 ? null : t("invalidCode")),
    },
  });

  return (
    <Flex h={"100vh"} display={"flex"} pt="lg" px="lg">
      <Logo
        style={{
          width: "2.3em",
          height: "2.3em",
          position: "absolute",
          top: "20px",
          left: "20px",
        }}
      />
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
      <Flex w={"100%"} h={"70vh"} mt="lg" align="center" justify="center">
        <form>
          <Title>{t("confirmEmail")}</Title>
          <Text>{t("pleaseCheck")}</Text>
          <TextInput mb={"sm"} {...form.getInputProps("code")} />
          <Button type="submit">{t("send")}</Button>
        </form>
      </Flex>
    </Flex>
  );
};

export default Confirm;
