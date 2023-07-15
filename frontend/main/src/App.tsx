import { MantineProvider, useMantineTheme, ColorScheme } from "@mantine/core";
import Router from "./components/router/Router";

function App() {
  const theme = useMantineTheme();

  return (
    <MantineProvider withGlobalStyles withNormalizeCSS>
      <Router />
    </MantineProvider>
  );
}

export default App;
