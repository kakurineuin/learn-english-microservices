import { Box } from '@chakra-ui/react';
import { useEffect, useState } from 'react';
import styles from './ShowText.module.css';

type As =
  | 'b'
  | 'i'
  | 'u'
  | 'abbr'
  | 'cite'
  | 'del'
  | 'em'
  | 'ins'
  | 'kbd'
  | 'mark'
  | 's'
  | 'samp'
  | 'sub'
  | 'sup';

type FontSize =
  | 'xs'
  | 'sm'
  | 'md'
  | 'lg'
  | 'xl'
  | '2xl'
  | '3xl'
  | '4xl'
  | '5xl'
  | '6xl';

type Props = {
  children: string;
  as?: As;
  fontSize?: FontSize;
  color?: string;
  mt?: string;
  mb?: string;
  ml?: string;
  mr?: string;
  mx?: string;
  my?: string;
};

function ShowText({
  children,
  as,
  fontSize,
  color,
  mt,
  mb,
  ml,
  mr,
  mx,
  my,
}: Props) {
  const [chars, setChars] = useState(children.split(''));
  const [showText, setShowText] = useState('');

  useEffect(() => {
    if (chars.length === 0) {
      return;
    }

    setTimeout(() => {
      setChars((prevChars) => {
        if (prevChars.length === 0) {
          return prevChars;
        }

        const char = prevChars.shift();
        setShowText((prevShowText) => `${prevShowText}${char}`);
        return [...prevChars];
      });
    }, 20);
  }, [chars]);

  return (
    <Box
      as={as}
      fontSize={fontSize}
      color={color}
      mt={mt}
      mb={mb}
      ml={ml}
      mr={mr}
      mx={mx}
      my={my}
    >
      <pre className={styles['show-text']}>{showText}</pre>
    </Box>
  );
}

export default ShowText;
