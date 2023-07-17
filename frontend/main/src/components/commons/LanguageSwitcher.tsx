import { ActionIcon, Flex, Menu, Select } from "@mantine/core";
import { CSSProperties } from "react";
import { useTranslation } from "react-i18next";
import { RU, US } from "country-flag-icons/react/3x2";
import { IconLanguage } from "@tabler/icons-react";
interface LanguageSwitcherProps {
  style?: CSSProperties;
}

const LanguageSwitcher = ({ style }: LanguageSwitcherProps) => {
  const { i18n } = useTranslation();

  const handleLanguageChange = (language: string) => {
    void i18n.changeLanguage(language);
  };

  const languageOptions = [
    {
      label: "English",
      value: "en",
      icon: <US style={{ height: "14px" }} />,
    },
    { label: "Русский", value: "ru", icon: <RU style={{ height: "14px" }} /> },
  ];

  return (
    <ActionIcon style={style}>
      <IconLanguage size={18} />
      <Menu shadow="md" withArrow width={200}>
        <Menu.Target>
          <ActionIcon size="lg" title="Toggle language" variant="default">
            <IconLanguage size={19} />
          </ActionIcon>
        </Menu.Target>
        <Menu.Dropdown>
          <Menu.Label>Change language</Menu.Label>
          <Menu.Divider />
          {languageOptions.map(({ label, value, icon }) => (
            <Menu.Item
              style={{
                fontWeight: value === i18n.language ? "bold" : "initial",
              }}
              icon={icon}
              key={value}
              onClick={() => handleLanguageChange(value)}
            >
              {label}
            </Menu.Item>
          ))}
        </Menu.Dropdown>
      </Menu>
    </ActionIcon>
  );
};

export default LanguageSwitcher;
