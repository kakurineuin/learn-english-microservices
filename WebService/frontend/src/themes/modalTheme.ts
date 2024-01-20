import { modalAnatomy as parts } from '@chakra-ui/anatomy';
import { createMultiStyleConfigHelpers } from '@chakra-ui/styled-system';

const { definePartsStyle, defineMultiStyleConfig } =
  createMultiStyleConfigHelpers(parts.keys);

const baseStyle = definePartsStyle({
  // define the part you're going to style
  overlay: {
    bg: 'none', // change the background
    backdropFilter: 'auto',
    backdropBlur: '4px',
  },
  dialog: {
    borderRadius: 'lg',
    bg: `none`,
    borderStyle: 'solid',
    borderWidth: '1px',
    borderColor: 'gray',
  },
});

const modalTheme = defineMultiStyleConfig({
  baseStyle,
});

export default modalTheme;
