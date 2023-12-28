import { extendTheme, type ThemeConfig } from '@chakra-ui/react';
import modalTheme from './modalTheme';
import cardTheme from './cardTheme';

const config: ThemeConfig = {
  initialColorMode: 'dark', // light, dark, system
  useSystemColorMode: false,
};

const theme = extendTheme({
  config,
  components: {
    Modal: modalTheme,
    Card: cardTheme,
  },
});

export default theme;
